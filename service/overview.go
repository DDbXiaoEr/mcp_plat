package service

// Author: deepseek-v4-pro / opencode

import (
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

	if database.AuditLogDB != nil {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 0, 1)
		var todayCalls int64
		if err := database.AuditLogDB.Model(&model.AuditLog{}).
			Where("created_at >= ? AND created_at < ?", start, end).
			Count(&todayCalls).Error; err != nil {
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
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -(days - 1))

	countByDay := map[string]int64{}
	if database.AuditLogDB != nil {
		type row struct {
			Day   string
			Count int64
		}
		var rows []row
		if err := database.AuditLogDB.Model(&model.AuditLog{}).
			Select("DATE(created_at) AS day, COUNT(*) AS count").
			Where("created_at >= ?", start).
			Group("DATE(created_at)").
			Scan(&rows).Error; err != nil {
			return nil, err
		}
		for _, r := range rows {
			countByDay[r.Day] = r.Count
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
