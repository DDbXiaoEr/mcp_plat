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

// Author: deepseek-v4-pro / opencode

import (
	"net/http"

	"mcp_plat-console/resp"
	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type MCPServerHandler struct{}

func NewMCPServerHandler() *MCPServerHandler {
	return &MCPServerHandler{}
}

func (h *MCPServerHandler) List(c *gin.Context) {
	servers, err := service.ListServers()
	if err != nil {
		resp.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}

	resp.OK(c, "success", servers)
}

func (h *MCPServerHandler) Create(c *gin.Context) {
	var input service.CreateServerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	server, err := service.CreateServer(input)
	if err != nil {
		resp.Fail(c, http.StatusInternalServerError, "创建失败")
		return
	}

	resp.OK(c, "创建成功", server)
}

func (h *MCPServerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var input service.UpdateServerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.UpdateServer(id, input); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "更新成功")
}

func (h *MCPServerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := service.DeleteServer(id); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "删除成功")
}

func (h *MCPServerHandler) FetchTools(c *gin.Context) {
	var input service.FetchToolsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	tools, err := service.FetchTools(input)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "获取成功", gin.H{"tools": tools})
}

func (h *MCPServerHandler) Publish(c *gin.Context) {
	var input service.PublishInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.PublishServers(input); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "发布成功")
}

func (h *MCPServerHandler) Maintenance(c *gin.Context) {
	var input service.MaintenanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.SetMaintenance(input); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "维护设置成功")
}
