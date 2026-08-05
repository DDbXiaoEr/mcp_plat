package service

import (
	"log"

	"mcp_plat-console/database"
	"mcp_plat-console/model"
)

type AuditLogQuery struct {
	AccessKey string `form:"access_key"`
	ServerID  string `form:"server_id"`
	ToolName  string `form:"tool_name"`
	Start     string `form:"start"`
	End       string `form:"end"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

type AuditLogOutput struct {
	List     []model.AuditLog `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

func ListAuditLogs(userID uint, query AuditLogQuery) (*AuditLogOutput, error) {
	log.Printf("audit_log: query user_id=%d access_key=%s server_id=%s tool_name=%s start=%s end=%s page=%d page_size=%d",
		userID, query.AccessKey, query.ServerID, query.ToolName, query.Start, query.End, query.Page, query.PageSize)

	if database.AuditLogDB == nil {
		log.Printf("audit_log: AuditLogDB is nil, returning empty")
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

	var userAccessKeys []string
	database.DB.Model(&model.AccessKey{}).Where("user_id = ?", userID).Pluck("access_key", &userAccessKeys)

	db := database.AuditLogDB.Model(&model.AuditLog{})

	if query.AccessKey != "" {
		belongs := false
		for _, ak := range userAccessKeys {
			if ak == query.AccessKey {
				belongs = true
				break
			}
		}
		if !belongs {
			db = db.Where("1 = 0")
		} else {
			db = db.Where("access_key = ?", query.AccessKey)
		}
	} else if len(userAccessKeys) > 0 {
		db = db.Where("access_key IN ?", userAccessKeys)
	} else {
		db = db.Where("user_id = ?", userID)
	}
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
		log.Printf("audit_log: query error: %v", err)
		return nil, err
	}

	log.Printf("audit_log: result total=%d returned=%d", total, len(list))

	return &AuditLogOutput{
		List:     list,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}
