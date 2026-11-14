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

type GatewayServer struct {
	ID              string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name            string    `gorm:"not null;size:128" json:"name"`
	Provider        string    `gorm:"size:32;default:apisix" json:"provider"`
	AdminURL        string    `gorm:"size:512" json:"admin_url"`
	AdminKey        string    `gorm:"type:text" json:"-"`
	AccesskeyHeader string    `gorm:"size:128" json:"accesskey_header"`
	AuthGrpcAddrs   string    `gorm:"type:text" json:"-"`
	AuditGrpcAddrs  string    `gorm:"type:text" json:"-"`
	Description     string    `gorm:"type:text" json:"description"`
	IsDefault       bool      `gorm:"default:false" json:"is_default"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
