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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"mcp_plat-console/model"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type accessKeyClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := "config.yaml"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		configPath = p
	}

	type dbConfig struct {
		Type     string `yaml:"type"`
		Postgres struct {
			Host, Port, User, Password, DBName string
		} `yaml:"postgres"`
		MySQL struct {
			Host, Port, User, Password, DBName string
		} `yaml:"mysql"`
	}
	type config struct {
		AccessKeySecret string   `yaml:"access_key_secret"`
		Database        dbConfig `yaml:"database"`
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("无法读取配置文件 %s: %w", configPath, err)
	}
	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	var dial gorm.Dialector
	switch cfg.Database.Type {
	case "postgres":
		pg := cfg.Database.Postgres
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			pg.Host, pg.User, pg.Password, pg.DBName, pg.Port)
		dial = postgres.Open(dsn)
	case "mysql":
		my := cfg.Database.MySQL
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			my.User, my.Password, my.Host, my.Port, my.DBName)
		dial = mysql.Open(dsn)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", cfg.Database.Type)
	}

	db, err := gorm.Open(dial, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	if len(os.Args) < 2 {
		fmt.Println("用法: accesskey-test <access_key> [tool_name]")
		fmt.Println("示例: accesskey-test ak-eyJhbGci... query_network_info")
		os.Exit(0)
	}

	rawKey := strings.TrimSpace(os.Args[1])
	if rawKey == "" {
		return fmt.Errorf("请输入 access key")
	}

	testTool := ""
	if len(os.Args) >= 3 {
		testTool = strings.TrimSpace(os.Args[2])
	}

	if !strings.HasPrefix(rawKey, "ak-") {
		fmt.Println("[WARN] access key 不以 'ak-' 开头，将按原始字符串查询")
	}

	jwtPart := rawKey
	if strings.HasPrefix(rawKey, "ak-") {
		jwtPart = rawKey[3:]
	}

	fmt.Println("========================================")
	fmt.Println("   AccessKey 诊断工具")
	fmt.Println("========================================")
	fmt.Println()

	fmt.Printf("● 输入 Key: %s\n", rawKey)
	fmt.Println()

	decodeJWT(jwtPart)
	fmt.Println()

	var ak model.AccessKey
	if err := db.Where("key = ?", rawKey).First(&ak).Error; err != nil {
		return fmt.Errorf("数据库中未找到此 AccessKey: %w", err)
	}

	showKeyInfo(ak)
	fmt.Println()

	serverTools := inspectServers(ak.Servers)

	if testTool != "" && serverTools != nil {
		fmt.Println()
		validateToolCall(serverTools, testTool)
	}

	return nil
}

func decodeJWT(jwtStr string) {
	parts := strings.Split(jwtStr, ".")
	if len(parts) != 3 {
		fmt.Println("[WARN] JWT 格式异常（段数不等于3）")
		return
	}

	headerJSON, err := base64URLDecode(parts[0])
	fmt.Printf("● JWT Header : %s", headerJSON)
	if err != nil {
		fmt.Printf("  (解码失败: %v)", err)
	}
	fmt.Println()

	payloadJSON, err := base64URLDecode(parts[1])
	if err != nil {
		fmt.Printf("● JWT Payload: 解码失败: %v\n", err)
		return
	}

	var claims accessKeyClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		fmt.Printf("● JWT Payload: %s (无法解析为 claims: %v)\n", string(payloadJSON), err)
		return
	}

	fmt.Printf("● JWT Payload: %s\n", string(payloadJSON))
	fmt.Printf("  - user_id: %d\n", claims.UserID)
	fmt.Printf("  - role   : %s\n", claims.Role)

	var raw struct {
		Exp int64 `json:"exp"`
		Iat int64 `json:"iat"`
	}
	if err := json.Unmarshal(payloadJSON, &raw); err == nil {
		if raw.Exp > 0 {
			expired := time.Now().After(time.Unix(raw.Exp, 0))
			status := "有效"
			if expired {
				status = "已过期"
			}
			fmt.Printf("  - exp    : %d (%s, %s)\n", raw.Exp, time.Unix(raw.Exp, 0).Format(time.RFC3339), status)
		}
		if raw.Iat > 0 {
			fmt.Printf("  - iat    : %d (%s)\n", raw.Iat, time.Unix(raw.Iat, 0).Format(time.RFC3339))
		}
	}
}

