package service

import (
	"errors"

	"mcp_plat-console/database"
	"mcp_plat-console/model"
)

type CreateServerInput struct {
	Name       string `json:"name" binding:"required"`
	Address    string `json:"address" binding:"required"`
	Department string `json:"department"`
	Protocol   string `json:"protocol"`
	Tools      string `json:"tools"`
}

type UpdateServerInput struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	Department string `json:"department"`
	Protocol   string `json:"protocol"`
	Tools      string `json:"tools"`
}

func ListServers() ([]model.MCPServer, error) {
	var servers []model.MCPServer
	err := database.DB.Order("created_at desc").Find(&servers).Error
	return servers, err
}

func CreateServer(input CreateServerInput) (*model.MCPServer, error) {
	protocol := input.Protocol
	if protocol == "" {
		protocol = "SSE"
	}
	server := model.MCPServer{
		Name:       input.Name,
		Address:    input.Address,
		Department: input.Department,
		Protocol:   protocol,
		Tools:      input.Tools,
	}
	if err := database.DB.Create(&server).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

func UpdateServer(id uint, input UpdateServerInput) error {
	var server model.MCPServer
	if err := database.DB.First(&server, id).Error; err != nil {
		return errors.New("MCP 服务器不存在")
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Address != "" {
		updates["address"] = input.Address
	}
	if input.Department != "" {
		updates["department"] = input.Department
	}
	if input.Protocol != "" {
		updates["protocol"] = input.Protocol
	}
	if input.Tools != "" {
		updates["tools"] = input.Tools
	}

	return database.DB.Model(&server).Updates(updates).Error
}

func DeleteServer(id uint) error {
	result := database.DB.Delete(&model.MCPServer{}, id)
	if result.RowsAffected == 0 {
		return errors.New("MCP 服务器不存在")
	}
	return result.Error
}
