package model

// Author: deepseek-v4-pro / opencode

type RoleServer struct {
	RoleID   uint `gorm:"primaryKey" json:"role_id"`
	ServerID string `gorm:"primaryKey;type:varchar(36)" json:"server_id"`
}
