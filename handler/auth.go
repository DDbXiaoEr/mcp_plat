package handler

import (
	"net/http"

	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c *gin.Context) {
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
