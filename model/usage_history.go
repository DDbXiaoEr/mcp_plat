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

type UsageHistory struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	AccessKeyID uint      `json:"access_key_id"`
	Server      string    `gorm:"size:128" json:"server"`
	Endpoint    string    `gorm:"size:255" json:"endpoint"`
	Status      string    `gorm:"size:32" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
