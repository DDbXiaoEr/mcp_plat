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

// Package resp centralizes HTTP JSON responses so every handler shares one
// output path. The unified envelope is {"code":<int>,"message":"...","data":...}.
// The message is localized according to the request locale (set by the Locale
// middleware, falling back to the Accept-Language header).
package resp

// Author: deepseek-v4-pro / opencode

import (
	"net/http"

	"mcp_plat-console/i18n"

	"github.com/gin-gonic/gin"
)

// Locale returns the resolved language tag ("zh" or "en") for the request.
func Locale(c *gin.Context) string {
	if v := c.GetString("locale"); v != "" {
		return v
	}
	return i18n.Lang(c.GetHeader("Accept-Language"))
}

// JSON writes a localized response. status doubles as the business "code" field.
// When data is non-empty it is embedded under the "data" key.
func JSON(c *gin.Context, status int, message string, data ...any) {
	body := gin.H{"code": status, "message": i18n.Translate(Locale(c), message)}
	if len(data) > 0 {
		body["data"] = data[0]
	}
	c.JSON(status, body)
}

// OK writes a 200 response with the given message (and optional data).
func OK(c *gin.Context, message string, data ...any) {
	JSON(c, http.StatusOK, message, data...)
}

// Fail writes an error response whose HTTP status equals its business code.
func Fail(c *gin.Context, status int, message string) {
	JSON(c, status, message)
}
