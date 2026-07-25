package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"mcp_plat-console/model"
	"mcp_plat-console/plugin"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Type     string `yaml:"type"`
	SQLite   struct {
		Path string `yaml:"path"`
	} `yaml:"sqlite"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"postgres"`
}

type Config struct {
	AccessKeySecret string         `yaml:"access_key_secret"`
	GrpcAddr        string         `yaml:"grpc_addr"`
	Database        DatabaseConfig `yaml:"database"`
}

var AppConfig *Config

func Load() {
	AppConfig = &Config{}

	configPath := "grpc_server.yaml"
	if p := os.Getenv("GRPC_CONFIG_PATH"); p != "" {
		configPath = p
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %s not found\n", configPath)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(data, AppConfig); err != nil {
		fmt.Fprintf(os.Stderr, "config: failed to parse %s: %v\n", configPath, err)
		os.Exit(1)
	}

	if AppConfig.GrpcAddr == "" {
		AppConfig.GrpcAddr = "127.0.0.1:9090"
	}
	if AppConfig.Database.Type == "" {
		AppConfig.Database.Type = "sqlite"
	}
	if AppConfig.Database.SQLite.Path == "" {
		AppConfig.Database.SQLite.Path = "data.db"
	}

	if AppConfig.AccessKeySecret == "" {
		fmt.Fprintf(os.Stderr, "config: access_key_secret is required in grpc_server.yaml\n")
		os.Exit(1)
	}
	if len(AppConfig.AccessKeySecret) < 32 {
		fmt.Fprintf(os.Stderr, "config: access_key_secret must be at least 32 characters\n")
		os.Exit(1)
	}
}

var serverIDPattern = regexp.MustCompile(`^/([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})(?:/|$)`)
var serverIDPatternGeneric = regexp.MustCompile(`^/([0-9a-fA-F\-]+)`)

type cacheEntry struct {
	serverIDs []string
	expiresAt time.Time
}

type roleCache struct {
	mu    sync.RWMutex
	items map[uint]cacheEntry
}

func newRoleCache() *roleCache {
	return &roleCache{
		items: make(map[uint]cacheEntry),
	}
}

func (c *roleCache) get(roleID uint) ([]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.items[roleID]
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.serverIDs, true
}

func (c *roleCache) set(roleID uint, serverIDs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[roleID] = cacheEntry{
		serverIDs: serverIDs,
		expiresAt: time.Now().Add(5 * time.Minute),
	}
}

func extractServerID(requestPath string) (string, bool) {
	if requestPath == "" {
		return "", false
	}

	if matches := serverIDPattern.FindStringSubmatch(requestPath); len(matches) >= 2 {
		return matches[1], true
	}
	if matches := serverIDPatternGeneric.FindStringSubmatch(requestPath); len(matches) >= 2 {
		return matches[1], true
	}

	return "", false
}

type accessKeyServer struct {
	plugin.UnimplementedAccessKeyServiceServer
	secret []byte
	db     *gorm.DB
	cache  *roleCache
}

func (s *accessKeyServer) Validate(ctx context.Context, req *plugin.ValidateRequest) (*plugin.ValidateResponse, error) {
	claims, err := plugin.ParseAccessKeyWithSecret(req.AccessKey, s.secret)
	if err != nil {
		return &plugin.ValidateResponse{
			Valid:   false,
			Message: err.Error(),
		}, nil
	}

	var ak model.AccessKey
	if err := s.db.Where("key = ?", req.AccessKey).First(&ak).Error; err != nil {
		return &plugin.ValidateResponse{
			Valid:   false,
			Message: "access key 不存在或已删除",
		}, nil
	}
	if !ak.Enabled {
		return &plugin.ValidateResponse{
			Valid:   false,
			Message: "access key 已被禁用",
		}, nil
	}
	if ak.ExpiredAt != nil && time.Now().After(*ak.ExpiredAt) {
		return &plugin.ValidateResponse{
			Valid:   false,
			Message: "access key 已过期",
		}, nil
	}

	var user model.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil || user.Status == 0 {
		return &plugin.ValidateResponse{
			Valid:   false,
			Message: "用户不存在或已被禁用",
		}, nil
	}

	roleName := claims.Role
	if user.RoleID != nil {
		var currentRole model.Role
		if err := s.db.First(&currentRole, *user.RoleID).Error; err == nil {
			roleName = currentRole.Name
		}
	}

	serverID := req.ServerId
	if serverID == "" {
		var ok bool
		serverID, ok = extractServerID(req.RequestPath)
		if !ok {
			return &plugin.ValidateResponse{
				Valid:   false,
				Message: "无法识别请求的 MCP 服务器",
			}, nil
		}
	}

	var role model.Role
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return &plugin.ValidateResponse{
			Valid:   false,
			Message: "角色不存在",
		}, nil
	}

	serverIDs, found := s.cache.get(role.ID)
	if !found {
		var roleServers []model.RoleServer
		if err := s.db.Where("role_id = ?", role.ID).Find(&roleServers).Error; err != nil {
			return &plugin.ValidateResponse{
				Valid:   false,
				Message: "查询角色权限失败",
			}, nil
		}

		serverIDs = make([]string, 0, len(roleServers))
		for _, rs := range roleServers {
			serverIDs = append(serverIDs, rs.ServerID)
		}

		s.cache.set(role.ID, serverIDs)
	}

	for _, allowedID := range serverIDs {
		if allowedID == serverID {
			return &plugin.ValidateResponse{
				Valid:  true,
				UserId: uint64(claims.UserID),
				Role:   roleName,
			}, nil
		}
	}

	return &plugin.ValidateResponse{
		Valid:   false,
		Message: "无权访问该 MCP 服务器",
	}, nil
}

func initDB(cfg DatabaseConfig) *gorm.DB {
	var dialector gorm.Dialector

	switch cfg.Type {
	case "postgres":
		pg := cfg.Postgres
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			pg.Host, pg.User, pg.Password, pg.DBName, pg.Port)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.SQLite.Path)
	default:
		fmt.Printf("unsupported database type: %s\n", cfg.Type)
		os.Exit(1)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Printf("failed to connect database: %v\n", err)
		os.Exit(1)
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.AccessKey{},
		&model.UsageHistory{},
		&model.MCPServer{},
		&model.Role{},
		&model.RoleServer{},
		&model.Setting{},
	)
	if err != nil {
		fmt.Printf("failed to migrate database: %v\n", err)
		os.Exit(1)
	}

	return db
}

func main() {
	Load()

	db := initDB(AppConfig.Database)

	lis, err := net.Listen("tcp", AppConfig.GrpcAddr)
	if err != nil {
		fmt.Printf("gRPC server 监听失败: %v\n", err)
		os.Exit(1)
	}

	srv := grpc.NewServer()
	plugin.RegisterAccessKeyServiceServer(srv, &accessKeyServer{
		secret: []byte(AppConfig.AccessKeySecret),
		db:     db,
		cache:  newRoleCache(),
	})

	go func() {
		fmt.Printf("gRPC accesskey-server 已启动，监听 %s\n", AppConfig.GrpcAddr)
		if err := srv.Serve(lis); err != nil {
			fmt.Printf("gRPC server 异常: %v\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("gRPC server 正在关闭...")
	srv.GracefulStop()
}
