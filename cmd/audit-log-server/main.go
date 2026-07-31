package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	GrpcAddr string         `yaml:"grpc_addr"`
	Database DatabaseConfig `yaml:"database"`
}

var AppConfig *Config

func Load() {
	AppConfig = &Config{}

	configPath := "audit_log_server.yml"
	if p := os.Getenv("AUDIT_LOG_CONFIG_PATH"); p != "" {
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
		AppConfig.GrpcAddr = "127.0.0.1:9091"
	}
	if AppConfig.Database.Type == "" {
		AppConfig.Database.Type = "sqlite"
	}
	if AppConfig.Database.SQLite.Path == "" {
		AppConfig.Database.SQLite.Path = "audit_log.db"
	}
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
		&model.AuditLog{},
	)
	if err != nil {
		fmt.Printf("failed to migrate database: %v\n", err)
		os.Exit(1)
	}

	return db
}

type auditLogServer struct {
	plugin.UnimplementedAuditLogServiceServer
	db *gorm.DB
}

func (s *auditLogServer) LogAccess(ctx context.Context, req *plugin.LogAccessRequest) (*plugin.LogAccessResponse, error) {
	entry := model.AuditLog{
		AccessKey: req.AccessKey,
		UserID:    uint(req.UserId),
		ServerID:  req.ServerId,
		ToolName:  req.ToolName,
		Success:   req.Success,
		Message:   req.Message,
	}

	if err := s.db.Create(&entry).Error; err != nil {
		return &plugin.LogAccessResponse{
			Ok:      false,
			Message: fmt.Sprintf("写入审计日志失败: %v", err),
		}, nil
	}

	return &plugin.LogAccessResponse{
		Ok:      true,
		Message: "ok",
	}, nil
}

func main() {
	Load()

	db := initDB(AppConfig.Database)

	lis, err := net.Listen("tcp", AppConfig.GrpcAddr)
	if err != nil {
		fmt.Printf("gRPC server 监听失败: %v\n", err)
		os.Exit(1)
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("[audit-log] 收到请求 method=%s\n", info.FullMethod)
		return handler(ctx, req)
	}))
	plugin.RegisterAuditLogServiceServer(srv, &auditLogServer{db: db})

	go func() {
		fmt.Printf("gRPC audit-log-server 已启动，监听 %s\n", AppConfig.GrpcAddr)
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
