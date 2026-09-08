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
	"fmt"
	"net/http"
	"sync"
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

// pendingAudit 保存请求阶段捕获的元数据，等待响应阶段根据真实结果写审计。
type pendingAudit struct {
	accessKey string
	serverID  string
	toolName  string
	clientIP  string
	lb        *plugin.GrpcLB
	createdAt time.Time
}

var (
	pendingMu   sync.Mutex
	pendingLogs = make(map[uint32]*pendingAudit)
	janitorOnce sync.Once
	// fallbackTTL 超过该时长仍未收到响应阶段回调的请求，由 janitor 兜底记为失败，
	// 避免被前置拦截/上游无响应等未触发响应阶段的请求漏记。
	fallbackTTL = 60 * time.Second
)

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

// RequestFilter 只捕获本次请求的上下文，不落审计；真实成功与否在响应阶段判定。
func (p *AuditLog) RequestFilter(conf interface{}, w http.ResponseWriter, r runnerHttp.Request) {
	cfg := conf.(AuditLogConf)
	if cfg.lb == nil || cfg.lb.Len() == 0 || cfg.ServerID == "" {
		return
	}

	clientIP := ""
	if ip := r.SrcIP(); ip != nil {
		clientIP = ip.String()
	}

	pe := &pendingAudit{
		accessKey: extractAccessKey(r, cfg.HeaderName),
		serverID:  cfg.ServerID,
		toolName:  extractToolName(r),
		clientIP:  clientIP,
		lb:        cfg.lb,
		createdAt: time.Now(),
	}

	pendingMu.Lock()
	pendingLogs[r.ID()] = pe
	pendingMu.Unlock()

	janitorOnce.Do(startJanitor)
}

// ResponseFilter 在拿到上游真实响应后落审计，并解析结果判定调用是否成功。
func (p *AuditLog) ResponseFilter(conf interface{}, w runnerHttp.Response) {
	id := w.ID()

	pendingMu.Lock()
	pe := pendingLogs[id]
	if pe != nil {
		delete(pendingLogs, id)
	}
	pendingMu.Unlock()

	if pe == nil {
		return
	}

	code := w.StatusCode()
	body, _ := w.ReadBody()
	success, message := resolveResult(code, body)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := plugin.LogAccess(ctx, pe.lb, &plugin.LogAccessRequest{
			AccessKey: pe.accessKey,
			ServerId:  pe.serverID,
			ToolName:  pe.toolName,
			Success:   success,
			Message:   message,
			ClientIp:  pe.clientIP,
		})
		if err != nil {
			log.Warnf("audit log 写入失败: %s", err)
		}
	}()
}

// resolveResult 依据 HTTP 状态码与 JSON-RPC 响应体判定调用是否成功。
// 上游 MCP 对工具级错误也返回 HTTP 200（error 成员或 result.isError=true），
// 因此不能只看状态码。
func resolveResult(status int, body []byte) (bool, string) {
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return false, fmt.Sprintf("HTTP %d %s", status, http.StatusText(status))
	}

	if len(body) == 0 {
		return true, "ok"
	}

	var jr struct {
		Error  *json.RawMessage `json:"error"`
		Result json.RawMessage  `json:"result"`
	}
	if err := json.Unmarshal(body, &jr); err != nil {
		return true, "ok"
	}

	if jr.Error != nil {
		msg := "rpc error"
		var e struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(*jr.Error, &e); err == nil {
			if e.Message != "" {
				msg = fmt.Sprintf("rpc error code=%d %s", e.Code, e.Message)
			} else {
				msg = fmt.Sprintf("rpc error code=%d", e.Code)
			}
		} else {
			msg = "rpc error " + string(*jr.Error)
		}
		return false, msg
	}

	if len(jr.Result) > 0 {
		var r struct {
			IsError bool `json:"isError"`
		}
		if err := json.Unmarshal(jr.Result, &r); err == nil && r.IsError {
			return false, "tool returned isError"
		}
	}

	return true, "ok"
}

// startJanitor 兜底处理从未触发响应阶段的请求（前置拦截、上游无响应等），
// 在 fallbackTTL 后记为失败，避免漏记。
func startJanitor() {
	go func() {
		for {
			time.Sleep(fallbackTTL)

			now := time.Now()
			var stale []*pendingAudit

			pendingMu.Lock()
			for id, pe := range pendingLogs {
				if now.Sub(pe.createdAt) > fallbackTTL {
					stale = append(stale, pe)
					delete(pendingLogs, id)
				}
			}
			pendingMu.Unlock()

			for _, pe := range stale {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				_, err := plugin.LogAccess(ctx, pe.lb, &plugin.LogAccessRequest{
					AccessKey: pe.accessKey,
					ServerId:  pe.serverID,
					ToolName:  pe.toolName,
					Success:   false,
					Message:   "未收到上游响应（可能被网关前置拦截）",
					ClientIp:  pe.clientIP,
				})
				cancel()
				if err != nil {
					log.Warnf("audit log fallback 写入失败: %s", err)
				}
			}
		}
	}()
}
