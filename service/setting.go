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
	"strings"

	"mcp_plat-console/database"
	"mcp_plat-console/model"

	"gorm.io/gorm/clause"
)

var settingKeys = map[string]bool{
	"log":              true,
	"smtp":             true,
	"auth":             true,
	"user_ops":         true,
	"api_gateway":      true,
	"platform":         true,
	"publish_domain":   true,
	"network_security": true,
	"audit_log":        true,
	"quick_access":     true,
}

func GetSettings() (map[string]json.RawMessage, error) {
	var settings []model.Setting
	if err := database.DB.Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]json.RawMessage, len(settings))
	for _, s := range settings {
		if s.Key == "api_gateway" {
			status, _ := json.Marshal(CheckGatewayStatus())
			result[s.Key] = status
		} else {
			result[s.Key] = json.RawMessage(s.Value)
		}
	}
	if domainsJSON, err := json.Marshal(PublishDomainSetting{Domains: GetPublishDomains()}); err == nil {
		result["publish_domain"] = domainsJSON
	}

	return result, nil
}

func GetSetting(key string, out any) error {
	var setting model.Setting
	if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		return err
	}
	return json.Unmarshal([]byte(setting.Value), out)
}

func GetSettingByKey(key string) (json.RawMessage, error) {
	if !settingKeys[key] {
		return nil, errors.New("不支持的设置项")
	}
	if key == "api_gateway" {
		status, _ := json.Marshal(CheckGatewayStatus())
		return status, nil
	}
	if key == "publish_domain" {
		domains, _ := json.Marshal(PublishDomainSetting{Domains: GetPublishDomains()})
		return domains, nil
	}
	var setting model.Setting
	if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		return nil, errors.New("设置项不存在")
	}
	return json.RawMessage(setting.Value), nil
}

type PlatformSetting struct {
	Name            string `json:"name"`
	LogoURL         string `json:"logoUrl"`
	SiteURL         string `json:"siteUrl"`
	LoginBackground string `json:"loginBackground"`
}

func GetPlatform() PlatformSetting {
	var p PlatformSetting
	if err := GetSetting("platform", &p); err != nil {
		return PlatformSetting{Name: "某某大学"}
	}
	if p.Name == "" {
		p.Name = "某某大学"
	}
	return p
}

type GatewayStatus struct {
	Configured      bool     `json:"configured"`
	Provider        string   `json:"provider"`
	AdminURL        string   `json:"adminUrl"`
	AccesskeyHeader string   `json:"accesskeyHeader"`
	AuthGrpcAddrs   []string `json:"authGrpcAddrs"`
}

func CheckGatewayStatus() GatewayStatus {
	var gw ApiGatewaySetting
	if err := GetSetting("api_gateway", &gw); err != nil {
		return GatewayStatus{Configured: false, AuthGrpcAddrs: []string{}}
	}
	configured := gw.AdminURL != ""
	if gw.Provider != "kong" {
		configured = configured && gw.AdminKey != ""
	}
	addrs := gw.AuthGrpcAddrs
	if len(addrs) == 0 && gw.AuthGrpcAddr != "" {
		addrs = []string{gw.AuthGrpcAddr}
	}
	if addrs == nil {
		addrs = []string{}
	}
	return GatewayStatus{
		Configured:      configured,
		Provider:        gw.Provider,
		AdminURL:        gw.AdminURL,
		AccesskeyHeader: gw.AccesskeyHeader,
		AuthGrpcAddrs:   addrs,
	}
}

type PublishDomainSetting struct {
	Domains []string `json:"domains"`
}

func NormalizePublishDomain(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "http://")
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimRight(raw, "/")
	if idx := strings.Index(raw, "/"); idx >= 0 {
		raw = raw[:idx]
	}
	return raw
}

func GetPublishDomains() []string {
	var setting PublishDomainSetting
	domains := []string{}
	if err := GetSetting("publish_domain", &setting); err == nil {
		for _, raw := range setting.Domains {
			domain := NormalizePublishDomain(raw)
			if domain != "" {
				domains = append(domains, domain)
			}
		}
		return domains
	}

	var legacy struct {
		DefaultPublishDomain string `json:"defaultPublishDomain"`
	}
	if err := GetSetting("api_gateway", &legacy); err == nil {
		if domain := NormalizePublishDomain(legacy.DefaultPublishDomain); domain != "" {
			domains = append(domains, domain)
		}
	}
	return domains
}

func SaveSetting(key string, value json.RawMessage) error {
	if !settingKeys[key] {
		return errors.New("不支持的设置项")
	}
	if !json.Valid(value) {
		return errors.New("设置内容格式错误")
	}

	setting := model.Setting{Key: key, Value: string(value)}
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&setting).Error
}
