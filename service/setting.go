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
}

func GetSettings() (map[string]json.RawMessage, error) {
	var settings []model.Setting
	if err := database.DB.Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]json.RawMessage, len(settings))
	for _, s := range settings {
		result[s.Key] = json.RawMessage(s.Value)
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

type GatewayStatus struct {
	Configured bool `json:"configured"`
}

func CheckGatewayStatus() GatewayStatus {
	var gw ApiGatewaySetting
	if err := GetSetting("api_gateway", &gw); err != nil {
		return GatewayStatus{Configured: false}
	}
	return GatewayStatus{Configured: gw.AdminURL != "" && gw.AdminKey != ""}
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
