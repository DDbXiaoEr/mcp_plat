package model

// Author: deepseek-v4-pro / opencode

import "time"

type MCPServer struct {
	ID             string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name           string    `gorm:"not null;size:128" json:"name"`
	Address        string    `gorm:"not null;size:512" json:"address"`
	ServiceAddress string    `gorm:"size:256" json:"service_address"`
	Department     string    `gorm:"size:128" json:"department"`
	Protocol       string    `gorm:"size:32;default:streamable http" json:"protocol"`
	Tools          string    `gorm:"type:text" json:"tools"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
