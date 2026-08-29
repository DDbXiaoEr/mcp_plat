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

import "time"

type AuditLog struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	AccessKeyID   uint      `gorm:"index" json:"access_key_id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	ServerID      string    `gorm:"size:36" json:"server_id"`
	ToolName      string    `gorm:"size:255" json:"tool_name"`
	ClientIP      string    `gorm:"size:64" json:"client_ip"`
	Success       bool      `json:"success"`
	Message       string    `gorm:"type:text" json:"message"`
	CreatedAt     time.Time `json:"created_at"`
	AccessKeyName string    `gorm:"-" json:"access_key_name"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
