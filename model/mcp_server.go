package model

import "time"

type MCPServer struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	Name       string    `gorm:"not null;size:128" json:"name"`
	Address    string    `gorm:"not null;size:512" json:"address"`
	Department string    `gorm:"size:128" json:"department"`
	Protocol   string    `gorm:"size:32;default:streamable http" json:"protocol"`
	Tools      string    `gorm:"type:text" json:"tools"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
