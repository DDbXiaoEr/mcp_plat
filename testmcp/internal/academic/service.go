package academic

import (
	"context"
	"fmt"

	"testmcp/pkg/mcp"
)

var students = map[string]string{
	"2021010101": "张三",
	"2021010102": "李四",
	"2021010103": "王五",
	"2021010104": "赵六",
	"2021010105": "陈七",
}

type course struct {
	Name     string `json:"name"`
	Teacher  string `json:"teacher"`
	Schedule string `json:"schedule"`
	Room     string `json:"room"`
	Credits  int    `json:"credits"`
}

type grade struct {
	Course string  `json:"course"`
	Score  float64 `json:"score"`
	Credits int    `json:"credits"`
	Term   string  `json:"term"`
}

type exam struct {
	Course string `json:"course"`
	Time   string `json:"time"`
	Room   string `json:"room"`
	Seat   int    `json:"seat"`
}

var courses = map[string][]course{
	"2021010101": {
		{Name: "高等数学A(2)", Teacher: "张教授", Schedule: "周一 08:00-09:40", Room: "教一楼301", Credits: 5},
		{Name: "大学物理B(1)", Teacher: "李教授", Schedule: "周二 10:00-11:40", Room: "教二楼205", Credits: 4},
		{Name: "程序设计基础", Teacher: "王教授", Schedule: "周三 08:00-09:40", Room: "实验楼A301", Credits: 3},
		{Name: "数据结构与算法", Teacher: "赵教授", Schedule: "周四 14:00-15:40", Room: "教一楼402", Credits: 4},
		{Name: "大学英语(4)", Teacher: "陈教授", Schedule: "周五 08:00-09:40", Room: "教三楼201", Credits: 2},
	},
	"2021010102": {
		{Name: "高等数学A(2)", Teacher: "张教授", Schedule: "周一 10:00-11:40", Room: "教一楼302", Credits: 5},
		{Name: "大学物理B(1)", Teacher: "李教授", Schedule: "周二 08:00-09:40", Room: "教二楼206", Credits: 4},
		{Name: "电路分析基础", Teacher: "周教授", Schedule: "周三 14:00-15:40", Room: "教一楼501", Credits: 3},
		{Name: "模拟电子技术", Teacher: "吴教授", Schedule: "周四 10:00-11:40", Room: "教一楼503", Credits: 4},
		{Name: "大学英语(4)", Teacher: "陈教授", Schedule: "周五 10:00-11:40", Room: "教三楼202", Credits: 2},
	},
	"2021010103": {
		{Name: "计算机网络", Teacher: "刘教授", Schedule: "周一 14:00-15:40", Room: "教一楼601", Credits: 3},
		{Name: "操作系统", Teacher: "周教授", Schedule: "周二 14:00-15:40", Room: "教一楼602", Credits: 4},
		{Name: "数据库系统概论", Teacher: "王教授", Schedule: "周三 10:00-11:40", Room: "教一楼701", Credits: 3},
		{Name: "软件工程", Teacher: "赵教授", Schedule: "周四 08:00-09:40", Room: "教一楼702", Credits: 3},
		{Name: "编译原理", Teacher: "吴教授", Schedule: "周五 14:00-15:40", Room: "教一楼703", Credits: 3},
	},
	"2021010104": {
		{Name: "高等数学A(2)", Teacher: "张教授", Schedule: "周一 08:00-09:40", Room: "教一楼301", Credits: 5},
		{Name: "线性代数", Teacher: "郑教授", Schedule: "周二 14:00-15:40", Room: "教二楼301", Credits: 3},
		{Name: "概率论与数理统计", Teacher: "孙教授", Schedule: "周三 08:00-09:40", Room: "教二楼302", Credits: 3},
		{Name: "大学物理B(1)", Teacher: "李教授", Schedule: "周四 10:00-11:40", Room: "教二楼205", Credits: 4},
	},
	"2021010105": {
		{Name: "程序设计基础", Teacher: "王教授", Schedule: "周一 10:00-11:40", Room: "实验楼A301", Credits: 3},
		{Name: "离散数学", Teacher: "钱教授", Schedule: "周三 14:00-15:40", Room: "教一楼801", Credits: 4},
		{Name: "数字逻辑", Teacher: "孙教授", Schedule: "周四 14:00-15:40", Room: "教一楼802", Credits: 3},
		{Name: "数据结构与算法", Teacher: "赵教授", Schedule: "周五 08:00-09:40", Room: "教一楼402", Credits: 4},
	},
}

