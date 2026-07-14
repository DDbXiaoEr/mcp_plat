package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/model"

	"golang.org/x/crypto/bcrypt"
)

var (
	table   string
	count   int
	showHelp bool
)

func init() {
	flag.StringVar(&table, "table", "", "table name: user, access_key, usage_history, mcp_server")
	flag.IntVar(&count, "count", 1, "number of records to generate")
	flag.BoolVar(&showHelp, "h", false, "show help for the specified table")
}

func main() {
	flag.Parse()

	if showHelp {
		printHelp(table)
		os.Exit(0)
	}

	if table == "" {
		fmt.Println("Usage: datagen -table <table> [-count <n>]")
		fmt.Println("Tables: user, access_key, usage_history, mcp_server")
		fmt.Println("Use -h -table <table> for detailed help on each table")
		os.Exit(1)
	}

	config.Load()
	database.Init()

	switch table {
	case "user":
		genUsers(count)
	case "access_key":
		genAccessKeys(count)
	case "usage_history":
		genUsageHistories(count)
	case "mcp_server":
		genMCPServers(count)
	default:
		fmt.Printf("unknown table: %s\n", table)
		fmt.Println("available tables: user, access_key, usage_history, mcp_server")
		os.Exit(1)
	}
}

func genUsers(n int) {
	for i := 0; i < n; i++ {
		uid := randomHex(8)
		username := fmt.Sprintf("testuser_%s", randomHex(4))
		hashed, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		email := fmt.Sprintf("%s@example.com", username)
		phone := fmt.Sprintf("138%08d", randomInt(1e8))
		org := randomChoice([]string{"计算机学院", "信息中心", "网络中心", "图书馆", "教务处", "学生处", "后勤管理处", "理学院", "文学院"})

		u := model.User{
			UID:          uid,
			Username:     username,
			Password:     string(hashed),
			Email:        email,
			Phone:        phone,
			Organization: org,
			RoleID:       nil,
			Status:       1,
		}
		if err := database.DB.Create(&u).Error; err != nil {
			fmt.Printf("FAIL: create user %s: %v\n", username, err)
		} else {
			fmt.Printf("OK:   created user id=%d username=%s password=123456\n", u.ID, u.Username)
		}
	}
}

func genAccessKeys(n int) {
	var userIDs []uint
	database.DB.Model(&model.User{}).Pluck("id", &userIDs)
	if len(userIDs) == 0 {
		fmt.Println("No users found, create a user first")
		os.Exit(1)
	}

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("ak_%s", randomHex(4))
		key := "ak-" + randomHex(24)
		userID := randomFromSlice(userIDs)
		servers := randomChoice([]string{"server-1", "server-2", "server-3"})
		var expiredAt *time.Time
		if randomInt(2) == 0 {
			t := time.Now().Add(time.Duration(randomInt(365*24)) * time.Hour)
			expiredAt = &t
		}

		ak := model.AccessKey{
			UserID:    userID,
			Name:      name,
			Key:       key,
			Enabled:   true,
			Servers:   servers,
			ExpiredAt: expiredAt,
		}
		if err := database.DB.Create(&ak).Error; err != nil {
			fmt.Printf("FAIL: create access_key %s: %v\n", name, err)
		} else {
			fmt.Printf("OK:   created access_key id=%d name=%s key=%s user_id=%d\n", ak.ID, ak.Name, ak.Key, ak.UserID)
		}
	}
}

func genUsageHistories(n int) {
	var userIDs []uint
	database.DB.Model(&model.User{}).Pluck("id", &userIDs)

	var keyIDs []uint
	database.DB.Model(&model.AccessKey{}).Pluck("id", &keyIDs)

	for i := 0; i < n; i++ {
		userID := randomFromSlice(userIDs)
		keyID := randomFromSlice(keyIDs)
		server := randomChoice([]string{"server-1", "server-2", "server-3"})
		endpoint := randomChoice([]string{"/api/tools", "/api/resources", "/api/prompts"})
		status := randomChoice([]string{"success", "failed", "timeout"})

		h := model.UsageHistory{
			UserID:      userID,
			AccessKeyID: keyID,
			Server:      server,
			Endpoint:    endpoint,
			Status:      status,
		}
		if err := database.DB.Create(&h).Error; err != nil {
			fmt.Printf("FAIL: create usage_history: %v\n", err)
		} else {
			fmt.Printf("OK:   created usage_history id=%d user_id=%d server=%s\n", h.ID, h.UserID, h.Server)
		}
	}
}

