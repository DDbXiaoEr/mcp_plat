package handler

import (
	"net/http"
	"sync"
	"time"

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
		c.JSON(http.StatusTooManyRequests, gin.H{"code": 429, "message": "登录请求过于频繁，请稍后再试"})
		return
	}

	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	output, err := service.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "登录成功", "data": output})
}

func (h *AuthHandler) CASValidate(c *gin.Context) {
	var input service.CASValidateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	output, err := service.CASLogin(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "登录成功", "data": output})
}

func (h *AuthHandler) GetAuthMethod(c *gin.Context) {
	output, err := service.GetAuthMethod()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": output})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	if userID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": gin.H{
			"username": c.GetString("username"),
			"role":     role,
		}})
		return
	}

	user, err := service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": user})
}
