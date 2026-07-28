package model

// Author: deepseek-v4-pro / opencode

import "time"

type UsageHistory struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	AccessKeyID uint      `json:"access_key_id"`
	Server      string    `gorm:"size:128" json:"server"`
	Endpoint    string    `gorm:"size:255" json:"endpoint"`
	Status      string    `gorm:"size:32" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
