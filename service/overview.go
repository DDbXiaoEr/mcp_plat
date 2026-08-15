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
	"sort"
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

type NameCount struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type CallTrendOutput struct {
	Dates             []string    `json:"dates"`
	Counts            []int64     `json:"counts"`
	ServerCalls       []NameCount `json:"server_calls"`
	UserGroups        []NameCount `json:"user_groups"`
	ServerNames       []string    `json:"server_names"`
	ServerCountsByDay [][]int64   `json:"server_counts_by_day"`
}

// topServerN 服务器调用量堆叠展示时，独立成系列的最大服务器数，其余合并为「其他」
const topServerN = 8

func GetCallTrend(days int) (*CallTrendOutput, error) {
	out := &CallTrendOutput{
		Dates:             []string{},
		Counts:            []int64{},
		ServerCalls:       []NameCount{},
		UserGroups:        []NameCount{},
		ServerNames:       []string{},
		ServerCountsByDay: [][]int64{},
	}
	if days <= 0 || days > 365 {
		days = 30
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -(days - 1))

	countByDay := map[string]int64{}
	if database.AuditStore != nil {
		dayCounts, err := database.AuditStore.CountByDay(context.Background(), start, now)
		if err != nil {
			return nil, err
		}
		for _, dc := range dayCounts {
			countByDay[dc.Day] = dc.Count
		}
	}

	dateTimes := make([]time.Time, 0, days)
	dates := make([]string, 0, days)
	counts := make([]int64, 0, days)
	for i := days - 1; i >= 0; i-- {
		d := now.AddDate(0, 0, -i)
		dateTimes = append(dateTimes, d)
		dates = append(dates, d.Format("01-02"))
		counts = append(counts, countByDay[d.Format("2006-01-02")])
	}
	out.Dates = dates
	out.Counts = counts

	if database.AuditStore != nil {
		nameByID, err := serverNameByID()
		if err != nil {
			return nil, err
		}
		if err := fillServerCalls(out, start, now, nameByID); err != nil {
			return nil, err
		}
		if err := fillServerCountsByDay(out, start, now, dateTimes, nameByID); err != nil {
			return nil, err
		}
	}
	if err := fillUserGroups(out); err != nil {
		return nil, err
	}

	return out, nil
}

func serverNameByID() (map[string]string, error) {
	var servers []model.MCPServer
	if err := database.DB.Model(&model.MCPServer{}).Select("id", "name").Find(&servers).Error; err != nil {
		return nil, err
	}
	nameByID := make(map[string]string, len(servers))
	for _, s := range servers {
		nameByID[s.ID] = s.Name
	}
	return nameByID, nil
}

func fillServerCalls(out *CallTrendOutput, start, end time.Time, nameByID map[string]string) error {
	if database.AuditStore == nil {
		return nil
	}
	serverCounts, err := database.AuditStore.CountByServer(context.Background(), start, end)
	if err != nil {
		return err
	}
	if len(serverCounts) == 0 {
		return nil
	}

	out.ServerCalls = make([]NameCount, 0, len(serverCounts))
	for _, sc := range serverCounts {
		name, ok := nameByID[sc.ServerID]
		if !ok || name == "" {
			name = "未知服务器"
		}
		out.ServerCalls = append(out.ServerCalls, NameCount{Name: name, Count: sc.Count})
	}
	sort.Slice(out.ServerCalls, func(i, j int) bool {
		return out.ServerCalls[i].Count > out.ServerCalls[j].Count
	})
	return nil
}

func fillServerCountsByDay(out *CallTrendOutput, start, end time.Time, dateTimes []time.Time, nameByID map[string]string) error {
	if database.AuditStore == nil {
		return nil
	}
	dayServerCounts, err := database.AuditStore.CountByDayServer(context.Background(), start, end)
	if err != nil {
		return err
	}

	dayTotals := map[string]map[string]int64{}
	serverTotals := map[string]int64{}
	for _, dc := range dayServerCounts {
		if dayTotals[dc.Day] == nil {
			dayTotals[dc.Day] = map[string]int64{}
		}
		dayTotals[dc.Day][dc.ServerID] += dc.Count
		serverTotals[dc.ServerID] += dc.Count
	}
	if len(serverTotals) == 0 {
		return nil
	}

	serverIDs := make([]string, 0, len(serverTotals))
	for id := range serverTotals {
		serverIDs = append(serverIDs, id)
	}
	sort.Slice(serverIDs, func(i, j int) bool {
		if serverTotals[serverIDs[i]] != serverTotals[serverIDs[j]] {
			return serverTotals[serverIDs[i]] > serverTotals[serverIDs[j]]
		}
		return serverIDs[i] < serverIDs[j]
	})

	var top []string
	var otherIDs []string
	if len(serverIDs) > topServerN {
		top = serverIDs[:topServerN]
		otherIDs = serverIDs[topServerN:]
	} else {
		top = serverIDs
	}

	seriesIDs := make([]string, 0, len(top)+1)
	serverNames := make([]string, 0, len(top)+1)
	for _, id := range top {
		name, ok := nameByID[id]
		if !ok || name == "" {
			name = "未知服务器"
		}
		seriesIDs = append(seriesIDs, id)
		serverNames = append(serverNames, name)
	}
	if len(otherIDs) > 0 {
		seriesIDs = append(seriesIDs, "")
		serverNames = append(serverNames, "其他")
	}
	out.ServerNames = serverNames

	matrix := make([][]int64, len(dateTimes))
	for i, d := range dateTimes {
		key := d.Format("2006-01-02")
		dayMap := dayTotals[key]
		row := make([]int64, len(seriesIDs))
		for j, id := range seriesIDs {
			if id == "" {
				for _, oid := range otherIDs {
					row[j] += dayMap[oid]
				}
			} else {
				row[j] = dayMap[id]
			}
		}
		matrix[i] = row
	}
	out.ServerCountsByDay = matrix
	return nil
}

func fillUserGroups(out *CallTrendOutput) error {
	type row struct {
		RoleName string
		Count    int64
	}
	var rows []row
	err := database.DB.Model(&model.User{}).
		Select("COALESCE(r.name, '未分组') AS role_name, COUNT(*) AS count").
		Joins("LEFT JOIN roles r ON r.id = users.role_id").
		Group("role_name").
		Scan(&rows).Error
	if err != nil {
		return err
	}

	out.UserGroups = make([]NameCount, 0, len(rows))
	for _, r := range rows {
		out.UserGroups = append(out.UserGroups, NameCount{Name: r.RoleName, Count: r.Count})
	}
	sort.Slice(out.UserGroups, func(i, j int) bool {
		return out.UserGroups[i].Count > out.UserGroups[j].Count
	})
	return nil
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
