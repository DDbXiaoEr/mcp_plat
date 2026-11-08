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

package middleware

// Author: deepseek-v4-pro / opencode

import (
	"errors"
	"net/http"

	"mcp_plat-console/config"
	"mcp_plat-console/resp"
	"mcp_plat-console/session"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := session.IDFromRequest(c)
		if sid == "" {
			resp.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}

		data, err := session.Get(sid)
		if err != nil {
			if errors.Is(err, session.ErrNotFound) {
				session.ClearCookie(c)
				resp.Fail(c, http.StatusUnauthorized, "会话无效或已过期")
			} else {
				resp.Fail(c, http.StatusInternalServerError, "会话存储不可用")
			}
			c.Abort()
			return
		}

		_ = session.Touch(sid)
		session.SetCookie(c, sid)

		c.Set("user_id", data.UserID)
		c.Set("username", data.Username)
		c.Set("role", data.Role)
		c.Set("session_id", sid)
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.AppConfig.AdminEnabled() {
			resp.Fail(c, http.StatusForbidden, "当前节点未启用管理员功能")
			c.Abort()
			return
		}

		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			resp.Fail(c, http.StatusForbidden, "仅管理员可操作")
			c.Abort()
			return
		}
		c.Next()
	}
}