var grades = map[string][]grade{
	"2021010101": {
		{Course: "高等数学A(1)", Score: 85, Credits: 5, Term: "2023-2024-1"},
		{Course: "大学物理B(1)", Score: 78, Credits: 4, Term: "2023-2024-1"},
		{Course: "程序设计基础", Score: 92, Credits: 3, Term: "2023-2024-1"},
		{Course: "线性代数", Score: 88, Credits: 3, Term: "2023-2024-1"},
		{Course: "大学英语(3)", Score: 82, Credits: 2, Term: "2023-2024-1"},
		{Course: "高等数学A(2)", Score: 90, Credits: 5, Term: "2023-2024-2"},
		{Course: "数字逻辑", Score: 76, Credits: 3, Term: "2023-2024-2"},
	},
	"2021010102": {
		{Course: "高等数学A(1)", Score: 72, Credits: 5, Term: "2023-2024-1"},
		{Course: "大学物理B(1)", Score: 80, Credits: 4, Term: "2023-2024-1"},
		{Course: "程序设计基础", Score: 65, Credits: 3, Term: "2023-2024-1"},
		{Course: "大学英语(3)", Score: 88, Credits: 2, Term: "2023-2024-1"},
		{Course: "高等数学A(2)", Score: 76, Credits: 5, Term: "2023-2024-2"},
		{Course: "电路分析基础", Score: 82, Credits: 3, Term: "2023-2024-2"},
	},
	"2021010103": {
		{Course: "高等数学A(1)", Score: 95, Credits: 5, Term: "2023-2024-1"},
		{Course: "大学物理B(1)", Score: 88, Credits: 4, Term: "2023-2024-1"},
		{Course: "程序设计基础", Score: 91, Credits: 3, Term: "2023-2024-1"},
		{Course: "大学英语(3)", Score: 85, Credits: 2, Term: "2023-2024-1"},
		{Course: "计算机网络", Score: 90, Credits: 3, Term: "2023-2024-2"},
		{Course: "操作系统", Score: 87, Credits: 4, Term: "2023-2024-2"},
	},
	"2021010104": {
		{Course: "高等数学A(1)", Score: 60, Credits: 5, Term: "2023-2024-1"},
		{Course: "大学物理B(1)", Score: 55, Credits: 4, Term: "2023-2024-1"},
		{Course: "程序设计基础", Score: 68, Credits: 3, Term: "2023-2024-1"},
		{Course: "大学英语(3)", Score: 62, Credits: 2, Term: "2023-2024-1"},
		{Course: "高等数学A(2)", Score: 58, Credits: 5, Term: "2023-2024-2"},
	},
	"2021010105": {
		{Course: "高等数学A(1)", Score: 83, Credits: 5, Term: "2023-2024-1"},
		{Course: "大学物理B(1)", Score: 79, Credits: 4, Term: "2023-2024-1"},
		{Course: "程序设计基础", Score: 90, Credits: 3, Term: "2023-2024-1"},
		{Course: "大学英语(3)", Score: 76, Credits: 2, Term: "2023-2024-1"},
		{Course: "离散数学", Score: 85, Credits: 4, Term: "2023-2024-2"},
		{Course: "数字逻辑", Score: 72, Credits: 3, Term: "2023-2024-2"},
	},
}

