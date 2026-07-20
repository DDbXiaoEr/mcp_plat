package service

import (
	"errors"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/model"

	"github.com/golang-jwt/jwt/v5"
)

type CreateAccessKeyInput struct {
	Name      string     `json:"name" binding:"required"`
	Servers   string     `json:"servers"`
	ExpiredAt *time.Time `json:"expired_at"`
}

type UpdateAccessKeyInput struct {
	Name      string     `json:"name"`
	Enabled   *bool      `json:"enabled"`
	Servers   string     `json:"servers"`
	ExpiredAt *time.Time `json:"expired_at"`
}

type AccessKeyClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func ListAccessKeys(userID uint) ([]model.AccessKey, error) {
	var keys []model.AccessKey
	err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&keys).Error
	return keys, err
}

func CreateAccessKey(userID uint, input CreateAccessKeyInput) (*model.AccessKey, error) {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	roleName := ""
	if user.RoleID != nil {
		var role model.Role
		if err := database.DB.First(&role, *user.RoleID).Error; err == nil {
			roleName = role.Name
		}
	}

	key, err := generateAccessKey(userID, roleName, input.ExpiredAt)
	if err != nil {
		return nil, err
	}

	ak := model.AccessKey{
		UserID:    userID,
		Name:      input.Name,
		Key:       key,
		Enabled:   true,
		ExpiredAt: input.ExpiredAt,
		Servers:   input.Servers,
	}

	if err := database.DB.Create(&ak).Error; err != nil {
		return nil, err
	}
	return &ak, nil
}

func UpdateAccessKey(id, userID uint, input UpdateAccessKeyInput) error {
	var ak model.AccessKey
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&ak).Error; err != nil {
		return errors.New("AccessKey 不存在")
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Enabled != nil {
		updates["enabled"] = *input.Enabled
	}
	if input.Servers != "" {
		updates["servers"] = input.Servers
	}
	if input.ExpiredAt != nil {
		updates["expired_at"] = input.ExpiredAt
	}

	return database.DB.Model(&ak).Updates(updates).Error
}

func DeleteAccessKey(id, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.AccessKey{})
	if result.RowsAffected == 0 {
		return errors.New("AccessKey 不存在")
	}
	return result.Error
}

func generateAccessKey(userID uint, role string, expiredAt *time.Time) (string, error) {
	claims := AccessKeyClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	if expiredAt != nil {
		claims.ExpiresAt = jwt.NewNumericDate(*expiredAt)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.AppConfig.AccessKeySecret))
	if err != nil {
		return "", err
	}
	return "ak-" + tokenStr, nil
}