func genMCPServers(n int) {
	depts := []string{"计算机学院", "信息中心", "网络中心", "图书馆", "教务处"}
	protocols := []string{"SSE", "streamable-http"}
	toolsList := []string{
		"[]",
		`[{"name":"search","description":"搜索工具"}]`,
		`[{"name":"query","description":"查询工具"},{"name":"update","description":"更新工具"}]`,
	}

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("mcp-server-%s", randomHex(4))
		port := 3000 + randomInt(6000)
		address := fmt.Sprintf("http://10.%d.%d.%d:%d", randomInt(256), randomInt(256), randomInt(256), port)
		department := randomChoice(depts)
		protocol := randomChoice(protocols)
		tools := randomChoice(toolsList)

		s := model.MCPServer{
			Name:       name,
			Address:    address,
			Department: department,
			Protocol:   protocol,
			Tools:      tools,
		}
		if err := database.DB.Create(&s).Error; err != nil {
			fmt.Printf("FAIL: create mcp_server %s: %v\n", name, err)
		} else {
			fmt.Printf("OK:   created mcp_server id=%d name=%s address=%s\n", s.ID, s.Name, s.Address)
		}
	}
}

func printHelp(table string) {
	switch table {
	case "user":
		fmt.Println(`=== datagen user ===
Insert randomly generated user records into the "users" table.

Fields generated:
  uid          random 16-char hex string
  username     "testuser_" + random 8-char hex
  password     bcrypt hash of "123456"
  email        {username}@example.com
  phone        138xxxxxxxx
  organization random: 计算机学院 / 信息中心 / 网络中心 / 图书馆 / 教务处 / 学生处 / 后勤管理处 / 理学院 / 文学院
  role_id      NULL (未分配角色)
  status       1 (enabled)

Usage:
  datagen -table user -count 5`)
	case "access_key":
		fmt.Println(`=== datagen access_key ===
Insert randomly generated access key records into the "access_keys" table.
Requires at least one user in the database (picks user_id randomly).

Fields generated:
  name       "ak_" + random 8-char hex
  key        "ak-" + random 48-char hex
  user_id    random existing user ID
  enabled    true
  servers    random: server-1 / server-2 / server-3
  expired_at randomly set or nil (50% chance)

Usage:
  datagen -table access_key -count 10`)
	case "usage_history":
		fmt.Println(`=== datagen usage_history ===
Insert randomly generated usage history records into the "usage_histories" table.
Requires at least one user and one access_key in the database.

Fields generated:
  user_id       random existing user ID
  access_key_id random existing access_key ID
  server        random: server-1 / server-2 / server-3
  endpoint      random: /api/tools / /api/resources / /api/prompts
  status        random: success / failed / timeout

Usage:
  datagen -table usage_history -count 20`)
	case "mcp_server":
		fmt.Println(`=== datagen mcp_server ===
Insert randomly generated MCP server records into the "mcp_servers" table.

Fields generated:
  name       "mcp-server-" + random 8-char hex
  address    http://10.x.x.x:random_port
  department random: 计算机学院 / 信息中心 / 网络中心 / 图书馆 / 教务处
  protocol   random: SSE / streamable-http
  tools      random JSON tool list (empty / single / multiple)

Usage:
  datagen -table mcp_server -count 3`)
	default:
		fmt.Println("datagen - test data generation tool")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  datagen -table <table> [-count <n>]")
		fmt.Println("  datagen -h -table <table>")
		fmt.Println()
		fmt.Println("Tables: user, access_key, usage_history, mcp_server")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func randomInt(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max))
	return n.Int64()
}

func randomChoice(options []string) string {
	return options[randomInt(int64(len(options)))]
}

func randomFromSlice[T any](slice []T) T {
	return slice[randomInt(int64(len(slice)))]
}
