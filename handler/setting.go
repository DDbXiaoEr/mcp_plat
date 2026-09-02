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

	"mcp_plat-console/i18n"
	"mcp_plat-console/logging"
	"mcp_plat-console/resp"
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
		resp.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}

	resp.OK(c, "success", settings)
}

func (h *SettingHandler) GatewayStatus(c *gin.Context) {
	status := service.CheckGatewayStatus()
	resp.OK(c, "success", status)
}

func (h *SettingHandler) GetByKey(c *gin.Context) {
	key := c.Param("key")
	val, err := service.GetSettingByKey(key)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	resp.OK(c, "success", val)
}

func (h *SettingHandler) Save(c *gin.Context) {
	key := c.Param("key")

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.SaveSetting(key, json.RawMessage(body)); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if key == "user_ops" {
		service.ReloadAccessKeyCron()
	}

	if key == "log" {
		logging.Reconfigure()
	}

	resp.OK(c, "保存成功")
}

func (h *SettingHandler) TestLdapMapping(c *gin.Context) {
	var input service.TestLdapInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	output := service.TestLdapMapping(input)
	output.Message = i18n.Translate(resp.Locale(c), output.Message)
	resp.OK(c, "success", output)
}

func (h *SettingHandler) TestSmtp(c *gin.Context) {
	var input struct {
		Smtp service.SmtpSetting `json:"smtp"`
		To   string              `json:"to" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "收件人邮箱不能为空")
		return
	}

	if err := service.TestSmtp(input.Smtp, input.To); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "测试邮件发送成功")
}
