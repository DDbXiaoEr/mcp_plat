package service

import (
	"errors"

	"mcp_plat-console/database"
	"mcp_plat-console/model"
)

type CreateRoleInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ServerIDs   []string `json:"server_ids"`
}

type UpdateRoleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServerIDs   []string `json:"server_ids"`
}

type RoleOutput struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServerIDs   []string `json:"server_ids"`
	UserCount   int    `json:"user_count"`
}

type AssignUsersInput struct {
	UserIDs []uint `json:"user_ids" binding:"required"`
}

func ListRoles() ([]RoleOutput, error) {
	var roles []model.Role
	if err := database.DB.Order("created_at desc").Find(&roles).Error; err != nil {
		return nil, err
	}

	result := make([]RoleOutput, 0, len(roles))
	for _, role := range roles {
		var serverIDs []string
		database.DB.Model(&model.RoleServer{}).
			Where("role_id = ?", role.ID).
			Pluck("server_id", &serverIDs)
		if serverIDs == nil {
			serverIDs = []string{}
		}

		var userCount int64
		database.DB.Model(&model.User{}).Where("role_id = ?", role.ID).Count(&userCount)

		result = append(result, RoleOutput{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			ServerIDs:   serverIDs,
			UserCount:   int(userCount),
		})
	}

	return result, nil
}

func CreateRole(input CreateRoleInput) (*RoleOutput, error) {
	role := model.Role{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		return nil, err
	}

	if len(input.ServerIDs) > 0 {
		for _, serverID := range input.ServerIDs {
			database.DB.Create(&model.RoleServer{
				RoleID:   role.ID,
				ServerID: serverID,
			})
		}
	}

	return &RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		ServerIDs:   input.ServerIDs,
		UserCount:   0,
	}, nil
}

func UpdateRole(id uint, input UpdateRoleInput) error {
	var role model.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		return errors.New("角色不存在")
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if len(updates) > 0 {
		if err := database.DB.Model(&role).Updates(updates).Error; err != nil {
			return err
		}
	}

	if input.ServerIDs != nil {
		database.DB.Where("role_id = ?", id).Delete(&model.RoleServer{})

		for _, serverID := range input.ServerIDs {
			database.DB.Create(&model.RoleServer{
				RoleID:   id,
				ServerID: serverID,
			})
		}
	}

	return nil
}

func DeleteRole(id uint) error {
	var userCount int64
	database.DB.Model(&model.User{}).Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		return errors.New("该角色下存在用户，请先解除用户绑定")
	}

	database.DB.Where("role_id = ?", id).Delete(&model.RoleServer{})

	result := database.DB.Delete(&model.Role{}, id)
	if result.RowsAffected == 0 {
		return errors.New("角色不存在")
	}

	return result.Error
}

func AssignUsersToRole(roleID uint, input AssignUsersInput) error {
	var role model.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		return errors.New("角色不存在")
	}

	if err := database.DB.Model(&model.User{}).Where("role_id = ?", roleID).Update("role_id", nil).Error; err != nil {
		return err
	}

	if len(input.UserIDs) > 0 {
		if err := database.DB.Model(&model.User{}).Where("id IN ?", input.UserIDs).Update("role_id", roleID).Error; err != nil {
			return err
		}
	}

	return nil
}

func GetRoleUsers(roleID uint) ([]model.User, error) {
	var role model.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		return nil, errors.New("角色不存在")
	}

	var users []model.User
	if err := database.DB.Where("role_id = ?", roleID).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
