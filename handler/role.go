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

type RoleHandler struct{}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{}
}

func (h *RoleHandler) List(c *gin.Context) {
	roles, err := service.ListRoles()
	if err != nil {
		resp.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}

	resp.OK(c, "success", roles)
}

func (h *RoleHandler) Create(c *gin.Context) {
	var input service.CreateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	role, err := service.CreateRole(input)
	if err != nil {
		resp.Fail(c, http.StatusInternalServerError, "创建失败")
		return
	}

	resp.OK(c, "创建成功", role)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	var input service.UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.UpdateRole(uint(id), input); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "更新成功")
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.DeleteRole(uint(id)); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "删除成功")
}

func (h *RoleHandler) AssignUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	var input service.AssignUsersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := service.AssignUsersToRole(uint(id), input); err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "分配成功")
}

func (h *RoleHandler) GetUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	users, err := service.GetRoleUsers(uint(id))
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	resp.OK(c, "success", users)
}
