package auditstore

// Author: deepseek-v4-pro / opencode

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/model"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const (
	ckDefaultAddr  = "127.0.0.1:9000"
	ckDefaultDB    = "default"
	ckDefaultTable = "audit_logs"
	ckDefaultTTL   = 90
)

type clickHouseStore struct {
	conn    driver.Conn
	db      string
	table   string
	ttlDays int
	idSeq   uint64
}

func newClickHouseStore(cfg config.ClickHouseConfig) (Store, error) {
	if cfg.Addr == "" {
		cfg.Addr = ckDefaultAddr
	}
	if cfg.DB == "" {
		cfg.DB = ckDefaultDB
	}
	if cfg.Table == "" {
		cfg.Table = ckDefaultTable
	}
	if cfg.TTLDays <= 0 {
		cfg.TTLDays = ckDefaultTTL
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.Addr},
		Auth: clickhouse.Auth{
			Database: cfg.DB,
			Username: cfg.User,
			Password: cfg.Password,
		},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("clickhouse 连接失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("clickhouse ping 失败: %w", err)
	}

	s := &clickHouseStore{
		conn:    conn,
		db:      cfg.DB,
		table:   cfg.Table,
		ttlDays: cfg.TTLDays,
	}
	if err := s.ensureTable(ctx, cfg.Engine, cfg.Cluster); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *clickHouseStore) tableName() string {
	return fmt.Sprintf("%s.%s", s.db, s.table)
}

func (s *clickHouseStore) ensureTable(ctx context.Context, engine, cluster string) error {
	onCluster := ""
	if cluster != "" {
		onCluster = fmt.Sprintf(" ON CLUSTER %s", cluster)
	}

	tableEngine := "MergeTree"
	if strings.EqualFold(engine, "replicated_merge_tree") {
		if cluster == "" {
			return errors.New("clickhouse engine=replicated_merge_tree 时必须配置 cluster")
		}
		tableEngine = fmt.Sprintf("ReplicatedMergeTree('/clickhouse/tables/{shard}/%s', '{replica}')", s.tableName())
	}

	ddl := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s%s (
  id UInt64,
  access_key_id UInt64,
  user_id UInt64,
  server_id String,
  tool_name String,
  success UInt8,
  message String,
  created_at DateTime64(3)
) ENGINE = %s
PARTITION BY toYYYYMM(created_at)
ORDER BY (created_at, id)
TTL toDateTime(created_at) + INTERVAL %d DAY`, s.tableName(), onCluster, tableEngine, s.ttlDays)

	return s.conn.Exec(ctx, ddl)
}

func (s *clickHouseStore) nextID(t time.Time) uint64 {
	seq := atomic.AddUint64(&s.idSeq, 1)
	return uint64(t.UnixNano()) + (seq % 4096)
}

func (s *clickHouseStore) Append(ctx context.Context, logs []model.AuditLog) error {
	if len(logs) == 0 {
		return nil
	}

	batch, err := s.conn.PrepareBatch(ctx, fmt.Sprintf(
		"INSERT INTO %s (id, access_key_id, user_id, server_id, tool_name, success, message, created_at)",
		s.tableName()))
	if err != nil {
		return err
	}

	for _, log := range logs {
		success := uint8(0)
		if log.Success {
			success = 1
		}
		if err := batch.Append(
			s.nextID(log.CreatedAt),
			log.AccessKeyID,
			log.UserID,
			log.ServerID,
			log.ToolName,
			success,
			log.Message,
			log.CreatedAt,
		); err != nil {
			return err
		}
	}
	return batch.Send()
}

func (s *clickHouseStore) List(ctx context.Context, q Query) ([]model.AuditLog, int64, error) {
	where, args := buildWhere(q)

	countQuery := fmt.Sprintf("SELECT count() FROM %s%s", s.tableName(), where)
	var total uint64
	if err := s.conn.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page, pageSize := normalizePage(q.Page, q.PageSize)
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(
		"SELECT id, access_key_id, user_id, server_id, tool_name, success, message, created_at FROM %s%s ORDER BY created_at DESC, id DESC LIMIT %d OFFSET %d",
		s.tableName(), where, pageSize, offset)

	rows, err := s.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []model.AuditLog{}
	for rows.Next() {
		var (
			log       model.AuditLog
			successU8 uint8
			createdAt time.Time
		)
		if err := rows.Scan(
			&log.ID,
			&log.AccessKeyID,
			&log.UserID,
			&log.ServerID,
			&log.ToolName,
			&successU8,
			&log.Message,
			&createdAt,
		); err != nil {
			return nil, 0, err
		}
		log.Success = successU8 != 0
		log.CreatedAt = createdAt
		list = append(list, log)
	}
	return list, int64(total), rows.Err()
}

func (s *clickHouseStore) CountRange(ctx context.Context, start, end time.Time) (int64, error) {
	var count uint64
	err := s.conn.QueryRow(ctx,
		fmt.Sprintf("SELECT count() FROM %s WHERE created_at >= ? AND created_at < ?", s.tableName()),
		start, end).Scan(&count)
	return int64(count), err
}

func (s *clickHouseStore) CountByDay(ctx context.Context, start, end time.Time) ([]DayCount, error) {
	rows, err := s.conn.Query(ctx,
		fmt.Sprintf("SELECT toDate(created_at) AS day, count() FROM %s WHERE created_at >= ? AND created_at < ? GROUP BY day ORDER BY day", s.tableName()),
		start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []DayCount{}
	for rows.Next() {
		var (
			day   time.Time
			count uint64
		)
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		out = append(out, DayCount{Day: day.Format("2006-01-02"), Count: int64(count)})
	}
	return out, rows.Err()
}

func buildWhere(q Query) (string, []interface{}) {
	var conds []string
	var args []interface{}

	if q.AccessKeyID > 0 {
		conds = append(conds, "access_key_id = ?")
		args = append(args, q.AccessKeyID)
	} else if len(q.AccessKeyIDs) > 0 {
		ph := strings.TrimSuffix(strings.Repeat("?,", len(q.AccessKeyIDs)), ",")
		conds = append(conds, "access_key_id IN ("+ph+")")
		for _, id := range q.AccessKeyIDs {
			args = append(args, id)
		}
	} else if q.UserID > 0 {
		conds = append(conds, "user_id = ?")
		args = append(args, q.UserID)
	}
	if q.ServerID != "" {
		conds = append(conds, "server_id = ?")
		args = append(args, q.ServerID)
	}
	if q.ToolName != "" {
		conds = append(conds, "tool_name = ?")
		args = append(args, q.ToolName)
	}
	if !q.Start.IsZero() {
		conds = append(conds, "created_at >= ?")
		args = append(args, q.Start)
	}
	if !q.End.IsZero() {
		conds = append(conds, "created_at <= ?")
		args = append(args, q.End)
	}

	if len(conds) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}
