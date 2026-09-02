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

package handler

import (
	"net/http"

	"mcp_plat-console/database"
	"mcp_plat-console/resp"

	"github.com/gin-gonic/gin"
)

// Healthz 存活探针：进程存活即返回 200
func Healthz(c *gin.Context) {
	resp.OK(c, "ok")
}

// Readyz 就绪探针：进程存活且数据库可达才返回 200，否则 503
func Readyz(c *gin.Context) {
	sqlDB, err := database.DB.DB()
	if err != nil {
		resp.Fail(c, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	if err := sqlDB.Ping(); err != nil {
		resp.Fail(c, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	resp.OK(c, "ok")
}
