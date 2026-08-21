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
	"encoding/json"
	"io"
	"net/http"

	"mcp_plat-console/logging"
	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct{}

func NewSettingHandler() *SettingHandler {
	return &SettingHandler{}
}

func (h *SettingHandler) Get(c *gin.Context) {
	settings, err := service.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": settings})
}

func (h *SettingHandler) GatewayStatus(c *gin.Context) {
	status := service.CheckGatewayStatus()
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": status})
}

func (h *SettingHandler) GetByKey(c *gin.Context) {
	key := c.Param("key")
	val, err := service.GetSettingByKey(key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": val})
}

func (h *SettingHandler) Save(c *gin.Context) {
	key := c.Param("key")

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := service.SaveSetting(key, json.RawMessage(body)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if key == "user_ops" {
		service.ReloadAccessKeyCron()
	}

	if key == "log" {
		logging.Reconfigure()
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功"})
}

func (h *SettingHandler) TestLdapMapping(c *gin.Context) {
	var input service.TestLdapInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	output := service.TestLdapMapping(input)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": output})
}
