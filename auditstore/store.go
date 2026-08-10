package auditstore

// Author: deepseek-v4-pro / opencode

import (
	"context"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/model"
)

// Query 使用历史查询条件
type Query struct {
	UserID       uint
	AccessKeyIDs []uint // 用户拥有的 AccessKey ID 集合
	AccessKeyID  uint   // 指定 AccessKey ID（优先于 AccessKeyIDs）
	ServerID     string
	ToolName     string
	Start        time.Time
	End          time.Time
	Page         int
	PageSize     int
}

// DayCount 按天调用数
type DayCount struct {
	Day   string // YYYY-MM-DD
	Count int64
}

// Store 审计日志存储抽象，支持关系库（Postgres/SQLite）与 ClickHouse 两种实现
type Store interface {
	// Append 批量写入审计日志
	Append(ctx context.Context, logs []model.AuditLog) error
	// List 分页查询，返回列表与总数
	List(ctx context.Context, q Query) ([]model.AuditLog, int64, error)
	// CountRange 统计 [start, end) 区间内的记录数
	CountRange(ctx context.Context, start, end time.Time) (int64, error)
	// CountByDay 按天分组统计 [start, end) 区间内的记录数
	CountByDay(ctx context.Context, start, end time.Time) ([]DayCount, error)
}

// Purgable 支持按时间清理旧数据（ClickHouse 用原生 TTL，无需实现）
type Purgable interface {
	Purge(ctx context.Context, before time.Time) (int64, error)
}

// New 按配置创建存储实现
func New(cfg config.DatabaseConfig) (Store, error) {
	switch cfg.Type {
	case "clickhouse":
		return newClickHouseStore(cfg.ClickHouse)
	default:
		return newGORMStore(cfg)
	}
}
