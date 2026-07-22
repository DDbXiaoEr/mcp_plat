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

type AccessKeyVerifyConf struct {
	HeaderName string `json:"header_name"`
	GrpcAddr   string `json:"grpc_addr"`
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

func (p *AccessKeyVerify) RequestFilter(conf interface{}, w http.ResponseWriter, r runnerHttp.Request) {
	cfg := conf.(AccessKeyVerifyConf)

	accessKey := r.Header().Get(cfg.HeaderName)
	if accessKey == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code":401,"message":"缺少 access key"}`))
		return
	}

	grpcAddr := cfg.GrpcAddr

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	requestPath := string(r.Path())

	resp, err := plugin.ValidateAccessKey(ctx, grpcAddr, accessKey, requestPath)
	if err != nil || !resp.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		msg := "access key 校验失败"
		if resp != nil && resp.Message != "" {
			msg = resp.Message
		}
		if err != nil {
			msg = "access key 校验异常: " + err.Error()
		}
		w.Write([]byte(`{"code":401,"message":"` + msg + `"}`))
		return
	}
}

func (p *AccessKeyVerify) ResponseFilter(conf interface{}, w runnerHttp.Response) {}