var exams = map[string][]exam{
	"2021010101": {
		{Course: "高等数学A(2)", Time: "2024-06-20 08:00-10:00", Room: "教一楼301", Seat: 15},
		{Course: "大学物理B(1)", Time: "2024-06-22 10:00-12:00", Room: "教二楼205", Seat: 8},
		{Course: "程序设计基础", Time: "2024-06-24 08:00-10:00", Room: "实验楼A301", Seat: 22},
		{Course: "数据结构与算法", Time: "2024-06-26 14:00-16:00", Room: "教一楼402", Seat: 5},
		{Course: "大学英语(4)", Time: "2024-06-28 08:00-10:00", Room: "教三楼201", Seat: 30},
	},
	"2021010102": {
		{Course: "高等数学A(2)", Time: "2024-06-20 10:00-12:00", Room: "教一楼302", Seat: 12},
		{Course: "大学物理B(1)", Time: "2024-06-22 08:00-10:00", Room: "教二楼206", Seat: 20},
		{Course: "电路分析基础", Time: "2024-06-25 14:00-16:00", Room: "教一楼501", Seat: 3},
		{Course: "模拟电子技术", Time: "2024-06-27 10:00-12:00", Room: "教一楼503", Seat: 18},
		{Course: "大学英语(4)", Time: "2024-06-28 10:00-12:00", Room: "教三楼202", Seat: 25},
	},
	"2021010103": {
		{Course: "计算机网络", Time: "2024-06-21 14:00-16:00", Room: "教一楼601", Seat: 7},
		{Course: "操作系统", Time: "2024-06-23 14:00-16:00", Room: "教一楼602", Seat: 14},
		{Course: "数据库系统概论", Time: "2024-06-25 10:00-12:00", Room: "教一楼701", Seat: 9},
		{Course: "软件工程", Time: "2024-06-27 08:00-10:00", Room: "教一楼702", Seat: 21},
		{Course: "编译原理", Time: "2024-06-29 14:00-16:00", Room: "教一楼703", Seat: 33},
	},
	"2021010104": {
		{Course: "高等数学A(2)", Time: "2024-06-20 08:00-10:00", Room: "教一楼301", Seat: 28},
		{Course: "线性代数", Time: "2024-06-23 14:00-16:00", Room: "教二楼301", Seat: 16},
		{Course: "概率论与数理统计", Time: "2024-06-26 08:00-10:00", Room: "教二楼302", Seat: 11},
		{Course: "大学物理B(1)", Time: "2024-06-28 14:00-16:00", Room: "教二楼205", Seat: 2},
	},
	"2021010105": {
		{Course: "程序设计基础", Time: "2024-06-24 08:00-10:00", Room: "实验楼A301", Seat: 35},
		{Course: "离散数学", Time: "2024-06-26 14:00-16:00", Room: "教一楼801", Seat: 6},
		{Course: "数字逻辑", Time: "2024-06-28 14:00-16:00", Room: "教一楼802", Seat: 19},
		{Course: "数据结构与算法", Time: "2024-06-30 08:00-10:00", Room: "教一楼402", Seat: 27},
	},
}

func RegisterTools(srv *mcp.Server) {
	srv.RegisterTool(mcp.Tool{
		Name:        "query_courses",
		Description: "查询学生课程表",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
				"term":       {Type: "string", Description: "学期，如2023-2024-2，可选"},
			},
			Required: []string{"student_id"},
		},
	}, queryCourses)

	srv.RegisterTool(mcp.Tool{
		Name:        "query_grades",
		Description: "查询学生成绩",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
				"term":       {Type: "string", Description: "学期，如2023-2024-2，可选"},
			},
			Required: []string{"student_id"},
		},
	}, queryGrades)

	srv.RegisterTool(mcp.Tool{
		Name:        "query_gpa",
		Description: "查询学生平均绩点(GPA)",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
			},
			Required: []string{"student_id"},
		},
	}, queryGPA)

	srv.RegisterTool(mcp.Tool{
		Name:        "query_exam_schedule",
		Description: "查询考试安排",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"student_id": {Type: "string", Description: "学号"},
			},
			Required: []string{"student_id"},
		},
	}, queryExamSchedule)
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

func scoreToGradePoint(score float64) float64 {
	if score >= 90 {
		return 4.0
	} else if score >= 85 {
		return 3.7
	} else if score >= 82 {
		return 3.3
	} else if score >= 78 {
		return 3.0
	} else if score >= 75 {
		return 2.7
	} else if score >= 72 {
		return 2.3
	} else if score >= 68 {
		return 2.0
	} else if score >= 64 {
		return 1.5
	} else if score >= 60 {
		return 1.0
	}
	return 0.0
}

