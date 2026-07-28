package handler

// Author: deepseek-v4-pro / opencode

import (
	"net/http"

	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct{}

func NewHistoryHandler() *HistoryHandler {
	return &HistoryHandler{}
}

func (h *HistoryHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	var query service.HistoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	output, err := service.ListHistory(userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": output})
}
