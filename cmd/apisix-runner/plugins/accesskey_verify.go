// Copyright (C) 2026 Zhaoquan Wang
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package plugins

// Author: deepseek-v4-pro / opencode

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
	HeaderName string         `json:"header_name"`
	GrpcAddr   string         `json:"grpc_addr"`
	GrpcAddrs  []string       `json:"grpc_addrs"`
	ServerID   string         `json:"server_id"`
	lb         *plugin.GrpcLB
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
	addrs := conf.GrpcAddrs
	if len(addrs) == 0 && conf.GrpcAddr != "" {
		addrs = []string{conf.GrpcAddr}
	}
	conf.lb = plugin.GetGrpcLB(addrs)
	return conf, nil
}

func (p *AccessKeyVerify) RequestFilter(conf interface{}, w http.ResponseWriter, r runnerHttp.Request) {
	cfg := conf.(AccessKeyVerifyConf)

	if cfg.lb == nil || cfg.lb.Len() == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"code":503,"message":"access key 校验服务未配置"}`))
		return
	}

	accessKey := extractAccessKey(r, cfg.HeaderName)
	if accessKey == "" {
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code":401,"message":"缺少 access key"}`))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	requestPath := string(r.Path())
	toolName := extractToolName(r)

	resp, err := plugin.ValidateAccessKey(ctx, cfg.lb, accessKey, requestPath, cfg.ServerID, toolName)
	if err != nil || !resp.Valid {
		msg := "access key 校验失败"
		if resp != nil && resp.Message != "" {
			msg = resp.Message
		}
		if err != nil {
			msg = "access key 校验异常: " + err.Error()
		}
		w.Header().Set("WWW-Authenticate", "Bearer")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code":401,"message":"` + msg + `"}`))
		return
	}
}

func (p *AccessKeyVerify) ResponseFilter(conf interface{}, w runnerHttp.Response) {}

