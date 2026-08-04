package service

import (
	"mcp_plat-console/database"
	"mcp_plat-console/model"
)

type AuditLogQuery struct {
	ServerID string `form:"server_id"`
	ToolName string `form:"tool_name"`
	Start    string `form:"start"`
	End      string `form:"end"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type AuditLogOutput struct {
	List     []model.AuditLog `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

func ListAuditLogs(userID uint, query AuditLogQuery) (*AuditLogOutput, error) {
	if database.AuditLogDB == nil {
		return &AuditLogOutput{
			List:     []model.AuditLog{},
			Total:    0,
			Page:     query.Page,
			PageSize: query.PageSize,
		}, nil
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 20
	}

	db := database.AuditLogDB.Model(&model.AuditLog{}).Where("user_id = ?", userID)

	if query.ServerID != "" {
		db = db.Where("server_id = ?", query.ServerID)
	}
	if query.ToolName != "" {
		db = db.Where("tool_name = ?", query.ToolName)
	}
	if query.Start != "" {
		db = db.Where("created_at >= ?", query.Start)
	}
	if query.End != "" {
		db = db.Where("created_at <= ?", query.End+" 23:59:59")
	}

	var total int64
	db.Count(&total)

	var list []model.AuditLog
	offset := (query.Page - 1) * query.PageSize
	err := db.Order("created_at desc").Offset(offset).Limit(query.PageSize).Find(&list).Error
	if err != nil {
		return nil, err
	}

	return &AuditLogOutput{
		List:     list,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}
