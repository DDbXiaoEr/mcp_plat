package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"mcp_plat-console/database"
	"mcp_plat-console/model"

	"golang.org/x/crypto/bcrypt"
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

type UserOutput struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	RoleID    *uint     `json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ListUsers(roleID *uint) ([]UserOutput, error) {
	query := database.DB.Model(&model.User{})

	if roleID != nil {
		query = query.Where("role_id = ?", *roleID)
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
		Username:  user.Username,
		RoleID:    user.RoleID,
		RoleName:  roleName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
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
