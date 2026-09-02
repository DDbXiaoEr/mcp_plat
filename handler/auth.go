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
	"sync"
	"time"

	"mcp_plat-console/resp"
	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
}

func newRateLimiter() *rateLimiter {
	rl := &rateLimiter{attempts: make(map[string][]time.Time)}
	go rl.cleanup(1 * time.Minute)
	return rl
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-1 * time.Minute)

	recent := make([]time.Time, 0)
	for _, t := range rl.attempts[key] {
		if t.After(windowStart) {
			recent = append(recent, t)
		}
	}
	rl.attempts[key] = recent

	if len(recent) >= 10 {
		return false
	}

	rl.attempts[key] = append(rl.attempts[key], now)
	return true
}

func (rl *rateLimiter) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		expire := now.Add(-2 * time.Minute)
		for key, times := range rl.attempts {
			valid := make([]time.Time, 0)
			for _, t := range times {
				if t.After(expire) {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.attempts, key)
			} else {
				rl.attempts[key] = valid
			}
		}
		rl.mu.Unlock()
	}
}

var loginLimiter = newRateLimiter()

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c *gin.Context) {
	clientIP := c.ClientIP()

	if !loginLimiter.allow(clientIP) {
		resp.Fail(c, http.StatusTooManyRequests, "登录请求过于频繁，请稍后再试")
		return
	}

	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	output, err := service.Login(input)
	if err != nil {
		resp.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}

	resp.OK(c, "登录成功", output)
}

func (h *AuthHandler) CASValidate(c *gin.Context) {
	var input service.CASValidateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	output, err := service.CASLogin(input)
	if err != nil {
		resp.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}

	resp.OK(c, "登录成功", output)
}

func (h *AuthHandler) GetPlatform(c *gin.Context) {
	resp.OK(c, "success", service.GetPlatform())
}

func (h *AuthHandler) GetAuthMethod(c *gin.Context) {
	output, err := service.GetAuthMethod()
	if err != nil {
		resp.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp.OK(c, "success", output)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	if userID == 0 {
		resp.OK(c, "success", gin.H{
			"username": c.GetString("username"),
			"role":     role,
		})
		return
	}

	user, err := service.GetProfile(userID)
	if err != nil {
		resp.Fail(c, http.StatusNotFound, err.Error())
		return
	}

	resp.OK(c, "success", user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		resp.Fail(c, http.StatusForbidden, "管理员账号不支持修改邮箱")
		return
	}

	var input service.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	user, err := service.UpdateProfile(userID, input)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功", "data": user})
}
