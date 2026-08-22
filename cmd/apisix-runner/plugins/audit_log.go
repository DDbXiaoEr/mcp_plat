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
	GrpcAddr   string         `json:"grpc_addr"`
	GrpcAddrs  []string       `json:"grpc_addrs"`
	HeaderName string         `json:"header_name"`
	ServerID   string         `json:"server_id"`
	lb         *plugin.GrpcLB
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
	addrs := conf.GrpcAddrs
	if len(addrs) == 0 && conf.GrpcAddr != "" {
		addrs = []string{conf.GrpcAddr}
	}
	conf.lb = plugin.GetGrpcLB(addrs)
	return conf, nil
}

func (p *AuditLog) RequestFilter(conf interface{}, w http.ResponseWriter, r runnerHttp.Request) {
	cfg := conf.(AuditLogConf)
	if cfg.lb == nil || cfg.lb.Len() == 0 || cfg.ServerID == "" {
		return
	}

	accessKey := extractAccessKey(r, cfg.HeaderName)
	toolName := extractToolName(r)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := plugin.LogAccess(ctx, cfg.lb, &plugin.LogAccessRequest{
			AccessKey: accessKey,
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
