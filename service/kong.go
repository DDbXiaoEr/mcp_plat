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

package service

// Author: deepseek-v4-pro / opencode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"mcp_plat-console/model"
)

// publishKongServer 发布 MCP 服务器到 Kong 网关。
// 仅创建/更新 上游（Upstream）+ Service + 路由（Route），不部署任何插件。
func publishKongServer(client *http.Client, gw ApiGatewaySetting, srv model.MCPServer) error {
	addresses := parseServiceAddresses(srv.ServiceAddress)
	if len(addresses) == 0 {
		return errors.New("未配置后端服务地址")
	}

	if err := kongRequest(client, gw, http.MethodPut, "/upstreams/"+srv.ID, map[string]interface{}{
		"name":      srv.ID,
		"algorithm": "round-robin",
	}); err != nil {
		return fmt.Errorf("创建上游失败: %w", err)
	}

	if err := kongReplaceTargets(client, gw, srv.ID, addresses); err != nil {
		return fmt.Errorf("设置上游节点失败: %w", err)
	}

	if err := kongRequest(client, gw, http.MethodPut, "/services/"+srv.ID, map[string]interface{}{
		"name":     srv.ID,
		"host":     srv.ID,
		"port":     80,
		"protocol": "http",
	}); err != nil {
		return fmt.Errorf("创建 Service 失败: %w", err)
	}

	path := strings.TrimRight(srv.Address, "/")
	if path == "" {
		path = "/" + srv.ID
	} else if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	routeBody := map[string]interface{}{
		"name":       srv.ID,
		"paths":      []string{path},
		"strip_path": false,
		"service": map[string]interface{}{
			"name": srv.ID,
		},
	}
	if gw.DefaultPublishDomain != "" {
		routeBody["hosts"] = []string{gw.DefaultPublishDomain}
	}
	if err := kongRequest(client, gw, http.MethodPut, "/routes/"+srv.ID, routeBody); err != nil {
		return fmt.Errorf("创建路由失败: %w", err)
	}

	return nil
}

func kongReplaceTargets(client *http.Client, gw ApiGatewaySetting, upstreamID string, addresses []string) error {
	// Kong 3.x 不支持 targets/bulk 与 target upsert，重复 POST 会 409，
	// 因此先删除旧节点，再逐个添加。
	var payload struct {
		Data []struct {
			Target string `json:"target"`
		} `json:"data"`
	}
	if err := kongGet(client, gw, "/upstreams/"+upstreamID+"/targets", &payload); err != nil {
		return err
	}
	for _, t := range payload.Data {
		target := url.PathEscape(t.Target)
		if err := kongRequest(client, gw, http.MethodDelete, "/upstreams/"+upstreamID+"/targets/"+target, nil); err != nil {
			return err
		}
	}
	for _, addr := range addresses {
		if err := kongRequest(client, gw, http.MethodPost, "/upstreams/"+upstreamID+"/targets", map[string]interface{}{
			"target": addr,
			"weight": 100,
		}); err != nil {
			return err
		}
	}
	return nil
}

func kongGet(client *http.Client, gw ApiGatewaySetting, apiPath string, out interface{}) error {
	urlStr := strings.TrimRight(gw.AdminURL, "/") + apiPath
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return fmt.Errorf("构建请求失败: %w", err)
	}
	if gw.AdminKey != "" {
		req.Header.Set("kong-admin-token", gw.AdminKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求 Kong Admin API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		io.Copy(io.Discard, resp.Body)
		return fmt.Errorf("Kong 返回错误 (status=%d)", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func kongRequest(client *http.Client, gw ApiGatewaySetting, method, apiPath string, body interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("构建请求失败: %w", err)
	}

	urlStr := strings.TrimRight(gw.AdminURL, "/") + apiPath
	req, err := http.NewRequest(method, urlStr, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("构建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if gw.AdminKey != "" {
		req.Header.Set("kong-admin-token", gw.AdminKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求 Kong Admin API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		io.Copy(io.Discard, resp.Body)
		return fmt.Errorf("Kong 返回错误 (status=%d)", resp.StatusCode)
	}

	return nil
}
