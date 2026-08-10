package model

import "time"

type AuditLog struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	AccessKeyID   uint      `gorm:"index" json:"access_key_id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	ServerID      string    `gorm:"size:36" json:"server_id"`
	ToolName      string    `gorm:"size:255" json:"tool_name"`
	Success       bool      `json:"success"`
	Message       string    `gorm:"type:text" json:"message"`
	CreatedAt     time.Time `json:"created_at"`
	AccessKeyName string    `gorm:"-" json:"access_key_name"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
