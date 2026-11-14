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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"mcp_plat-console/database"
	"mcp_plat-console/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GatewayServerInput struct {
	Name            string   `json:"name"`
	Provider        string   `json:"provider"`
	AdminURL        string   `json:"admin_url"`
	AdminKey        string   `json:"admin_key"`
	AccesskeyHeader string   `json:"accesskey_header"`
	AuthGrpcAddrs   []string `json:"auth_grpc_addrs"`
	Description     string   `json:"description"`
}

type GatewayServerView struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Provider        string    `json:"provider"`
	AdminURL        string    `json:"admin_url"`
	AccesskeyHeader string    `json:"accesskey_header"`
	AuthGrpcAddrs   []string  `json:"auth_grpc_addrs"`
	Description     string    `json:"description"`
	IsDefault       bool      `json:"is_default"`
	HasAdminKey     bool      `json:"has_admin_key"`
	Status          string    `json:"status"`
	StatusMessage   string    `json:"status_message"`
	LatencyMS       int64     `json:"latency_ms"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type GatewayCheckResult struct {
	Online    bool   `json:"online"`
	Message   string `json:"message"`
	LatencyMS int64  `json:"latency_ms"`
}

type GatewayCheckInput struct {
	ID string `json:"id"`
	GatewayServerInput
}

func normalizeGatewayProvider(provider string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "", "apisix":
		return "apisix", nil
	case "kong":
		return "kong", nil
	}
	return "", errors.New("不支持的网关类型")
}

func normalizeGatewayInput(input *GatewayServerInput) error {
	input.Name = strings.TrimSpace(input.Name)
	input.Provider = strings.ToLower(strings.TrimSpace(input.Provider))
	input.AdminURL = strings.TrimRight(strings.TrimSpace(input.AdminURL), "/")
	input.AdminKey = strings.TrimSpace(input.AdminKey)
	input.AccesskeyHeader = strings.TrimSpace(input.AccesskeyHeader)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return errors.New("网关名称不能为空")
	}

	provider, err := normalizeGatewayProvider(input.Provider)
	if err != nil {
		return err
	}
	input.Provider = provider

	if input.AdminURL == "" {
		return errors.New("API 网关 Admin API 地址未配置")
	}
	u, err := url.Parse(input.AdminURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("API 网关 Admin API 地址格式错误")
	}

	addrs := make([]string, 0, len(input.AuthGrpcAddrs))
	for _, addr := range input.AuthGrpcAddrs {
		if addr = strings.TrimSpace(addr); addr != "" {
			addrs = append(addrs, addr)
		}
	}
	input.AuthGrpcAddrs = addrs

	return nil
}

func parseGatewayAuthAddrs(raw string) []string {
	addrs := []string{}
	if raw == "" {
		return addrs
	}
	if err := json.Unmarshal([]byte(raw), &addrs); err != nil {
		return []string{}
	}
	if addrs == nil {
		return []string{}
	}
	return addrs
}

func gatewayServerView(gw model.GatewayServer) GatewayServerView {
	return GatewayServerView{
		ID:              gw.ID,
		Name:            gw.Name,
		Provider:        gw.Provider,
		AdminURL:        gw.AdminURL,
		AccesskeyHeader: gw.AccesskeyHeader,
		AuthGrpcAddrs:   parseGatewayAuthAddrs(gw.AuthGrpcAddrs),
		Description:     gw.Description,
		IsDefault:       gw.IsDefault,
		HasAdminKey:     gw.AdminKey != "",
		CreatedAt:       gw.CreatedAt,
		UpdatedAt:       gw.UpdatedAt,
	}
}

func CheckGatewayConnectivity(provider, adminURL, adminKey string, timeout time.Duration) GatewayCheckResult {
	provider, err := normalizeGatewayProvider(provider)
	if err != nil {
		return GatewayCheckResult{Message: err.Error()}
	}
	adminURL = strings.TrimRight(strings.TrimSpace(adminURL), "/")
	if adminURL == "" {
		return GatewayCheckResult{Message: "API 网关 Admin API 地址未配置"}
	}

	var reqURL, headerName string
	switch provider {
	case "kong":
		reqURL = adminURL + "/status"
		headerName = "Kong-Admin-Token"
	default:
		reqURL = adminURL + "/apisix/admin/routes?page=1&page_size=10"
		headerName = "X-API-KEY"
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return GatewayCheckResult{Message: fmt.Sprintf("连接网关失败: %v", err)}
	}
	if adminKey != "" {
		req.Header.Set(headerName, adminKey)
	}

	client := &http.Client{Timeout: timeout}
	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return GatewayCheckResult{LatencyMS: latency, Message: fmt.Sprintf("连接网关失败: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return GatewayCheckResult{
			LatencyMS: latency,
			Message:   fmt.Sprintf("网关返回错误 (status=%d): %s", resp.StatusCode, strings.TrimSpace(string(body))),
		}
	}

	return GatewayCheckResult{Online: true, LatencyMS: latency, Message: "连接成功"}
}

func ListGatewayServers() ([]GatewayServerView, error) {
	var gateways []model.GatewayServer
	if err := database.DB.Order("is_default DESC, created_at ASC").Find(&gateways).Error; err != nil {
		return nil, err
	}

	views := make([]GatewayServerView, len(gateways))
	var wg sync.WaitGroup
	for i, gw := range gateways {
		views[i] = gatewayServerView(gw)
		wg.Add(1)
		go func(i int, gw model.GatewayServer) {
			defer wg.Done()
			result := CheckGatewayConnectivity(gw.Provider, gw.AdminURL, gw.AdminKey, 3*time.Second)
			if result.Online {
				views[i].Status = "online"
			} else {
				views[i].Status = "offline"
			}
			views[i].StatusMessage = result.Message
			views[i].LatencyMS = result.LatencyMS
		}(i, gw)
	}
	wg.Wait()

	return views, nil
}

func CreateGatewayServer(input GatewayServerInput) (GatewayServerView, error) {
	if err := normalizeGatewayInput(&input); err != nil {
		return GatewayServerView{}, err
	}
	if input.Provider != "kong" && input.AdminKey == "" {
		return GatewayServerView{}, errors.New("API 网关 Admin API Key 未配置")
	}

	var count int64
	if err := database.DB.Model(&model.GatewayServer{}).Count(&count).Error; err != nil {
		return GatewayServerView{}, errors.New("查询网关失败")
	}

	addrsJSON, err := json.Marshal(input.AuthGrpcAddrs)
	if err != nil {
		return GatewayServerView{}, errors.New("认证 gRPC 地址格式错误")
	}

	gw := model.GatewayServer{
		ID:              uuid.NewString(),
		Name:            input.Name,
		Provider:        input.Provider,
		AdminURL:        input.AdminURL,
		AdminKey:        input.AdminKey,
		AccesskeyHeader: input.AccesskeyHeader,
		AuthGrpcAddrs:   string(addrsJSON),
		Description:     input.Description,
		IsDefault:       count == 0,
	}
	if err := database.DB.Create(&gw).Error; err != nil {
		return GatewayServerView{}, errors.New("创建网关失败")
	}

	return gatewayServerView(gw), nil
}

func UpdateGatewayServer(id string, input GatewayServerInput) error {
	var gw model.GatewayServer
	if err := database.DB.Where("id = ?", id).First(&gw).Error; err != nil {
		return errors.New("网关服务器不存在")
	}
	if err := normalizeGatewayInput(&input); err != nil {
		return err
	}

	if input.AdminKey == "" {
		if gw.AdminKey == "" && input.Provider != "kong" {
			return errors.New("API 网关 Admin API Key 未配置")
		}
		input.AdminKey = gw.AdminKey
	}

	addrsJSON, err := json.Marshal(input.AuthGrpcAddrs)
	if err != nil {
		return errors.New("认证 gRPC 地址格式错误")
	}

	if err := database.DB.Model(&model.GatewayServer{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":             input.Name,
			"provider":         input.Provider,
			"admin_url":        input.AdminURL,
			"admin_key":        input.AdminKey,
			"accesskey_header": input.AccesskeyHeader,
			"auth_grpc_addrs":  string(addrsJSON),
			"description":      input.Description,
		}).Error; err != nil {
		return errors.New("更新网关失败")
	}
	return nil
}

func DeleteGatewayServer(id string) error {
	var gw model.GatewayServer
	if err := database.DB.Where("id = ?", id).First(&gw).Error; err != nil {
		return errors.New("网关服务器不存在")
	}

	if err := database.DB.Delete(&model.GatewayServer{}, "id = ?", id).Error; err != nil {
		return errors.New("删除网关失败")
	}

	if gw.IsDefault {
		var next model.GatewayServer
		if err := database.DB.Order("created_at ASC").First(&next).Error; err == nil {
			if err := database.DB.Model(&model.GatewayServer{}).Where("id = ?", next.ID).
				Update("is_default", true).Error; err != nil {
				return errors.New("设置默认网关失败")
			}
		}
	}

	return nil
}

func SetDefaultGatewayServer(id string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var gw model.GatewayServer
		if err := tx.Where("id = ?", id).First(&gw).Error; err != nil {
			return errors.New("网关服务器不存在")
		}
		if err := tx.Model(&model.GatewayServer{}).Where("is_default = ?", true).
			Update("is_default", false).Error; err != nil {
			return errors.New("设置默认网关失败")
		}
		if err := tx.Model(&model.GatewayServer{}).Where("id = ?", id).
			Update("is_default", true).Error; err != nil {
			return errors.New("设置默认网关失败")
		}
		return nil
	})
}

func TestGatewayServer(id string) (GatewayCheckResult, error) {
	var gw model.GatewayServer
	if err := database.DB.Where("id = ?", id).First(&gw).Error; err != nil {
		return GatewayCheckResult{}, errors.New("网关服务器不存在")
	}
	return CheckGatewayConnectivity(gw.Provider, gw.AdminURL, gw.AdminKey, 5*time.Second), nil
}

func TestGatewayServerConfig(input GatewayCheckInput) (GatewayCheckResult, error) {
	if err := normalizeGatewayInput(&input.GatewayServerInput); err != nil {
		return GatewayCheckResult{}, err
	}
	if input.AdminKey == "" && input.ID != "" {
		var gw model.GatewayServer
		if err := database.DB.Where("id = ?", input.ID).First(&gw).Error; err == nil {
			input.AdminKey = gw.AdminKey
		}
	}
	return CheckGatewayConnectivity(input.Provider, input.AdminURL, input.AdminKey, 5*time.Second), nil
}

func resolveGateway(id string) (model.GatewayServer, error) {
	var gw model.GatewayServer
	query := database.DB
	if id != "" {
		query = query.Where("id = ?", id)
	}
	if err := query.Order("is_default DESC, created_at ASC").First(&gw).Error; err != nil {
		if id != "" {
			return gw, errors.New("网关服务器不存在")
		}
		return gw, errors.New("请先在 API 网关服务器管理中添加网关")
	}
	if gw.AdminURL == "" {
		return gw, errors.New("API 网关 Admin API 地址未配置")
	}
	if gw.Provider != "kong" && gw.AdminKey == "" {
		return gw, errors.New("API 网关 Admin API Key 未配置")
	}
	return gw, nil
}

func gatewaySettingOf(gw model.GatewayServer) ApiGatewaySetting {
	return ApiGatewaySetting{
		Provider:        gw.Provider,
		AdminURL:        gw.AdminURL,
		AdminKey:        gw.AdminKey,
		AuthGrpcAddrs:   parseGatewayAuthAddrs(gw.AuthGrpcAddrs),
		AccesskeyHeader: gw.AccesskeyHeader,
	}
}

func GetDefaultGatewaySetting() (ApiGatewaySetting, error) {
	gw, err := resolveGateway("")
	if err != nil {
		return ApiGatewaySetting{}, err
	}
	return gatewaySettingOf(gw), nil
}

func MigrateLegacyGatewaySetting() {
	var count int64
	if err := database.DB.Model(&model.GatewayServer{}).Count(&count).Error; err != nil || count > 0 {
		return
	}

	var setting model.Setting
	if err := database.DB.Where("key = ?", "api_gateway").First(&setting).Error; err != nil {
		return
	}

	var legacySetting struct {
		ApiGatewaySetting
		DefaultPublishDomain string `json:"defaultPublishDomain"`
	}
	if err := json.Unmarshal([]byte(setting.Value), &legacySetting); err != nil {
		return
	}
	legacy := legacySetting.ApiGatewaySetting
	if strings.TrimSpace(legacy.AdminURL) == "" {
		return
	}

	provider, err := normalizeGatewayProvider(legacy.Provider)
	if err != nil {
		provider = "apisix"
	}

	addrs := legacy.AuthGrpcAddrs
	if len(addrs) == 0 && legacy.AuthGrpcAddr != "" {
		addrs = []string{legacy.AuthGrpcAddr}
	}
	addrsJSON, _ := json.Marshal(addrs)

	gw := model.GatewayServer{
		ID:              uuid.NewString(),
		Name:            "默认网关",
		Provider:        provider,
		AdminURL:        strings.TrimRight(strings.TrimSpace(legacy.AdminURL), "/"),
		AdminKey:        legacy.AdminKey,
		AccesskeyHeader: legacy.AccesskeyHeader,
		AuthGrpcAddrs:   string(addrsJSON),
		IsDefault:       true,
	}
	if err := database.DB.Create(&gw).Error; err != nil {
		return
	}

	if domain := NormalizePublishDomain(legacySetting.DefaultPublishDomain); domain != "" {
		var existing model.Setting
		if err := database.DB.Where("key = ?", "publish_domain").First(&existing).Error; err != nil {
			if domainsJSON, err := json.Marshal(PublishDomainSetting{Domains: []string{domain}}); err == nil {
				_ = SaveSetting("publish_domain", domainsJSON)
			}
		}
	}

	database.DB.Delete(&model.Setting{}, "key = ?", "api_gateway")
}
