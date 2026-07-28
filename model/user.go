package model

// Author: deepseek-v4-pro / opencode

import "time"

type User struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UID          string    `gorm:"uniqueIndex;size:32" json:"uid"`
	Username     string    `gorm:"uniqueIndex;not null;size:64" json:"username"`
	Password     string    `gorm:"not null;size:255" json:"-"`
	Email        string    `gorm:"size:128" json:"email"`
	Phone        string    `gorm:"size:32" json:"phone"`
	Organization string    `gorm:"size:128" json:"organization"`
	RoleID       *uint     `json:"role_id"`
	Status       int       `gorm:"default:1" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
