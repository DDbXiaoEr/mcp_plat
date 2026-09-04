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

package service

// Author: deepseek-v4-pro / opencode

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"mcp_plat-console/database"
	"mcp_plat-console/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleID   *uint  `json:"role_id"`
}

type UpdateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   *uint  `json:"role_id"`
}

type BatchUserItem struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	UID          string `json:"uid"`
	Phone        string `json:"phone"`
	Organization string `json:"organization"`
}

type BatchCreateUsersInput struct {
	Users  []BatchUserItem `json:"users"`
	RoleID *uint           `json:"role_id"`
}

type UserOutput struct {
	ID        uint      `json:"id"`
	UID       string    `json:"uid"`
	Username  string    `json:"username"`
	RoleID    *uint     `json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ListUsers(roleID *uint, keyword string) ([]UserOutput, error) {
	query := database.DB.Model(&model.User{})

	if roleID != nil {
		query = query.Where("role_id = ?", *roleID)
	}

	if keyword != "" {
		query = query.Where("uid LIKE ?", "%"+keyword+"%")
	}

	var users []model.User
	if err := query.Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}

	var roles []model.Role
	database.DB.Find(&roles)
	roleNameMap := make(map[uint]string, len(roles))
	for _, r := range roles {
		roleNameMap[r.ID] = r.Name
	}

	result := make([]UserOutput, 0, len(users))
	for _, u := range users {
		roleName := ""
		if u.RoleID != nil {
			roleName = roleNameMap[*u.RoleID]
		}
		result = append(result, UserOutput{
			ID:        u.ID,
			UID:       u.UID,
			Username:  u.Username,
			RoleID:    u.RoleID,
			RoleName:  roleName,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return result, nil
}

func CreateUser(input CreateUserInput) (*UserOutput, error) {
	var existing model.User
	if err := database.DB.Where("username = ?", input.Username).First(&existing).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	uid := generateUID()

	user := model.User{
		UID:      uid,
		Username: input.Username,
		Password: string(hashedPassword),
		RoleID:   input.RoleID,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	roleName := ""
	if user.RoleID != nil {
		var role model.Role
		if err := database.DB.First(&role, *user.RoleID).Error; err == nil {
			roleName = role.Name
		}
	}

	return &UserOutput{
		ID:        user.ID,
		UID:       user.UID,
		Username:  user.Username,
		RoleID:    user.RoleID,
		RoleName:  roleName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func BatchCreateUsers(input BatchCreateUsersInput) ([]UserOutput, error) {
	if len(input.Users) == 0 {
		return nil, errors.New("用户列表不能为空")
	}

	seenUsername := make(map[string]bool, len(input.Users))
	seenUID := make(map[string]bool)
	seenEmail := make(map[string]bool)
	for _, item := range input.Users {
		item.Username = strings.TrimSpace(item.Username)
		item.UID = strings.TrimSpace(item.UID)
		item.Email = strings.TrimSpace(item.Email)

		if item.Username == "" || strings.TrimSpace(item.Password) == "" {
			return nil, errors.New("用户名和密码不能为空")
		}
		if item.Email != "" && !emailRegex.MatchString(item.Email) {
			return nil, errors.New("邮箱格式不正确")
		}
		if seenUsername[item.Username] {
			return nil, fmt.Errorf("列表中存在重复的用户名 '%s'", item.Username)
		}
		seenUsername[item.Username] = true
		if item.UID != "" {
			if seenUID[item.UID] {
				return nil, fmt.Errorf("学号/工号 '%s' 已被使用", item.UID)
			}
			seenUID[item.UID] = true
		}
		if item.Email != "" {
			if seenEmail[item.Email] {
				return nil, fmt.Errorf("邮箱 '%s' 已被其他账号使用", item.Email)
			}
			seenEmail[item.Email] = true
		}
	}

	results := make([]UserOutput, 0, len(input.Users))
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		roleName := ""
		if input.RoleID != nil {
			var role model.Role
			if err := tx.First(&role, *input.RoleID).Error; err == nil {
				roleName = role.Name
			}
		}

		for _, item := range input.Users {
			item.Username = strings.TrimSpace(item.Username)
			item.UID = strings.TrimSpace(item.UID)
			item.Email = strings.TrimSpace(item.Email)

			var dup model.User
			if err := tx.Where("username = ?", item.Username).First(&dup).Error; err == nil {
				return fmt.Errorf("用户名 '%s' 已存在", item.Username)
			}
			if item.Email != "" {
				if err := tx.Where("email = ? AND email <> ''", item.Email).First(&dup).Error; err == nil {
					return fmt.Errorf("邮箱 '%s' 已被其他账号使用", item.Email)
				}
			}

			uid := item.UID
			if uid == "" {
				for attempt := 0; attempt < 10; attempt++ {
					candidate := generateUID()
					if err := tx.Where("uid = ?", candidate).First(&dup).Error; err != nil {
						uid = candidate
						break
					}
				}
				if uid == "" {
					return errors.New("生成学号/工号失败，请重试")
				}
			} else if err := tx.Where("uid = ?", uid).First(&dup).Error; err == nil {
				return fmt.Errorf("学号/工号 '%s' 已被使用", uid)
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			user := model.User{
				UID:          uid,
				Username:     item.Username,
				Password:     string(hashedPassword),
				Name:         item.Name,
				Email:        item.Email,
				Phone:        item.Phone,
				Organization: item.Organization,
				RoleID:       input.RoleID,
			}
			if err := tx.Create(&user).Error; err != nil {
				return err
			}

			results = append(results, UserOutput{
				ID:        user.ID,
				UID:       user.UID,
				Username:  user.Username,
				RoleID:    user.RoleID,
				RoleName:  roleName,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func UpdateUser(id uint, input UpdateUserInput) error {
	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	updates := map[string]interface{}{}

	if input.Username != "" {
		var dup model.User
		if err := database.DB.Where("username = ? AND id != ?", input.Username, id).First(&dup).Error; err == nil {
			return errors.New("用户名已存在")
		}
		updates["username"] = input.Username
	}
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hashedPassword)
	}

	if input.RoleID != nil {
		if *input.RoleID == 0 {
			updates["role_id"] = nil
		} else {
			updates["role_id"] = *input.RoleID
		}
	}

	if len(updates) > 0 {
		return database.DB.Model(&user).Updates(updates).Error
	}

	return nil
}

func DeleteUser(id uint) error {
	result := database.DB.Delete(&model.User{}, id)
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return result.Error
}

func generateUID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
