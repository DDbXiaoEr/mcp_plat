package library

import (
	"context"
	"fmt"
	"strings"
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

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Category string `json:"category"`
	ISBN     string `json:"isbn"`
	Location string `json:"location"`
	Status   string `json:"status"`
}

type loan struct {
	BookID     string `json:"book_id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	BorrowDate string `json:"borrow_date"`
	DueDate    string `json:"due_date"`
	Renewed    int    `json:"renewed"`
}

type history struct {
	BookID     string `json:"book_id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	BorrowDate string `json:"borrow_date"`
	ReturnDate string `json:"return_date"`
}

var libraryBooks = map[string]book{
	"B001": {ID: "B001", Title: "数据结构与算法分析", Author: "Mark Allen Weiss", Category: "计算机科学", ISBN: "978-7-111-35802-1", Location: "科技图书区A-3-12", Status: "在馆"},
	"B002": {ID: "B002", Title: "计算机网络：自顶向下方法", Author: "James F. Kurose", Category: "计算机科学", ISBN: "978-7-111-45362-2", Location: "科技图书区A-2-8", Status: "在馆"},
	"B003": {ID: "B003", Title: "深入理解计算机系统", Author: "Randal E. Bryant", Category: "计算机科学", ISBN: "978-7-111-54493-7", Location: "科技图书区A-1-15", Status: "在馆"},
	"B004": {ID: "B004", Title: "C程序设计语言", Author: "Brian W. Kernighan", Category: "计算机科学", ISBN: "978-7-111-19626-5", Location: "科技图书区A-3-5", Status: "已借出"},
	"B005": {ID: "B005", Title: "操作系统概念", Author: "Abraham Silberschatz", Category: "计算机科学", ISBN: "978-7-111-50473-5", Location: "科技图书区A-2-20", Status: "在馆"},
	"B006": {ID: "B006", Title: "数据库系统概念", Author: "Abraham Silberschatz", Category: "计算机科学", ISBN: "978-7-111-52173-8", Location: "科技图书区A-4-3", Status: "在馆"},
	"B007": {ID: "B007", Title: "算法导论", Author: "Thomas H. Cormen", Category: "计算机科学", ISBN: "978-7-111-40701-6", Location: "科技图书区A-3-18", Status: "在馆"},
	"B008": {ID: "B008", Title: "高等数学（第七版）", Author: "同济大学数学系", Category: "数学", ISBN: "978-7-04-039663-8", Location: "自然科学区B-1-10", Status: "在馆"},
	"B009": {ID: "B009", Title: "线性代数及其应用", Author: "David C. Lay", Category: "数学", ISBN: "978-7-111-48582-5", Location: "自然科学区B-1-12", Status: "在馆"},
	"B010": {ID: "B010", Title: "概率论与数理统计", Author: "盛骤", Category: "数学", ISBN: "978-7-04-023896-9", Location: "自然科学区B-2-5", Status: "在馆"},
	"B011": {ID: "B011", Title: "人工智能：一种现代方法", Author: "Stuart Russell", Category: "计算机科学", ISBN: "978-7-302-47236-1", Location: "科技图书区A-5-8", Status: "在馆"},
	"B012": {ID: "B012", Title: "编译原理", Author: "Alfred V. Aho", Category: "计算机科学", ISBN: "978-7-111-25121-7", Location: "科技图书区A-4-15", Status: "在馆"},
	"B013": {ID: "B013", Title: "软件工程", Author: "Roger S. Pressman", Category: "计算机科学", ISBN: "978-7-111-44776-5", Location: "科技图书区A-5-3", Status: "在馆"},
	"B014": {ID: "B014", Title: "大学物理", Author: "马文蔚", Category: "物理学", ISBN: "978-7-04-039012-4", Location: "自然科学区C-3-8", Status: "在馆"},
	"B015": {ID: "B015", Title: "离散数学及其应用", Author: "Kenneth H. Rosen", Category: "数学", ISBN: "978-7-111-48846-8", Location: "自然科学区B-3-7", Status: "在馆"},
}

var currentLoans = map[string][]loan{
	"2021010101": {
		{BookID: "B001", Title: "数据结构与算法分析", Author: "Mark Allen Weiss", BorrowDate: "2024-03-01", DueDate: "2024-04-15", Renewed: 1},
		{BookID: "B002", Title: "计算机网络：自顶向下方法", Author: "James F. Kurose", BorrowDate: "2024-03-15", DueDate: "2024-04-15", Renewed: 0},
		{BookID: "B003", Title: "深入理解计算机系统", Author: "Randal E. Bryant", BorrowDate: "2024-04-01", DueDate: "2024-05-01", Renewed: 0},
	},
	"2021010102": {
		{BookID: "B004", Title: "C程序设计语言", Author: "Brian W. Kernighan", BorrowDate: "2024-03-20", DueDate: "2024-04-20", Renewed: 0},
		{BookID: "B012", Title: "编译原理", Author: "Alfred V. Aho", BorrowDate: "2024-04-05", DueDate: "2024-05-05", Renewed: 0},
	},
	"2021010103": {
		{BookID: "B007", Title: "算法导论", Author: "Thomas H. Cormen", BorrowDate: "2024-02-15", DueDate: "2024-03-15", Renewed: 2},
		{BookID: "B011", Title: "人工智能：一种现代方法", Author: "Stuart Russell", BorrowDate: "2024-04-01", DueDate: "2024-05-01", Renewed: 0},
		{BookID: "B013", Title: "软件工程", Author: "Roger S. Pressman", BorrowDate: "2024-04-10", DueDate: "2024-05-10", Renewed: 0},
	},
	"2021010105": {
		{BookID: "B005", Title: "操作系统概念", Author: "Abraham Silberschatz", BorrowDate: "2024-04-08", DueDate: "2024-05-08", Renewed: 0},
	},
}

