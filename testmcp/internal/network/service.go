package network

import (
	"context"
	"fmt"
	"time"

	"testmcp/pkg/mcp"
)

var students = map[string]string{
	"2021010101": "张三",
	"2021010102": "李四",
	"2021010103": "王五",
	"2021010104": "赵六",
	"2021010105": "陈七",
}

type accountInfo struct {
	Account        string  `json:"account"`
	Status         string  `json:"status"`
	IP             string  `json:"ip"`
	MAC            string  `json:"mac"`
	Dorm           string  `json:"dorm"`
	MonthlyQuotaGB float64 `json:"monthly_quota_gb"`
	UsedGB         float64 `json:"used_gb"`
	RemainingGB    float64 `json:"remaining_gb"`
	ExpireDate     string  `json:"expire_date"`
}

var accounts = map[string]accountInfo{
	"2021010101": {
		Account:        "2021010101@stu.xaut.edu.cn",
		Status:         "active",
		IP:             "10.10.23.45",
		MAC:            "AA:BB:CC:DD:EE:01",
		Dorm:           "学生公寓3号楼512室",
		MonthlyQuotaGB: 100,
		UsedGB:         65.3,
		RemainingGB:    34.7,
		ExpireDate:     "2026-07-15",
	},
	"2021010102": {
		Account:        "2021010102@stu.xaut.edu.cn",
		Status:         "active",
		IP:             "10.10.23.46",
		MAC:            "AA:BB:CC:DD:EE:02",
		Dorm:           "学生公寓5号楼208室",
		MonthlyQuotaGB: 100,
		UsedGB:         88.1,
		RemainingGB:    11.9,
		ExpireDate:     "2026-07-15",
	},
	"2021010103": {
		Account:        "2021010103@stu.xaut.edu.cn",
		Status:         "active",
		IP:             "10.10.24.10",
		MAC:            "AA:BB:CC:DD:EE:03",
		Dorm:           "学生公寓2号楼315室",
		MonthlyQuotaGB: 100,
		UsedGB:         45.2,
		RemainingGB:    54.8,
		ExpireDate:     "2026-07-15",
	},
	"2021010104": {
		Account:        "2021010104@stu.xaut.edu.cn",
		Status:         "suspended",
		IP:             "10.10.23.47",
		MAC:            "AA:BB:CC:DD:EE:04",
		Dorm:           "学生公寓7号楼102室",
		MonthlyQuotaGB: 100,
		UsedGB:         99.9,
		RemainingGB:    0.1,
		ExpireDate:     "2026-07-15",
	},
	"2021010105": {
		Account:        "2021010105@stu.xaut.edu.cn",
		Status:         "active",
		IP:             "10.10.23.48",
		MAC:            "AA:BB:CC:DD:EE:05",
		Dorm:           "学生公寓3号楼604室",
		MonthlyQuotaGB: 100,
		UsedGB:         32.7,
		RemainingGB:    67.3,
		ExpireDate:     "2026-07-15",
	},
}

func RegisterTools(srv *mcp.Server) {
	srv.RegisterTool(mcp.Tool{
		Name:        "query_network_info",
		Description: "查询校园网络账号信息，包括IP、MAC、流量使用情况",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
			},
			Required: []string{"student_id"},
		},
	}, queryNetworkInfo)

	srv.RegisterTool(mcp.Tool{
		Name:        "query_network_usage",
		Description: "查询校园网流量使用详情",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
			},
			Required: []string{"student_id"},
		},
	}, queryNetworkUsage)

	srv.RegisterTool(mcp.Tool{
		Name:        "reset_network_password",
		Description: "重置校园网络账号密码",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
			},
			Required: []string{"student_id"},
		},
	}, resetPassword)

	srv.RegisterTool(mcp.Tool{
		Name:        "suspend_account",
		Description: "暂停/启用校园网络账号",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
				"action":     {Type: "string", Description: "操作类型: suspend(停用) 或 resume(恢复)"},
			},
			Required: []string{"student_id", "action"},
		},
	}, suspendAccount)
}

func getStr(args map[string]interface{}, key string) string {
	if v, ok := args[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func textF(format string, a ...interface{}) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf(format, a...)}},
	}
}

func errT(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{{Type: "text", Text: msg}},
		IsError: true,
	}
}

