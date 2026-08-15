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

package main

// Author: deepseek-v4-pro / opencode

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/model"
)

var (
	days   int
	perDay int
	dryRun bool
)

func init() {
	flag.IntVar(&days, "days", 30, "生成最近 N 天的审计数据")
	flag.IntVar(&perDay, "per-day", 200, "每天目标生成条数（会做 ±25% 抖动）")
	flag.BoolVar(&dryRun, "dry-run", false, "仅预览不写入")
}

type serverRef struct {
	ID    string
	Name  string
	Tools []string
}

type keyRef struct {
	ID     uint
	UserID uint
	Name   string
}

var failMessages = []string{
	"connection refused",
	"timeout: deadline exceeded",
	"tool not found",
	"invalid arguments",
	"server internal error",
}

var fallbackTools = []string{"query", "search", "get", "list", "create", "update", "delete", "generate"}

func main() {
	flag.Parse()

	config.Load()
	database.Init()

	if database.AuditStore == nil {
		fmt.Println("audit_log_db 未配置，无法生成审计数据")
		os.Exit(1)
	}

	servers := loadServers()
	users := loadUsers()
	keys := loadKeys()

	if len(servers) == 0 {
		fmt.Println("mcp_servers 表为空，请先运行 datagen 生成服务器数据")
		os.Exit(1)
	}
	if len(users) == 0 {
		fmt.Println("users 表为空，请先运行 datagen 生成用户数据")
		os.Exit(1)
	}

	now := time.Now()
	var all []model.AuditLog
	total := 0
	for d := days - 1; d >= 0; d-- {
		day := now.AddDate(0, 0, -d)
		n := perDayWithJitter(perDay)
		for i := 0; i < n; i++ {
			all = append(all, buildAuditLog(day, servers, users, keys, now))
		}
		total += n
	}

	if dryRun {
		fmt.Printf("dry-run: 将生成 %d 条审计日志（最近 %d 天，每天约 %d 条）\n", total, days, perDay)
		printSample(all)
		return
	}

	if err := writeAll(all); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("完成: 生成并写入 %d 条审计日志（存储类型=%s）\n", total, config.AppConfig.AuditLogDB.Type)
}

func loadServers() []serverRef {
	var servers []model.MCPServer
	if err := database.DB.Find(&servers).Error; err != nil {
		fmt.Printf("查询 mcp_servers 失败: %v\n", err)
		os.Exit(1)
	}
	refs := make([]serverRef, 0, len(servers))
	for _, s := range servers {
		refs = append(refs, serverRef{ID: s.ID, Name: s.Name, Tools: parseToolNames(s.Tools)})
	}
	return refs
}

func parseToolNames(raw string) []string {
	if raw == "" {
		return nil
	}
	var items []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil
	}
	names := make([]string, 0, len(items))
	for _, it := range items {
		if it.Name != "" {
			names = append(names, it.Name)
		}
	}
	return names
}

func loadUsers() []uint {
	var ids []uint
	if err := database.DB.Model(&model.User{}).Pluck("id", &ids).Error; err != nil {
		fmt.Printf("查询 users 失败: %v\n", err)
		os.Exit(1)
	}
	return ids
}

func loadKeys() []keyRef {
	var keys []model.AccessKey
	if err := database.DB.Model(&model.AccessKey{}).Select("id", "user_id", "name").Find(&keys).Error; err != nil {
		fmt.Printf("查询 access_keys 失败: %v\n", err)
		os.Exit(1)
	}
	refs := make([]keyRef, 0, len(keys))
	for _, k := range keys {
		refs = append(refs, keyRef{ID: k.ID, UserID: k.UserID, Name: k.Name})
	}
	return refs
}

func buildAuditLog(day time.Time, servers []serverRef, users []uint, keys []keyRef, now time.Time) model.AuditLog {
	s := servers[randomInt(len(servers))]

	var userID, keyID uint
	var keyName string
	if len(keys) > 0 {
		k := keys[randomInt(len(keys))]
		userID, keyID, keyName = k.UserID, k.ID, k.Name
	} else {
		userID = users[randomInt(len(users))]
	}

	success := randomInt(100) < 90
	msg := "ok"
	if !success {
		msg = failMessages[randomInt(len(failMessages))]
	}

	ts := day.Add(time.Duration(randomHour())*time.Hour +
		time.Duration(randomInt(60))*time.Minute +
		time.Duration(randomInt(60))*time.Second)
	if ts.After(now) {
		ts = now
	}

	return model.AuditLog{
		AccessKeyID:   keyID,
		UserID:        userID,
		ServerID:      s.ID,
		ToolName:      pickTool(s),
		Success:       success,
		Message:       msg,
		CreatedAt:     ts,
		AccessKeyName: keyName,
	}
}

func pickTool(s serverRef) string {
	if len(s.Tools) > 0 {
		return s.Tools[randomInt(len(s.Tools))]
	}
	return fallbackTools[randomInt(len(fallbackTools))]
}

func writeAll(logs []model.AuditLog) error {
	const batch = 500
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	for i := 0; i < len(logs); i += batch {
		end := i + batch
		if end > len(logs) {
			end = len(logs)
		}
		if err := database.AuditStore.Append(ctx, logs[i:end]); err != nil {
			return err
		}
		fmt.Printf("已写入 %d/%d 条\n", end, len(logs))
	}
	return nil
}

func printSample(logs []model.AuditLog) {
	n := len(logs)
	if n > 5 {
		n = 5
	}
	for i := 0; i < n; i++ {
		l := logs[i]
		fmt.Printf("  样例[%d] time=%s server=%s tool=%s user_id=%d access_key_id=%d success=%v\n",
			i, l.CreatedAt.Format("2006-01-02 15:04:05"), l.ServerID, l.ToolName, l.UserID, l.AccessKeyID, l.Success)
	}
}

func randomInt(max int) int {
	if max <= 0 {
		return 0
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func perDayWithJitter(base int) int {
	if base <= 0 {
		base = 200
	}
	n := base + randomInt(base/2+1) - base/4
	if n < 10 {
		n = 10
	}
	return n
}

func randomHour() int {
	if randomInt(100) < 75 {
		return 8 + randomInt(13) // 8..20，工作日高峰
	}
	return randomInt(24)
}
