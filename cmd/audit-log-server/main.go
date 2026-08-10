package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mcp_plat-console/auditstore"
	"mcp_plat-console/config"
	"mcp_plat-console/model"
	"mcp_plat-console/plugin"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

type WriterConfig struct {
	BatchSize       int `yaml:"batch_size"`
	FlushIntervalMS int `yaml:"flush_interval_ms"`
}

type Config struct {
	GrpcAddr        string                `yaml:"grpc_addr"`
	Database        config.DatabaseConfig `yaml:"database"`
	AccessKeySecret string                `yaml:"access_key_secret"`
	Writer          WriterConfig          `yaml:"writer"`
	RetentionDays   int                   `yaml:"retention_days"`
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
	if AppConfig.AccessKeySecret == "" {
		fmt.Fprintf(os.Stderr, "config: access_key_secret is required in audit_log_server.yml\n")
		os.Exit(1)
	}
	if AppConfig.Database.Type == "" {
		AppConfig.Database.Type = "sqlite"
	}
	if AppConfig.Database.SQLite.Path == "" {
		AppConfig.Database.SQLite.Path = "audit_log.db"
	}
	if AppConfig.Writer.BatchSize <= 0 {
		AppConfig.Writer.BatchSize = 200
	}
	if AppConfig.Writer.FlushIntervalMS <= 0 {
		AppConfig.Writer.FlushIntervalMS = 1000
	}
}

func initStore(cfg config.DatabaseConfig) auditstore.Store {
	store, err := auditstore.New(cfg)
	if err != nil {
		fmt.Printf("failed to init audit store: %v\n", err)
		os.Exit(1)
	}
	return store
}

// auditWriter 内部缓冲批量写入，降低高频单条落库压力
type auditWriter struct {
	store         auditstore.Store
	ch            chan model.AuditLog
	batchSize     int
	flushInterval time.Duration
}

func newAuditWriter(store auditstore.Store, batchSize int, flushInterval time.Duration) *auditWriter {
	w := &auditWriter{
		store:         store,
		ch:            make(chan model.AuditLog, 4096),
		batchSize:     batchSize,
		flushInterval: flushInterval,
	}
	go w.run()
	return w
}

func (w *auditWriter) submit(entry model.AuditLog) {
	select {
	case w.ch <- entry:
	default:
		fmt.Println("audit-log: 写入缓冲已满，丢弃一条日志")
	}
}

func (w *auditWriter) run() {
	buf := make([]model.AuditLog, 0, w.batchSize)
	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case e := <-w.ch:
			buf = append(buf, e)
			if len(buf) >= w.batchSize {
				w.flush(buf)
				buf = buf[:0]
			}
		case <-ticker.C:
			if len(buf) > 0 {
				w.flush(buf)
				buf = buf[:0]
			}
		}
	}
}

func (w *auditWriter) flush(logs []model.AuditLog) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := w.store.Append(ctx, logs); err != nil {
		fmt.Printf("audit-log: flush %d entries failed: %v\n", len(logs), err)
		return
	}
	fmt.Printf("audit-log: flushed %d entries at %s\n", len(logs), time.Now().Format("2006-01-02 15:04:05"))
}

// startRetention 定期清理关系库中的过期审计日志（ClickHouse 由表 TTL 负责，无需此逻辑）
func startRetention(store auditstore.Store, retentionDays int) {
	if retentionDays <= 0 {
		return
	}
	p, ok := store.(auditstore.Purgable)
	if !ok {
		return
	}

	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			before := time.Now().AddDate(0, 0, -retentionDays)
			n, err := p.Purge(ctx, before)
			cancel()
			if err != nil {
				fmt.Printf("audit-log: retention purge error: %v\n", err)
			} else if n > 0 {
				fmt.Printf("audit-log: retention purged %d entries\n", n)
			}
		}
	}()
}

type auditLogServer struct {
	plugin.UnimplementedAuditLogServiceServer
	writer *auditWriter
	secret []byte
}

func (s *auditLogServer) LogAccess(ctx context.Context, req *plugin.LogAccessRequest) (*plugin.LogAccessResponse, error) {
	userID := uint(req.UserId)
	var keyID uint

	if req.AccessKey != "" {
		if claims, err := plugin.ParseAccessKeyWithSecret(req.AccessKey, s.secret); err == nil {
			if userID == 0 {
				userID = claims.UserID
			}
			keyID = claims.KeyID
		}
	}

	entry := model.AuditLog{
		AccessKeyID: keyID,
		UserID:      userID,
		ServerID:    req.ServerId,
		ToolName:    req.ToolName,
		Success:     req.Success,
		Message:     req.Message,
		CreatedAt:   time.Now(),
	}

	s.writer.submit(entry)

	return &plugin.LogAccessResponse{
		Ok:      true,
		Message: "ok",
	}, nil
}

func main() {
	Load()

	store := initStore(AppConfig.Database)
	writer := newAuditWriter(store, AppConfig.Writer.BatchSize, time.Duration(AppConfig.Writer.FlushIntervalMS)*time.Millisecond)
	startRetention(store, AppConfig.RetentionDays)

	lis, err := net.Listen("tcp", AppConfig.GrpcAddr)
	if err != nil {
		fmt.Printf("gRPC server 监听失败: %v\n", err)
		os.Exit(1)
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("[audit-log] 收到请求 method=%s\n", info.FullMethod)
		return handler(ctx, req)
	}))
	plugin.RegisterAuditLogServiceServer(srv, &auditLogServer{writer: writer, secret: []byte(AppConfig.AccessKeySecret)})

	go func() {
		fmt.Printf("gRPC audit-log-server 已启动，监听 %s（存储类型=%s）\n", AppConfig.GrpcAddr, AppConfig.Database.Type)
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
