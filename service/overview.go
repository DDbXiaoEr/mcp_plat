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

package service

// Author: deepseek-v4-pro / opencode

import (
	"context"
	"encoding/json"
	"time"

	"mcp_plat-console/database"
	"mcp_plat-console/model"
)

type OverviewStats struct {
	Users      int64 `json:"users"`
	Servers    int64 `json:"servers"`
	Tools      int64 `json:"tools"`
	TodayCalls int64 `json:"today_calls"`
}

func GetOverviewStats() (*OverviewStats, error) {
	stats := &OverviewStats{}

	var userCount int64
	if err := database.DB.Model(&model.User{}).Count(&userCount).Error; err != nil {
		return nil, err
	}
	stats.Users = userCount

	var serverCount int64
	if err := database.DB.Model(&model.MCPServer{}).Count(&serverCount).Error; err != nil {
		return nil, err
	}
	stats.Servers = serverCount

	var servers []model.MCPServer
	if err := database.DB.Model(&model.MCPServer{}).Select("tools").Find(&servers).Error; err != nil {
		return nil, err
	}
	var toolCount int64
	for _, s := range servers {
		toolCount += countTools(s.Tools)
	}
	stats.Tools = toolCount

	if database.AuditStore != nil {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 0, 1)
		todayCalls, err := database.AuditStore.CountRange(context.Background(), start, end)
		if err != nil {
			return nil, err
		}
		stats.TodayCalls = todayCalls
	}

	return stats, nil
}

type CallTrendOutput struct {
	Dates  []string `json:"dates"`
	Counts []int64  `json:"counts"`
}

func GetCallTrend(days int) (*CallTrendOutput, error) {
	out := &CallTrendOutput{Dates: []string{}, Counts: []int64{}}
	if days <= 0 || days > 365 {
		days = 30
	}

	now := time.Now()

	countByDay := map[string]int64{}
	if database.AuditStore != nil {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -(days - 1))
		dayCounts, err := database.AuditStore.CountByDay(context.Background(), start, now)
		if err != nil {
			return nil, err
		}
		for _, dc := range dayCounts {
			countByDay[dc.Day] = dc.Count
		}
	}

	dates := make([]string, 0, days)
	counts := make([]int64, 0, days)
	for i := days - 1; i >= 0; i-- {
		d := now.AddDate(0, 0, -i)
		dates = append(dates, d.Format("01-02"))
		counts = append(counts, countByDay[d.Format("2006-01-02")])
	}
	out.Dates = dates
	out.Counts = counts

	return out, nil
}

func countTools(raw string) int64 {
	if raw == "" {
		return 0
	}
	var tools []json.RawMessage
	if err := json.Unmarshal([]byte(raw), &tools); err != nil {
		return 0
	}
	return int64(len(tools))
}
