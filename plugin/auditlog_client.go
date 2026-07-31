package plugin

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func LogAccess(ctx context.Context, grpcAddr string, req *LogAccessRequest) (*LogAccessResponse, error) {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := NewAuditLogServiceClient(conn)
	return client.LogAccess(ctx, req)
}
