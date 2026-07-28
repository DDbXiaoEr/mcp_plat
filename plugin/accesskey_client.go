package plugin

// Author: deepseek-v4-pro / opencode

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ValidateAccessKey(ctx context.Context, grpcAddr string, key string, requestPath string, serverID string, toolName string) (*ValidateResponse, error) {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := NewAccessKeyServiceClient(conn)
	return client.Validate(ctx, &ValidateRequest{AccessKey: key, RequestPath: requestPath, ServerId: serverID, ToolName: toolName})
}