func queryCourses(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	clist, ok := courses[sid]
	if !ok || len(clist) == 0 {
		return textF("学号 %s (%s) 暂无课程信息", sid, students[sid]), nil
	}

	out := fmt.Sprintf("学生: %s (%s)  课程表\n\n", sid, students[sid])
	for _, c := range clist {
		out += fmt.Sprintf("%-20s %-8s %-22s %-12s %d学分\n", c.Name, c.Teacher, c.Schedule, c.Room, c.Credits)
	}
	return textF(out), nil
}

func queryGrades(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	glist, ok := grades[sid]
	if !ok || len(glist) == 0 {
		return textF("学号 %s (%s) 暂无成绩记录", sid, students[sid]), nil
	}

	termFilter := getStr(args, "term")

	out := fmt.Sprintf("学生: %s (%s)  成绩单\n\n", sid, students[sid])
	out += fmt.Sprintf("%-20s %8s %6s %s\n", "课程", "成绩", "学分", "学期")
	out += "----------------------------------------------\n"

	totalCredits := 0
	totalPoints := 0.0

	for _, g := range glist {
		if termFilter != "" && g.Term != termFilter {
			continue
		}
		out += fmt.Sprintf("%-20s %8.0f %6d %s\n", g.Course, g.Score, g.Credits, g.Term)
		totalCredits += g.Credits
		totalPoints += scoreToGradePoint(g.Score) * float64(g.Credits)
	}

	if totalCredits > 0 {
		gpa := totalPoints / float64(totalCredits)
		out += fmt.Sprintf("----------------------------------------------\n")
		out += fmt.Sprintf("加权GPA: %.2f", gpa)
	}

	return textF(out), nil
}

func queryGPA(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	glist, ok := grades[sid]
	if !ok || len(glist) == 0 {
		return textF("学号 %s (%s) 暂无成绩记录，无法计算GPA", sid, students[sid]), nil
	}

	out := fmt.Sprintf("学生: %s (%s)  GPA明细\n\n", sid, students[sid])
	out += fmt.Sprintf("%-20s %8s %6s %8s %6s\n", "课程", "成绩", "学分", "绩点", "学期")
	out += "--------------------------------------------------------\n"

	totalCredits := 0
	totalPoints := 0.0

	for _, g := range glist {
		gp := scoreToGradePoint(g.Score)
		out += fmt.Sprintf("%-20s %8.0f %6d %8.2f %6s\n", g.Course, g.Score, g.Credits, gp, g.Term)
		totalCredits += g.Credits
		totalPoints += gp * float64(g.Credits)
	}

	gpa := totalPoints / float64(totalCredits)
	out += "--------------------------------------------------------\n"
	out += fmt.Sprintf("总学分: %d  加权GPA: %.2f\n", totalCredits, gpa)

	var rank string
	switch {
	case gpa >= 3.7:
		rank = "优秀"
	case gpa >= 3.0:
		rank = "良好"
	case gpa >= 2.0:
		rank = "中等"
	default:
		rank = "需努力"
	}
	out += fmt.Sprintf("评级: %s", rank)

	return textF(out), nil
}

func queryExamSchedule(_ context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
	sid, errResp := checkStudent(args)
	if errResp != nil {
		return errResp, nil
	}

	elist, ok := exams[sid]
	if !ok || len(elist) == 0 {
		return textF("学号 %s (%s) 暂无考试安排", sid, students[sid]), nil
	}

	out := fmt.Sprintf("学生: %s (%s)  考试安排\n\n", sid, students[sid])
	out += fmt.Sprintf("%-20s %-22s %-14s %s\n", "课程", "时间", "教室", "座位")
	out += "--------------------------------------------------------------\n"
	for _, e := range elist {
		out += fmt.Sprintf("%-20s %-22s %-14s %d号\n", e.Course, e.Time, e.Room, e.Seat)
	}
	out += "--------------------------------------------------------------\n"
	out += fmt.Sprintf("共 %d 门考试，请提前15分钟到达考场，携带学生证。", len(elist))

	return textF(out), nil
}
