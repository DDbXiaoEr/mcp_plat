package card

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"testmcp/pkg/mcp"
)

var (
	balances = map[string]float64{
		"2021010101": 250.80,
		"2021010102": 180.50,
		"2021010103": 520.00,
		"2021010104": 85.30,
		"2021010105": 320.00,
	}

	cardStatus = map[string]string{
		"2021010101": "normal",
		"2021010102": "normal",
		"2021010103": "normal",
		"2021010104": "lost",
		"2021010105": "normal",
	}

	transactions = map[string][]map[string]interface{}{
		"2021010101": {
			{"time": "2024-03-15 12:30", "location": "食堂一楼", "type": "消费", "amount": 15.00, "balance": 235.80},
			{"time": "2024-03-15 18:00", "location": "食堂二楼", "type": "消费", "amount": 18.50, "balance": 217.30},
			{"time": "2024-03-16 08:00", "location": "校园超市", "type": "消费", "amount": 12.00, "balance": 205.30},
			{"time": "2024-03-16 12:30", "location": "食堂一楼", "type": "消费", "amount": 16.00, "balance": 189.30},
			{"time": "2024-03-17 10:00", "location": "图书馆打印室", "type": "消费", "amount": 3.00, "balance": 186.30},
			{"time": "2024-03-18 12:30", "location": "食堂三楼", "type": "消费", "amount": 22.00, "balance": 164.30},
		},
		"2021010102": {
			{"time": "2024-03-15 08:00", "location": "食堂一楼", "type": "消费", "amount": 8.00, "balance": 172.50},
			{"time": "2024-03-15 13:00", "location": "食堂二楼", "type": "消费", "amount": 15.00, "balance": 157.50},
			{"time": "2024-03-16 09:00", "location": "校园超市", "type": "消费", "amount": 25.00, "balance": 132.50},
			{"time": "2024-03-20 14:00", "location": "在线充值", "type": "充值", "amount": 100.00, "balance": 232.50},
			{"time": "2024-03-21 12:00", "location": "食堂三楼", "type": "消费", "amount": 20.00, "balance": 212.50},
			{"time": "2024-03-22 18:00", "location": "浴室", "type": "消费", "amount": 3.00, "balance": 209.50},
		},
		"2021010103": {
			{"time": "2024-03-10 07:30", "location": "食堂一楼", "type": "消费", "amount": 10.00, "balance": 510.00},
			{"time": "2024-03-10 11:30", "location": "食堂二楼", "type": "消费", "amount": 18.00, "balance": 492.00},
			{"time": "2024-03-11 07:30", "location": "食堂一楼", "type": "消费", "amount": 8.00, "balance": 484.00},
			{"time": "2024-03-11 12:00", "location": "食堂三楼", "type": "消费", "amount": 25.00, "balance": 459.00},
			{"time": "2024-03-15 10:00", "location": "校园超市", "type": "消费", "amount": 35.00, "balance": 424.00},
			{"time": "2024-03-18 16:00", "location": "在线充值", "type": "充值", "amount": 200.00, "balance": 624.00},
		},
		"2021010105": {
			{"time": "2024-03-12 08:00", "location": "食堂二楼", "type": "消费", "amount": 12.00, "balance": 308.00},
			{"time": "2024-03-13 12:00", "location": "食堂一楼", "type": "消费", "amount": 16.00, "balance": 292.00},
			{"time": "2024-03-14 18:00", "location": "食堂三楼", "type": "消费", "amount": 22.00, "balance": 270.00},
			{"time": "2024-03-16 14:00", "location": "校园超市", "type": "消费", "amount": 18.00, "balance": 252.00},
			{"time": "2024-03-20 09:00", "location": "图书馆打印室", "type": "消费", "amount": 5.00, "balance": 247.00},
		},
	}
)

func RegisterTools(srv *mcp.Server) {
	srv.RegisterTool(mcp.Tool{
		Name:        "query_balance",
		Description: "查询一卡通账户余额",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
			},
			Required: []string{"student_id"},
		},
	}, queryBalance)

	srv.RegisterTool(mcp.Tool{
		Name:        "query_transactions",
		Description: "查询一卡通消费流水记录",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
				"date":       {Type: "string", Description: "查询日期，格式YYYY-MM-DD，可选"},
				"limit":      {Type: "number", Description: "返回记录条数，默认20"},
			},
			Required: []string{"student_id"},
		},
	}, queryTransactions)

	srv.RegisterTool(mcp.Tool{
		Name:        "report_loss",
		Description: "一卡通挂失，挂失后卡片无法使用",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
			},
			Required: []string{"student_id"},
		},
	}, reportLoss)

	srv.RegisterTool(mcp.Tool{
		Name:        "cancel_loss",
		Description: "一卡通解挂，恢复正常使用",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
			},
			Required: []string{"student_id"},
		},
	}, cancelLoss)

	srv.RegisterTool(mcp.Tool{
		Name:        "recharge",
		Description: "一卡通充值",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
				"amount":     {Type: "number", Description: "充值金额(元)"},
			},
			Required: []string{"student_id", "amount"},
		},
	}, recharge)
}

