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

	"mcp_plat-console/resp"
	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct{}

func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{}
}

func (h *HistoryHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	var query service.AuditLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	output, err := service.ListAuditLogs(userID, query)
	if err != nil {
		resp.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}

	resp.OK(c, "success", output)
}
