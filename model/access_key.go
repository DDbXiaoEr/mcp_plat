package model

// Author: deepseek-v4-pro / opencode

import "time"

type AccessKey struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	Name      string     `gorm:"not null;size:128" json:"name"`
	Key       string     `gorm:"uniqueIndex;not null;size:255" json:"key"`
	Enabled   bool       `gorm:"default:true" json:"enabled"`
	ExpiredAt *time.Time `json:"expired_at"`
	Servers   string     `gorm:"type:text" json:"servers"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