func showKeyInfo(ak model.AccessKey) {
	fmt.Printf("● 数据库记录:\n")
	fmt.Printf("  - ID      : %d\n", ak.ID)
	fmt.Printf("  - UserID  : %d\n", ak.UserID)
	fmt.Printf("  - Name    : %s\n", ak.Name)
	fmt.Printf("  - Enabled : %v\n", ak.Enabled)
	if ak.ExpiredAt != nil {
		expired := time.Now().After(*ak.ExpiredAt)
		status := "有效"
		if expired {
			status = "已过期"
		}
		fmt.Printf("  - 过期时间: %s (%s)\n", ak.ExpiredAt.Format(time.RFC3339), status)
	} else {
		fmt.Println("  - 过期时间: 永不过期")
	}
	fmt.Printf("  - 创建时间: %s\n", ak.CreatedAt.Format(time.RFC3339))
	fmt.Printf("  - 更新时间: %s\n", ak.UpdatedAt.Format(time.RFC3339))
	fmt.Printf("  - Servers 长度: %d 字符\n", len(ak.Servers))
	if ak.Servers == "" {
		fmt.Println("  - Servers 内容: (空)")
		fmt.Println("  - 说明: Servers 为空, auth-server 将只执行角色级权限检查")
	}
}

func inspectServers(raw string) map[string][]string {
	fmt.Printf("● Servers 原始内容:\n")
	if raw == "" {
		fmt.Println("  (为空，无内容)")
		return nil
	}

	const maxPreview = 500
	preview := raw
	if len(preview) > maxPreview {
		preview = preview[:maxPreview] + " ...(截断)"
	}
	fmt.Printf("  %s\n\n", preview)

	fmt.Println("● 按 map 格式解析 (map[string][]string):")
	var asMap map[string][]string
	if err := json.Unmarshal([]byte(raw), &asMap); err != nil {
		fmt.Printf("  [失败] %v\n", err)
	} else {
		fmt.Printf("  [成功] %d 个服务器\n", len(asMap))
		for sid, tools := range asMap {
			if len(tools) == 0 {
				fmt.Printf("    - %s (空列表 → 禁止所有工具)\n", sid)
			} else {
				fmt.Printf("    - %s -> [%s]\n", sid, strings.Join(tools, ", "))
			}
		}
	}

	return asMap
}

func validateToolCall(serverTools map[string][]string, toolName string) {
	fmt.Printf("● 模拟 auth-server 工具校验:\n")
	fmt.Printf("  - 目标工具: %s\n", toolName)

	if len(serverTools) == 0 {
		fmt.Println("  - 结果: [禁止] 该 AccessKey 未授权任何服务器（Servers 配置为空）")
		return
	}

	for sid, allowedTools := range serverTools {
		if len(allowedTools) == 0 {
			fmt.Printf("  - 服务器 %s: [禁止] 工具列表为空，禁止所有工具\n", sid)
			continue
		}

		matched := false
		for _, t := range allowedTools {
			if t == toolName {
				matched = true
				break
			}
		}
		if matched {
			fmt.Printf("  - 服务器 %s: [允许] 工具 %s 在授权列表中\n", sid, toolName)
		} else {
			fmt.Printf("  - 服务器 %s: [禁止] 工具 %s 不在授权列表 %v 中\n", sid, toolName, allowedTools)
		}
	}

	fmt.Println()
	fmt.Println("  如果运行时未被拦截，请检查:")
	fmt.Println("  1. auth-server 是否已用最新代码编译部署")
	fmt.Println("  2. 请求的工具名是否与数据库中的完全一致（注意大小写/下划线）")
	fmt.Println("  3. X-Access-Key 是否在请求头中正确传递")
	fmt.Println("  4. APISIX 插件 accesskey_verify 是否正确配置")
}

func base64URLDecode(s string) ([]byte, error) {
	pad := len(s) % 4
	if pad > 0 {
		s += strings.Repeat("=", 4-pad)
	}
	return base64.URLEncoding.DecodeString(s)
}
