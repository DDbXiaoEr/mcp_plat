package handler

// Author: deepseek-v4-pro / opencode

import (
	"net/http"
	"strconv"

	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type OverviewHandler struct{}

func NewOverviewHandler() *OverviewHandler {
	return &OverviewHandler{}
}

func (h *OverviewHandler) Stats(c *gin.Context) {
	stats, err := service.GetOverviewStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": stats})
}

func (h *OverviewHandler) CallTrend(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	trend, err := service.GetCallTrend(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": trend})
}
