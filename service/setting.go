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

type GatewayStatus struct {
	Configured bool   `json:"configured"`
	Provider   string `json:"provider"`
	AdminURL   string `json:"adminUrl"`
}

func CheckGatewayStatus() GatewayStatus {
	var gw ApiGatewaySetting
	if err := GetSetting("api_gateway", &gw); err != nil {
		return GatewayStatus{Configured: false}
	}
	configured := gw.AdminURL != ""
	if gw.Provider != "kong" {
		configured = configured && gw.AdminKey != ""
	}
	return GatewayStatus{
		Configured: configured,
		Provider:   gw.Provider,
		AdminURL:   gw.AdminURL,
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
