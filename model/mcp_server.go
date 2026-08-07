package model

// Author: deepseek-v4-pro / opencode

import "time"

type MCPServer struct {
	ID             string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name           string    `gorm:"not null;size:128" json:"name"`
	Address        string    `gorm:"not null;size:512" json:"address"`
	ServiceAddress string    `gorm:"type:text" json:"service_address"`
	Department     string    `gorm:"size:128" json:"department"`
	Protocol        string    `gorm:"size:32;default:streamable http" json:"protocol"`
	ProtocolVersion string    `gorm:"size:16;default:2026-07-28" json:"protocol_version"`
	Tools          string    `gorm:"type:text" json:"tools"`
	Description    string    `gorm:"type:text" json:"description"`
	Status         string    `gorm:"size:32;default:unpublished" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