var borrowHistory = map[string][]history{
	"2021010101": {
		{BookID: "B004", Title: "C程序设计语言", Author: "Brian W. Kernighan", BorrowDate: "2023-09-15", ReturnDate: "2023-10-15"},
		{BookID: "B005", Title: "操作系统概念", Author: "Abraham Silberschatz", BorrowDate: "2023-10-01", ReturnDate: "2023-11-01"},
		{BookID: "B006", Title: "数据库系统概念", Author: "Abraham Silberschatz", BorrowDate: "2023-11-15", ReturnDate: "2023-12-15"},
		{BookID: "B008", Title: "高等数学（第七版）", Author: "同济大学数学系", BorrowDate: "2023-09-01", ReturnDate: "2024-01-15"},
		{BookID: "B009", Title: "线性代数及其应用", Author: "David C. Lay", BorrowDate: "2023-09-15", ReturnDate: "2024-01-15"},
	},
	"2021010102": {
		{BookID: "B008", Title: "高等数学（第七版）", Author: "同济大学数学系", BorrowDate: "2023-09-05", ReturnDate: "2024-01-10"},
		{BookID: "B014", Title: "大学物理", Author: "马文蔚", BorrowDate: "2023-10-10", ReturnDate: "2023-12-20"},
	},
	"2021010103": {
		{BookID: "B001", Title: "数据结构与算法分析", Author: "Mark Allen Weiss", BorrowDate: "2023-09-20", ReturnDate: "2023-10-20"},
		{BookID: "B006", Title: "数据库系统概念", Author: "Abraham Silberschatz", BorrowDate: "2023-10-25", ReturnDate: "2023-11-25"},
		{BookID: "B012", Title: "编译原理", Author: "Alfred V. Aho", BorrowDate: "2023-11-01", ReturnDate: "2023-12-01"},
		{BookID: "B015", Title: "离散数学及其应用", Author: "Kenneth H. Rosen", BorrowDate: "2023-09-01", ReturnDate: "2024-01-15"},
		{BookID: "B010", Title: "概率论与数理统计", Author: "盛骤", BorrowDate: "2023-12-01", ReturnDate: "2024-01-20"},
	},
	"2021010104": {
		{BookID: "B008", Title: "高等数学（第七版）", Author: "同济大学数学系", BorrowDate: "2023-09-01", ReturnDate: "2024-01-15"},
	},
}

func RegisterTools(srv *mcp.Server) {
	srv.RegisterTool(mcp.Tool{
		Name:        "query_current_loans",
		Description: "查询当前借阅的图书",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
			},
			Required: []string{"student_id"},
		},
	}, queryCurrentLoans)

	srv.RegisterTool(mcp.Tool{
		Name:        "query_borrow_history",
		Description: "查询历史借阅记录",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
				"limit":      {Type: "number", Description: "返回记录数，默认10"},
			},
			Required: []string{"student_id"},
		},
	}, queryBorrowHistory)

	srv.RegisterTool(mcp.Tool{
		Name:        "renew_book",
		Description: "续借图书",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
				"book_id":    {Type: "string", Description: "图书ID"},
			},
			Required: []string{"student_id", "book_id"},
		},
	}, renewBook)

	srv.RegisterTool(mcp.Tool{
		Name:        "search_books",
		Description: "搜索图书馆藏书",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"keyword": {Type: "string", Description: "搜索关键词（书名或作者）"},
			},
			Required: []string{"keyword"},
		},
	}, searchBooks)
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

