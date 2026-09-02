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

package resp_test

// Author: deepseek-v4-pro / opencode

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mcp_plat-console/middleware"
	"mcp_plat-console/resp"

	"github.com/gin-gonic/gin"
)

func perform(header string) (int, map[string]any) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.Locale())
	r.GET("/test", func(c *gin.Context) {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
	})
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	if header != "" {
		req.Header.Set("Accept-Language", header)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	return w.Code, body
}

func TestFailLocalized(t *testing.T) {
	code, body := perform("en-US,en;q=0.9")
	if code != http.StatusBadRequest {
		t.Fatalf("status = %d", code)
	}
	if body["code"] != float64(400) {
		t.Errorf("code = %v", body["code"])
	}
	if body["message"] != "Invalid parameters" {
		t.Errorf("en message = %v", body["message"])
	}
	if _, has := body["data"]; has {
		t.Errorf("data should be omitted on failure")
	}
}

func TestFailKeepsChinese(t *testing.T) {
	_, body := perform("zh-CN,zh;q=0.9")
	if body["message"] != "参数错误" {
		t.Errorf("zh message = %v", body["message"])
	}
	_, body = perform("")
	if body["message"] != "参数错误" {
		t.Errorf("default message = %v", body["message"])
	}
}

func TestOKWithData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.Locale())
	r.GET("/t", func(c *gin.Context) {
		resp.OK(c, "创建成功", gin.H{"id": 1})
	})
	req := httptest.NewRequest(http.MethodGet, "/t", nil)
	req.Header.Set("Accept-Language", "en")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "Created successfully" {
		t.Errorf("message = %v", body["message"])
	}
	data, ok := body["data"].(map[string]any)
	if !ok || data["id"] != float64(1) {
		t.Errorf("data = %v", body["data"])
	}
}
