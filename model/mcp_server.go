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

type MCPServer struct {
	ID              string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name            string    `gorm:"not null;size:128" json:"name"`
	Address         string    `gorm:"not null;size:512" json:"address"`
	ServiceAddress  string    `gorm:"type:text" json:"service_address"`
	Department      string    `gorm:"size:128" json:"department"`
	Protocol        string    `gorm:"size:32;default:streamable http" json:"protocol"`
	ProtocolVersion string    `gorm:"size:16;default:2026-07-28" json:"protocol_version"`
	Tools           string    `gorm:"type:text" json:"tools"`
	Description     string    `gorm:"type:text" json:"description"`
	Status          string    `gorm:"size:32;default:unpublished" json:"status"`
	AuthEnabled     bool      `gorm:"default:true" json:"auth_enabled"`
	GatewayRoute    string    `gorm:"type:text" json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
