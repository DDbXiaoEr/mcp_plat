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

package auditstore

// Author: deepseek-v4-pro / opencode

import (
	"context"
	"fmt"
	"log"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type gormStore struct {
	db *gorm.DB
}

// newGORMStore 打开关系库并迁移建表
func newGORMStore(cfg config.DatabaseConfig) (Store, error) {
	db := openGORM(cfg)
	if err := db.AutoMigrate(&model.AuditLog{}); err != nil {
		return nil, err
	}
	return &gormStore{db: db}, nil
}

// NewGORM 包装已打开并迁移好的 GORM 连接（主控制台使用）
func NewGORM(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func openGORM(cfg config.DatabaseConfig) *gorm.DB {
	var dialector gorm.Dialector

	switch cfg.Type {
	case "postgres":
		pg := cfg.Postgres
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			pg.Host, pg.User, pg.Password, pg.DBName, pg.Port)
		dialector = postgres.Open(dsn)
	case "mysql":
		my := cfg.MySQL
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			my.User, my.Password, my.Host, my.Port, my.DBName)
		dialector = mysql.Open(dsn)
	default:
		log.Fatalf("unsupported audit database type: %s", cfg.Type)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect audit database: %v", err)
	}
	return db
}

func (s *gormStore) Append(ctx context.Context, logs []model.AuditLog) error {
	if len(logs) == 0 {
		return nil
	}
	return s.db.WithContext(ctx).CreateInBatches(logs, len(logs)).Error
}

func (s *gormStore) List(ctx context.Context, q Query) ([]model.AuditLog, int64, error) {
	db := s.db.WithContext(ctx).Model(&model.AuditLog{})

	if q.AccessKeyID > 0 {
		if q.UserID > 0 {
			db = db.Where("(access_key_id = ? OR access_key_id = 0) AND user_id = ?", q.AccessKeyID, q.UserID)
		} else {
			db = db.Where("access_key_id = ?", q.AccessKeyID)
		}
	} else if len(q.AccessKeyIDs) > 0 {
		keyIDs := q.AccessKeyIDs
		if q.UserID > 0 {
			db = db.Where("(access_key_id IN ? OR access_key_id = 0) AND user_id = ?", keyIDs, q.UserID)
		} else {
			db = db.Where("access_key_id IN ?", keyIDs)
		}
	} else if q.UserID > 0 {
		db = db.Where("user_id = ?", q.UserID)
	}
	if q.ServerID != "" {
		db = db.Where("server_id = ?", q.ServerID)
	}
	if q.ToolName != "" {
		db = db.Where("tool_name = ?", q.ToolName)
	}
	if !q.Start.IsZero() {
		db = db.Where("created_at >= ?", q.Start)
	}
	if !q.End.IsZero() {
		db = db.Where("created_at <= ?", q.End)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page, pageSize := normalizePage(q.Page, q.PageSize)
	var list []model.AuditLog
	offset := (page - 1) * pageSize
	err := db.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *gormStore) CountRange(ctx context.Context, start, end time.Time) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&model.AuditLog{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&count).Error
	return count, err
}

func (s *gormStore) CountByDay(ctx context.Context, start, end time.Time) ([]DayCount, error) {
	type row struct {
		Day   string
		Count int64
	}
	var rows []row
	err := s.db.WithContext(ctx).Model(&model.AuditLog{}).
		Select("DATE(created_at) AS day, COUNT(*) AS count").
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("DATE(created_at)").
		Order("day").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]DayCount, 0, len(rows))
	for _, r := range rows {
		out = append(out, DayCount{Day: r.Day, Count: r.Count})
	}
	return out, nil
}

func (s *gormStore) CountByServer(ctx context.Context, start, end time.Time) ([]ServerCount, error) {
	type row struct {
		ServerID string
		Count    int64
	}
	var rows []row
	err := s.db.WithContext(ctx).Model(&model.AuditLog{}).
		Select("server_id, COUNT(*) AS count").
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("server_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]ServerCount, 0, len(rows))
	for _, r := range rows {
		out = append(out, ServerCount{ServerID: r.ServerID, Count: r.Count})
	}
	return out, nil
}

func (s *gormStore) CountByDayServer(ctx context.Context, start, end time.Time) ([]DayServerCount, error) {
	type row struct {
		Day      string
		ServerID string
		Count    int64
	}
	var rows []row
	err := s.db.WithContext(ctx).Model(&model.AuditLog{}).
		Select("DATE(created_at) AS day, server_id, COUNT(*) AS count").
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("DATE(created_at), server_id").
		Order("day, server_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]DayServerCount, 0, len(rows))
	for _, r := range rows {
		out = append(out, DayServerCount{Day: r.Day, ServerID: r.ServerID, Count: r.Count})
	}
	return out, nil
}

func (s *gormStore) Purge(ctx context.Context, before time.Time) (int64, error) {
	var deleted int64
	for {
		result := s.db.WithContext(ctx).
			Where("created_at < ?", before).
			Limit(1000).
			Delete(&model.AuditLog{})
		if result.Error != nil {
			return deleted, result.Error
		}
		if result.RowsAffected == 0 {
			break
		}
		deleted += result.RowsAffected
	}
	return deleted, nil
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}