func checkStudent(args map[string]interface{}) (string, *mcp.CallToolResult) {
	sid := getStr(args, "student_id")
	if sid == "" {
		return "", errT("缺少参数: student_id")
	}
	if _, ok := students[sid]; !ok {
		return "", errT(fmt.Sprintf("学号 %s 不存在", sid))
	}
	return sid, nil
}

func queryNetworkInfo(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	info := accounts[sid]
	out := fmt.Sprintf("学生: %s (%s)\n\n", sid, students[sid])
	out += fmt.Sprintf("网络账号: %s\n", info.Account)
	out += fmt.Sprintf("账号状态: %s\n", info.Status)
	out += fmt.Sprintf("IP地址:   %s\n", info.IP)
	out += fmt.Sprintf("MAC地址:  %s\n", info.MAC)
	out += fmt.Sprintf("宿舍:     %s\n", info.Dorm)
	out += fmt.Sprintf("月流量配额: %.0f GB\n", info.MonthlyQuotaGB)
	out += fmt.Sprintf("已使用:     %.1f GB\n", info.UsedGB)
	out += fmt.Sprintf("剩余:       %.1f GB\n", info.RemainingGB)
	out += fmt.Sprintf("账号有效期: %s\n", info.ExpireDate)
	out += fmt.Sprintf("\n使用率: %.1f%% | %s",
		info.UsedGB/info.MonthlyQuotaGB*100,
		usageWarning(info.UsedGB, info.MonthlyQuotaGB))

	return textF(out), nil
}

func queryNetworkUsage(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	info := accounts[sid]
	usagePct := info.UsedGB / info.MonthlyQuotaGB * 100

	out := fmt.Sprintf("学生: %s (%s)  流量使用详情\n\n", sid, students[sid])
	out += fmt.Sprintf("账号: %s\n", info.Account)
	out += "----------------------------------------\n"
	out += fmt.Sprintf("月配额:     %8.0f GB\n", info.MonthlyQuotaGB)
	out += fmt.Sprintf("已使用:     %8.1f GB\n", info.UsedGB)
	out += fmt.Sprintf("剩余:       %8.1f GB\n", info.RemainingGB)
	out += fmt.Sprintf("使用率:     %7.1f%%\n", usagePct)
	out += "----------------------------------------\n"

	barLen := 30
	filled := int(usagePct / 100.0 * float64(barLen))
	if filled > barLen {
		filled = barLen
	}
	bar := ""
	for i := 0; i < barLen; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	out += fmt.Sprintf("[%s]\n\n", bar)
	out += usageWarning(info.UsedGB, info.MonthlyQuotaGB)

	return textF(out), nil
}

func usageWarning(used, quota float64) string {
	pct := used / quota * 100
	switch {
	case pct >= 90:
		return "⚠️  流量即将用尽，请注意节约使用或申请额外流量！"
	case pct >= 70:
		return "⚡ 流量使用已超过70%，建议合理安排使用。"
	default:
		return "✅ 流量充足，正常使用中。"
	}
}

func resetPassword(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	info := accounts[sid]
	newPass := fmt.Sprintf("Xaut@%s#%d", sid[4:8], time.Now().UnixNano()%9000+1000)

	out := fmt.Sprintf("网络账号密码重置成功!\n\n")
	out += fmt.Sprintf("学号:     %s\n", sid)
	out += fmt.Sprintf("姓名:     %s\n", students[sid])
	out += fmt.Sprintf("账号:     %s\n", info.Account)
	out += fmt.Sprintf("新密码:   %s\n", newPass)
	out += fmt.Sprintf("重置时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	out += "请妥善保管新密码，首次登录后建议修改。"

	return textF(out), nil
}

func suspendAccount(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	action := getStr(args, "action")
	info := accounts[sid]

	switch action {
	case "suspend":
		if info.Status == "suspended" {
			return textF("账号 %s 已处于停用状态，无需重复操作", info.Account), nil
		}
		info.Status = "suspended"
		accounts[sid] = info
		return textF("账号 %s (%s) 已成功停用。\n如需恢复，请使用 resume 操作。", info.Account, students[sid]), nil

	case "resume":
		if info.Status == "active" {
			return textF("账号 %s 已处于启用状态，无需重复操作", info.Account), nil
		}
		info.Status = "active"
		accounts[sid] = info
		return textF("账号 %s (%s) 已成功恢复使用。", info.Account, students[sid]), nil

	default:
		return errT("无效的操作类型，请使用 suspend 或 resume"), nil
	}
}
