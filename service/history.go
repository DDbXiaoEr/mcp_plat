package service

// Author: deepseek-v4-pro / opencode

import (
	"mcp_plat-console/database"
	"mcp_plat-console/model"
)

type HistoryQuery struct {
	Server      string `form:"server"`
	AccessKeyID uint   `form:"access_key_id"`
	StartDate   string `form:"start_date"`
	EndDate     string `form:"end_date"`
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
}

type HistoryListOutput struct {
	List     []model.UsageHistory `json:"list"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

func ListHistory(userID uint, query HistoryQuery) (*HistoryListOutput, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 20
	}

	db := database.DB.Model(&model.UsageHistory{}).Where("user_id = ?", userID)

	if query.Server != "" {
		db = db.Where("server = ?", query.Server)
	}
	if query.AccessKeyID > 0 {
		db = db.Where("access_key_id = ?", query.AccessKeyID)
	}
	if query.StartDate != "" {
		db = db.Where("created_at >= ?", query.StartDate)
	}
	if query.EndDate != "" {
		db = db.Where("created_at <= ?", query.EndDate+" 23:59:59")
	}

	var total int64
	db.Count(&total)

	var list []model.UsageHistory
	offset := (query.Page - 1) * query.PageSize
	err := db.Order("created_at desc").Offset(offset).Limit(query.PageSize).Find(&list).Error
	if err != nil {
		return nil, err
	}

	return &HistoryListOutput{
		List:     list,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}
