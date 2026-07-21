package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mcp_plat-console/plugin"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AccessKeySecret string `yaml:"access_key_secret"`
	GrpcAddr        string `yaml:"grpc_addr"`
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
		// 请修改为自定义密钥
		AppConfig.AccessKeySecret = "a8k3x9m2p7q1r6w4v5y0b3n8t2h7j1k5"
	}
	if AppConfig.GrpcAddr == "" {
		AppConfig.GrpcAddr = ":9090"
	}
}

func setDefaults() {
	AppConfig = &Config{
		AccessKeySecret: "a8k3x9m2p7q1r6w4v5y0b3n8t2h7j1k5", // 请修改为自定义密钥
		GrpcAddr:        ":9090",
	}
}

type accessKeyServer struct {
	plugin.UnimplementedAccessKeyServiceServer
	secret []byte
}

func (s *accessKeyServer) Validate(ctx context.Context, req *plugin.ValidateRequest) (*plugin.ValidateResponse, error) {
	claims, err := plugin.ParseAccessKeyWithSecret(req.AccessKey, s.secret)
	if err != nil {
		return &plugin.ValidateResponse{
			Valid:   false,
			Message: err.Error(),
		}, nil
	}

	return &plugin.ValidateResponse{
		Valid:  true,
		UserId: uint64(claims.UserID),
		Role:   claims.Role,
	}, nil
}

func main() {
	Load()

	lis, err := net.Listen("tcp", AppConfig.GrpcAddr)
	if err != nil {
		fmt.Printf("gRPC server 监听失败: %v\n", err)
		os.Exit(1)
	}

	srv := grpc.NewServer()
	plugin.RegisterAccessKeyServiceServer(srv, &accessKeyServer{
		secret: []byte(AppConfig.AccessKeySecret),
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
