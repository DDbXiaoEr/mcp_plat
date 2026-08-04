package plugins

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"mcp_plat-console/plugin"

	runnerHttp "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	runnerPlugin "github.com/apache/apisix-go-plugin-runner/pkg/plugin"
)

type AuditLogConf struct {
	GrpcAddr        string `json:"grpc_addr"`
	HeaderName      string `json:"header_name"`
	ServerID        string `json:"server_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

type AuditLog struct{}

func init() {
	err := runnerPlugin.RegisterPlugin(&AuditLog{})
	if err != nil {
		log.Fatalf("注册 audit_log 插件失败: %s", err)
	}
}

func (p *AuditLog) Name() string {
	return "audit_log"
}

func (p *AuditLog) ParseConf(in []byte) (interface{}, error) {
	conf := AuditLogConf{
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

func (p *AuditLog) RequestFilter(conf interface{}, w http.ResponseWriter, r runnerHttp.Request) {
	cfg := conf.(AuditLogConf)
	if cfg.GrpcAddr == "" || cfg.ServerID == "" {
		return
	}

	accessKey := extractAccessKey(r, cfg.HeaderName)
	toolName := extractToolName(r)

	var userID uint64
	if cfg.AccessKeySecret != "" && accessKey != "" {
		if claims, err := plugin.ParseAccessKeyWithSecret(accessKey, []byte(cfg.AccessKeySecret)); err == nil {
			userID = uint64(claims.UserID)
		}
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := plugin.LogAccess(ctx, cfg.GrpcAddr, &plugin.LogAccessRequest{
			AccessKey: accessKey,
			UserId:    userID,
			ServerId:  cfg.ServerID,
			ToolName:  toolName,
			Success:   true,
			Message:   "ok",
		})
		if err != nil {
			log.Warnf("audit log 写入失败: %s", err)
		}
	}()
}

func (p *AuditLog) ResponseFilter(conf interface{}, w runnerHttp.Response) {}