func queryCurrentLoans(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	loans, ok := currentLoans[sid]
	if !ok || len(loans) == 0 {
		return textF("学生: %s (%s) 当前无在借图书", sid, students[sid]), nil
	}

	now := time.Now()
	out := fmt.Sprintf("学生: %s (%s)  当前借阅\n\n", sid, students[sid])
	out += fmt.Sprintf("%-6s %-28s %-22s %-12s %s\n", "书号", "书名", "作者", "应还日期", "状态")
	out += strings.Repeat("-", 95) + "\n"

	for _, l := range loans {
		dueDate, _ := time.Parse("2006-01-02", l.DueDate)
		status := "正常"
		if now.After(dueDate) {
			status = fmt.Sprintf("超期%d天", int(now.Sub(dueDate).Hours()/24))
		} else if now.Add(7*24*time.Hour).After(dueDate) {
			status = "即将到期"
		}

		renewed := ""
		if l.Renewed > 0 {
			renewed = fmt.Sprintf(" (已续%d次)", l.Renewed)
		}
		out += fmt.Sprintf("%-6s %-28s %-22s %-12s %s%s\n",
			l.BookID, l.Title, l.Author, l.DueDate, status, renewed)
	}
	out += strings.Repeat("-", 95) + "\n"
	out += fmt.Sprintf("共 %d 本书在借\n", len(loans))
	out += "续借规则: 每本书最多续借2次，每次续借延长15天。请在到期前操作。"
	return textF(out), nil
}

func queryBorrowHistory(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	hist, ok := borrowHistory[sid]
	if !ok || len(hist) == 0 {
		return textF("学生: %s (%s) 暂无借阅历史", sid, students[sid]), nil
	}

	limit := 10
	if v, ok := args["limit"]; ok {
		switch val := v.(type) {
		case float64:
			limit = int(val)
		case int:
			limit = val
		}
	}

	out := fmt.Sprintf("学生: %s (%s)  借阅历史\n\n", sid, students[sid])
	out += fmt.Sprintf("%-6s %-28s %-22s %-12s %-12s\n", "书号", "书名", "作者", "借阅日期", "归还日期")
	out += strings.Repeat("-", 85) + "\n"

	count := 0
	for _, h := range hist {
		if count >= limit {
			break
		}
		out += fmt.Sprintf("%-6s %-28s %-22s %-12s %-12s\n",
			h.BookID, h.Title, h.Author, h.BorrowDate, h.ReturnDate)
		count++
	}
	out += strings.Repeat("-", 85) + "\n"
	out += fmt.Sprintf("共 %d 条记录", len(hist))
	return textF(out), nil
}

func renewBook(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	bookID := getStr(args, "book_id")
	if bookID == "" {
		return errT("缺少参数: book_id"), nil
	}

	loans, ok := currentLoans[sid]
	if !ok || len(loans) == 0 {
		return errT(fmt.Sprintf("学生 %s 当前无在借图书", sid)), nil
	}

	for i, l := range loans {
		if l.BookID == bookID {
			if l.Renewed >= 2 {
				return errT(fmt.Sprintf("图书 [%s]《%s》已续借 %d 次，达到最大续借次数，无法再次续借。", l.BookID, l.Title, l.Renewed)), nil
			}

			oldDue := l.DueDate
			dueDate, _ := time.Parse("2006-01-02", l.DueDate)
			newDue := dueDate.Add(15 * 24 * time.Hour)
			loans[i].Renewed++
			loans[i].DueDate = newDue.Format("2006-01-02")
			currentLoans[sid] = loans

			out := fmt.Sprintf("续借成功!\n\n")
			out += fmt.Sprintf("学号: %s (%s)\n", sid, students[sid])
			out += fmt.Sprintf("书号: %s\n", l.BookID)
			out += fmt.Sprintf("书名: 《%s》\n", l.Title)
			out += fmt.Sprintf("作者: %s\n", l.Author)
			out += fmt.Sprintf("原到期日: %s\n", oldDue)
			out += fmt.Sprintf("新到期日: %s\n", newDue.Format("2006-01-02"))
			out += fmt.Sprintf("续借次数: %d/2\n\n", loans[i].Renewed)
			out += "请在新到期日前归还或再次续借。"
			return textF(out), nil
		}
	}

	return errT(fmt.Sprintf("未找到图书 %s 的借阅记录，请确认书号是否正确", bookID)), nil
}

func searchBooks(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	keyword := strings.ToLower(getStr(args, "keyword"))
	if keyword == "" {
		return errT("缺少参数: keyword"), nil
	}

	var results []book
	for _, b := range libraryBooks {
		if strings.Contains(strings.ToLower(b.Title), keyword) ||
			strings.Contains(strings.ToLower(b.Author), keyword) ||
			strings.Contains(strings.ToLower(b.Category), keyword) {
			results = append(results, b)
		}
	}

	if len(results) == 0 {
		return textF("未找到与 \"%s\" 相关的图书", getStr(args, "keyword")), nil
	}

	out := fmt.Sprintf("搜索 \"%s\" 的结果，共 %d 条:\n\n", getStr(args, "keyword"), len(results))
	out += fmt.Sprintf("%-6s %-32s %-24s %-10s %s\n", "书号", "书名", "作者", "状态", "馆藏位置")
	out += strings.Repeat("-", 100) + "\n"
	for _, b := range results {
		out += fmt.Sprintf("%-6s %-32s %-24s %-10s %s\n",
			b.ID, b.Title, b.Author, b.Status, b.Location)
	}
	out += strings.Repeat("-", 100) + "\n"
	out += "如需借阅，请持学生证到图书馆办理。"

	return textF(out), nil
}
