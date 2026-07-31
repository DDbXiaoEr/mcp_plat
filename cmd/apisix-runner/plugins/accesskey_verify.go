package plugins

// Author: deepseek-v4-pro / opencode

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"mcp_plat-console/plugin"

	runnerHttp "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	runnerPlugin "github.com/apache/apisix-go-plugin-runner/pkg/plugin"
)

type jsonRpcRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type toolCallParams struct {
	Name string `json:"name"`
}

type AccessKeyVerifyConf struct {
	HeaderName    string `json:"header_name"`
	GrpcAddr      string `json:"grpc_addr"`
	AuditGrpcAddr string `json:"audit_grpc_addr"`
	ServerID      string `json:"server_id"`
}

type AccessKeyVerify struct{}

func init() {
	err := runnerPlugin.RegisterPlugin(&AccessKeyVerify{})
	if err != nil {
		log.Fatalf("注册 accesskey_verify 插件失败: %s", err)
	}
}

func (p *AccessKeyVerify) Name() string {
	return "accesskey_verify"
}

func (p *AccessKeyVerify) ParseConf(in []byte) (interface{}, error) {
	conf := AccessKeyVerifyConf{
		HeaderName: "X-Access-Key",
	}
	if err := json.Unmarshal(in, &conf); err != nil {
		return conf, err
	}
	if conf.HeaderName == "" {
		conf.HeaderName = "X-Access-Key"
	}
	return conf, nil
}

func extractAccessKey(r runnerHttp.Request, headerName string) string {
	authHeader := r.Header().Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return r.Header().Get(headerName)
}

func (p *AccessKeyVerify) RequestFilter(conf interface{}, w http.ResponseWriter, r runnerHttp.Request) {
	cfg := conf.(AccessKeyVerifyConf)

	accessKey := extractAccessKey(r, cfg.HeaderName)
	if accessKey == "" {
		if cfg.AuditGrpcAddr != "" {
			go logAuditAsync(cfg.AuditGrpcAddr, accessKey, 0, cfg.ServerID, "", false, "缺少 access key")
		}
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code":401,"message":"缺少 access key"}`))
		return
	}

	grpcAddr := cfg.GrpcAddr

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	requestPath := string(r.Path())

	toolName := ""
	if r.Method() == "POST" {
		body, bodyErr := r.Body()
		if bodyErr == nil && len(body) > 0 {
			var jrpc jsonRpcRequest
			if json.Unmarshal(body, &jrpc) == nil && jrpc.Method == "tools/call" {
				var params toolCallParams
				if json.Unmarshal(jrpc.Params, &params) == nil {
					toolName = params.Name
				}
			}
		}
		if bodyErr != nil {
			log.Warnf("读取请求体失败: %s", bodyErr)
		}
	}

	resp, err := plugin.ValidateAccessKey(ctx, grpcAddr, accessKey, requestPath, cfg.ServerID, toolName)
	if err != nil || !resp.Valid {
		msg := "access key 校验失败"
		if resp != nil && resp.Message != "" {
			msg = resp.Message
		}
		if err != nil {
			msg = "access key 校验异常: " + err.Error()
		}
		if cfg.AuditGrpcAddr != "" {
			userID := uint64(0)
			if resp != nil {
				userID = resp.UserId
			}
			go logAuditAsync(cfg.AuditGrpcAddr, accessKey, userID, cfg.ServerID, toolName, false, msg)
		}
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code":401,"message":"` + msg + `"}`))
		return
	}

	if cfg.AuditGrpcAddr != "" {
		go logAuditAsync(cfg.AuditGrpcAddr, accessKey, resp.UserId, cfg.ServerID, toolName, true, "ok")
	}
}

func logAuditAsync(grpcAddr string, accessKey string, userID uint64, serverID string, toolName string, success bool, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := plugin.LogAccess(ctx, grpcAddr, &plugin.LogAccessRequest{
		AccessKey: accessKey,
		UserId:    userID,
		ServerId:  serverID,
		ToolName:  toolName,
		Success:   success,
		Message:   message,
	})
	if err != nil {
		log.Warnf("audit log 写入失败: %s", err)
	}
}

func (p *AccessKeyVerify) ResponseFilter(conf interface{}, w runnerHttp.Response) {}

