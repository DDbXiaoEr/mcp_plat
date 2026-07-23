package handler

import (
	"net/http"

	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type MCPServerHandler struct{}

func NewMCPServerHandler() *MCPServerHandler {
	return &MCPServerHandler{}
}

func (h *MCPServerHandler) List(c *gin.Context) {
	servers, err := service.ListServers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": servers})
}

func (h *MCPServerHandler) Create(c *gin.Context) {
	var input service.CreateServerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	server, err := service.CreateServer(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": server})
}

func (h *MCPServerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var input service.UpdateServerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := service.UpdateServer(id, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *MCPServerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := service.DeleteServer(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *MCPServerHandler) FetchTools(c *gin.Context) {
	var input service.FetchToolsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	tools, err := service.FetchTools(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取成功", "data": gin.H{"tools": tools}})
}

func (h *MCPServerHandler) Publish(c *gin.Context) {
	var input service.PublishInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := service.PublishServers(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "发布成功"})
}
