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
	"strconv"

	"mcp_plat-console/resp"
	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type AccessKeyHandler struct{}

func NewAccessKeyHandler() *AccessKeyHandler {
	return &AccessKeyHandler{}
}

func (h *AccessKeyHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	keys, err := service.ListAccessKeys(userID)
	if err != nil {
		resp.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}

	resp.OK(c, "success", keys)
}

func (h *AccessKeyHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input service.CreateAccessKeyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	ak, err := service.CreateAccessKey(userID, input)
	if err != nil {
		resp.Fail(c, http.StatusInternalServerError, "创建失败")
		return
	}

	resp.OK(c, "创建成功", ak)
}

func (h *AccessKeyHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	var input service.UpdateAccessKeyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.UpdateAccessKey(uint(id), userID, input); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "更新成功")
}

func (h *AccessKeyHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.DeleteAccessKey(uint(id), userID); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "删除成功")
}
