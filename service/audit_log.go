package service

// Author: deepseek-v4-pro / opencode

import (
	"context"
	"log"
	"time"

	"mcp_plat-console/auditstore"
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
	if database.AuditStore == nil {
		log.Printf("audit_log: AuditStore is nil, returning empty")
		return &AuditLogOutput{
			List:     []model.AuditLog{},
			Total:    0,
			Page:     query.Page,
			PageSize: query.PageSize,
		}, nil
	}

	var keys []model.AccessKey
	if err := database.DB.Where("user_id = ?", userID).Find(&keys).Error; err != nil {
		log.Printf("audit_log: query access keys error: %v", err)
		return nil, err
	}

	idToName := make(map[uint]string, len(keys))
	keyIDs := make([]uint, 0, len(keys))
	for _, k := range keys {
		idToName[k.ID] = k.Name
		keyIDs = append(keyIDs, k.ID)
	}

	storeQuery := auditstore.Query{
		UserID:       userID,
		AccessKeyIDs: keyIDs,
		ServerID:     query.ServerID,
		ToolName:     query.ToolName,
		Page:         query.Page,
		PageSize:     query.PageSize,
	}

	if query.Start != "" {
		if t, err := time.ParseInLocation("2006-01-02", query.Start, time.Local); err == nil {
			storeQuery.Start = t
		}
	}
	if query.End != "" {
		if t, err := time.ParseInLocation("2006-01-02", query.End, time.Local); err == nil {
			storeQuery.End = t.AddDate(0, 0, 1)
		}
	}

	if query.AccessKey != "" {
		var matched *model.AccessKey
		for i := range keys {
			if keys[i].Key == query.AccessKey {
				matched = &keys[i]
				break
			}
		}
		if matched == nil {
			return &AuditLogOutput{
				List:     []model.AuditLog{},
				Total:    0,
				Page:     query.Page,
				PageSize: query.PageSize,
			}, nil
		}
		storeQuery.AccessKeyID = matched.ID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, total, err := database.AuditStore.List(ctx, storeQuery)
	if err != nil {
		log.Printf("audit_log: query error: %v", err)
		return nil, err
	}

	for i := range list {
		list[i].AccessKeyName = idToName[list[i].AccessKeyID]
		if list[i].AccessKeyName == "" && len(keys) == 1 {
			// 旧版 access key 无 key_id 声明，记录 access_key_id=0，单 key 用户可直接归因
			list[i].AccessKeyName = keys[0].Name
		}
	}

	return &AuditLogOutput{
		List:     list,
		Total:    total,
		Page:     storeQuery.Page,
		PageSize: storeQuery.PageSize,
	}, nil
}
