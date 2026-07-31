package model

import "time"

type AuditLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	AccessKey string    `gorm:"index;size:255" json:"access_key"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ServerID  string    `gorm:"size:36" json:"server_id"`
	ToolName  string    `gorm:"size:255" json:"tool_name"`
	Success   bool      `json:"success"`
	Message   string    `gorm:"type:text" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