func getStringArg(args map[string]interface{}, key string) string {
	if v, ok := args[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func textResult(format string, a ...interface{}) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf(format, a...)}},
	}
}

func mustStudent(args map[string]interface{}) (string, *mcp.CallToolResult) {
	sid := getStringArg(args, "student_id")
	if sid == "" {
		return "", &mcp.CallToolResult{
			Content: []mcp.Content{{Type: "text", Text: "缺少参数: student_id"}},
			IsError: true,
		}
	}
	if _, ok := balances[sid]; !ok {
		return "", &mcp.CallToolResult{
			Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf("学号 %s 不存在", sid)}},
			IsError: true,
		}
	}
	return sid, nil
}

func queryBalance(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := mustStudent(args)
	if errResp != nil {
		return errResp, nil
	}
	bal := balances[sid]
	status := cardStatus[sid]
	return textResult("学号: %s\n余额: ¥%.2f\n状态: %s", sid, bal, status), nil
}

func queryTransactions(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := mustStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	txns, ok := transactions[sid]
	if !ok || len(txns) == 0 {
		return textResult("学号 %s 暂无消费记录", sid), nil
	}

	dateFilter := getStringArg(args, "date")

	out := fmt.Sprintf("学号: %s 一卡通消费流水\n", sid)
	out += "----------------------------------------\n"
	count := 0

	for _, txn := range txns {
		t := fmt.Sprintf("%v", txn["time"])
		if dateFilter != "" && len(t) >= 10 && t[:10] != dateFilter {
			continue
		}
		out += fmt.Sprintf("%s | %s | %s | ¥%.2f | 余额¥%.2f\n",
			txn["time"], txn["location"], txn["type"], txn["amount"], txn["balance"])
		count++
	}
	out += fmt.Sprintf("----------------------------------------\n共 %d 条记录", count)
	return textResult(out), nil
}

func reportLoss(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := mustStudent(args)
	if errResp != nil {
		return errResp, nil
	}
	if cardStatus[sid] == "lost" {
		return textResult("学号 %s 的一卡通已处于挂失状态，无需重复挂失", sid), nil
	}
	cardStatus[sid] = "lost"
	return textResult("学号 %s 的一卡通已成功挂失，挂失时间: %s\n请及时到一卡通中心补办新卡。", sid, time.Now().Format("2006-01-02 15:04:05")), nil
}

func cancelLoss(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := mustStudent(args)
	if errResp != nil {
		return errResp, nil
	}
	if cardStatus[sid] == "normal" {
		return textResult("学号 %s 的一卡通状态正常，无需解挂", sid), nil
	}
	cardStatus[sid] = "normal"
	return textResult("学号 %s 的一卡通已成功解挂，现在可以正常使用。", sid), nil
}

func recharge(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := mustStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	amount := 0.0
	if v, ok := args["amount"]; ok {
		switch val := v.(type) {
		case float64:
			amount = val
		case int:
			amount = float64(val)
		default:
			return &mcp.CallToolResult{
				Content: []mcp.Content{{Type: "text", Text: "参数amount必须为数字"}},
				IsError: true,
			}, nil
		}
	}
	if amount <= 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{{Type: "text", Text: "充值金额必须大于0"}},
			IsError: true,
		}, nil
	}

	oldBal := balances[sid]
	balances[sid] = oldBal + amount

	rec := map[string]interface{}{
		"time":     time.Now().Format("2006-01-02 15:04:05"),
		"location": "在线充值",
		"type":     "充值",
		"amount":   amount,
		"balance":  balances[sid],
	}
	transactions[sid] = append(transactions[sid], rec)

	return textResult("充值成功!\n学号: %s\n充值金额: ¥%.2f\n充值前余额: ¥%.2f\n当前余额: ¥%.2f\n交易流水号: %d",
		sid, amount, oldBal, balances[sid], rand.Intn(900000)+100000), nil
}
