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

	"mcp_plat-console/i18n"
	"mcp_plat-console/resp"
	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type GatewayServerHandler struct{}

func NewGatewayServerHandler() *GatewayServerHandler {
	return &GatewayServerHandler{}
}

func (h *GatewayServerHandler) List(c *gin.Context) {
	gateways, err := service.ListGatewayServers()
	if err != nil {
		resp.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}

	locale := resp.Locale(c)
	for i := range gateways {
		gateways[i].StatusMessage = i18n.Translate(locale, gateways[i].StatusMessage)
	}

	resp.OK(c, "success", gateways)
}

func (h *GatewayServerHandler) Create(c *gin.Context) {
	var input service.GatewayServerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	gateway, err := service.CreateGatewayServer(input)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "创建成功", gateway)
}

func (h *GatewayServerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var input service.GatewayServerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.UpdateGatewayServer(id, input); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "更新成功")
}

func (h *GatewayServerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := service.DeleteGatewayServer(id); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "删除成功")
}

func (h *GatewayServerHandler) SetDefault(c *gin.Context) {
	id := c.Param("id")

	if err := service.SetDefaultGatewayServer(id); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "设置成功")
}

func (h *GatewayServerHandler) Test(c *gin.Context) {
	id := c.Param("id")

	result, err := service.TestGatewayServer(id)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	result.Message = i18n.Translate(resp.Locale(c), result.Message)
	resp.OK(c, "success", result)
}

func (h *GatewayServerHandler) Check(c *gin.Context) {
	var input service.GatewayCheckInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	result, err := service.TestGatewayServerConfig(input)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	result.Message = i18n.Translate(resp.Locale(c), result.Message)
	resp.OK(c, "success", result)
}
