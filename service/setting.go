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
	Configured           bool     `json:"configured"`
	Provider             string   `json:"provider"`
	AdminURL             string   `json:"adminUrl"`
	DefaultPublishDomain string   `json:"defaultPublishDomain"`
	AccesskeyHeader      string   `json:"accesskeyHeader"`
	AuthGrpcAddrs        []string `json:"authGrpcAddrs"`
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
		Configured:           configured,
		Provider:             gw.Provider,
		AdminURL:             gw.AdminURL,
		DefaultPublishDomain: gw.DefaultPublishDomain,
		AccesskeyHeader:      gw.AccesskeyHeader,
		AuthGrpcAddrs:        addrs,
	}
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
