// Copyright (C) 2026 Zhaoquan Wang
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package model

// Author: deepseek-v4-pro / opencode

import "time"

type User struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UID          string    `gorm:"uniqueIndex;size:32" json:"uid"`
	Username     string    `gorm:"uniqueIndex;not null;size:64" json:"username"`
	Password     string    `gorm:"not null;size:255" json:"-"`
	Name         string    `gorm:"size:64" json:"name"`
	Email        string    `gorm:"size:128" json:"email"`
	Phone        string    `gorm:"size:32" json:"phone"`
	Organization string    `gorm:"size:128" json:"organization"`
	RoleID       *uint     `json:"role_id"`
	Status       int       `gorm:"default:1" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
