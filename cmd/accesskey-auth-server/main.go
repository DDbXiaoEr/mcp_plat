package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
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
		fmt.Printf("config: %s not found, using defaults\n", configPath)
		setDefaults()
		return
	}

	if err := yaml.Unmarshal(data, AppConfig); err != nil {
		fmt.Printf("config: failed to parse %s: %v, using defaults\n", configPath, err)
		setDefaults()
		return
	}

	if AppConfig.AccessKeySecret == "" {
		AppConfig.AccessKeySecret = "a8k3x9m2p7q1r6w4v5y0b3n8t2h7j1k5"
	}
	if AppConfig.GrpcAddr == "" {
		AppConfig.GrpcAddr = ":9090"
	}
	if AppConfig.Database.Type == "" {
		AppConfig.Database.Type = "sqlite"
	}
	if AppConfig.Database.SQLite.Path == "" {
		AppConfig.Database.SQLite.Path = "data.db"
	}
}

func setDefaults() {
	AppConfig = &Config{
		AccessKeySecret: "a8k3x9m2p7q1r6w4v5y0b3n8t2h7j1k5",
		GrpcAddr:        ":9090",
		Database: DatabaseConfig{
			Type: "sqlite",
			SQLite: struct {
				Path string `yaml:"path"`
			}{
				Path: "data.db",
			},
		},
	}
}

var serverIDPattern = regexp.MustCompile(`/api/servers/(\d+)`)

type cacheEntry struct {
	serverIDs []uint
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

func (c *roleCache) get(roleID uint) ([]uint, bool) {
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

func (c *roleCache) set(roleID uint, serverIDs []uint) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[roleID] = cacheEntry{
		serverIDs: serverIDs,
		expiresAt: time.Now().Add(5 * time.Minute),
	}
}

func extractServerID(requestPath string) (uint, bool) {
	if requestPath == "" {
		return 0, false
	}

	matches := serverIDPattern.FindStringSubmatch(requestPath)
	if len(matches) < 2 {
		return 0, false
	}

	id, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(id), true
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

	serverID, hasServerID := extractServerID(req.RequestPath)

	if !hasServerID {
		return &plugin.ValidateResponse{
			Valid:  true,
			UserId: uint64(claims.UserID),
			Role:   claims.Role,
		}, nil
	}

	var role model.Role
	if err := s.db.Where("name = ?", claims.Role).First(&role).Error; err != nil {
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

		serverIDs = make([]uint, 0, len(roleServers))
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
				Role:   claims.Role,
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
