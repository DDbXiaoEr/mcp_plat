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
	"errors"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	KeyID  uint   `json:"key_id"`
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func ListAccessKeys(userID uint) ([]model.AccessKey, error) {
	var keys []model.AccessKey
	err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&keys).Error
	if err != nil {
		return nil, err
	}
	now := time.Now()
	for i := range keys {
		if keys[i].ExpiredAt != nil && keys[i].ExpiredAt.Before(now) {
			keys[i].IsExpired = true
		}
	}
	return keys, nil
}

func CreateAccessKey(userID uint, input CreateAccessKeyInput) (*model.AccessKey, error) {
	var count int64
	database.DB.Model(&model.AccessKey{}).Where("user_id = ?", userID).Count(&count)
	limit := getMaxAccessKeys()
	if limit > 0 && int(count) >= limit {
		return nil, fmt.Errorf("AccessKey 数量已达上限（%d 个），请先删除再创建", limit)
	}

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

	// 先落库拿到自增 ID，再签发包含 key_id 的 key（审计日志用 ID 存储，减少数据量）
	ak := model.AccessKey{
		UserID:    userID,
		Name:      input.Name,
		Key:       "pending-" + uuid.NewString(),
		Enabled:   true,
		ExpiredAt: input.ExpiredAt,
		Servers:   input.Servers,
	}

	if err := database.DB.Create(&ak).Error; err != nil {
		return nil, err
	}

	key, err := generateAccessKey(ak.ID, userID, roleName, input.ExpiredAt)
	if err != nil {
		return nil, err
	}
	if err := database.DB.Model(&ak).Update("key", key).Error; err != nil {
		return nil, err
	}
	ak.Key = key
	return &ak, nil
}

func UpdateAccessKey(id, userID uint, input UpdateAccessKeyInput) error {
	var ak model.AccessKey
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&ak).Error; err != nil {
		return errors.New("AccessKey 不存在")
	}

	if input.Enabled != nil && *input.Enabled && ak.ExpiredAt != nil && ak.ExpiredAt.Before(time.Now()) {
		return errors.New("该 AccessKey 已过期，无法启用，请删除后重新创建")
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

func generateAccessKey(keyID, userID uint, role string, expiredAt *time.Time) (string, error) {
	claims := AccessKeyClaims{
		KeyID:  keyID,
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

func DisableExpiredAccessKeys() {
	var keys []model.AccessKey
	if err := database.DB.
		Where("enabled = ? AND expired_at IS NOT NULL AND expired_at < ?", true, time.Now()).
		Find(&keys).Error; err != nil {
		log.Printf("access_key scheduler: query expired keys failed: %v", err)
		return
	}
	if len(keys) == 0 {
		return
	}

	ids := make([]uint, 0, len(keys))
	for _, k := range keys {
		ids = append(ids, k.ID)
	}
	if err := database.DB.Model(&model.AccessKey{}).Where("id IN ?", ids).Update("enabled", false).Error; err != nil {
		log.Printf("access_key scheduler: disable expired keys failed: %v", err)
		return
	}
	log.Printf("access_key scheduler: disabled %d expired key(s)", len(keys))

	notifyExpiredAccessKeys(keys)
}

func notifyExpiredAccessKeys(keys []model.AccessKey) {
	byUser := make(map[uint][]model.AccessKey)
	for _, k := range keys {
		byUser[k.UserID] = append(byUser[k.UserID], k)
	}

	platform := GetPlatform().Name
	for userID, list := range byUser {
		var user model.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			log.Printf("access_key scheduler: load user %d failed: %v", userID, err)
			continue
		}
		if user.Email == "" {
			continue
		}

		var rows strings.Builder
		for _, k := range list {
			name := html.EscapeString(k.Name)
			exp := "永不过期"
			if k.ExpiredAt != nil {
				exp = k.ExpiredAt.Format("2006-01-02 15:04")
			}
			rows.WriteString(fmt.Sprintf("<li><strong>%s</strong>（到期时间：%s）</li>", name, html.EscapeString(exp)))
		}

		subject := fmt.Sprintf("%s：您的 AccessKey 已因过期被禁用", platform)
		body := fmt.Sprintf(
			"<p>您好，您有以下 AccessKey 因超过有效期已被自动禁用：</p><ul>%s</ul><p>请登录平台，在「AccessKey 管理」中删除这些已过期的 Key。如需继续使用，请重新创建。</p>",
			rows.String(),
		)

		if err := SendNotificationMail(user.Email, subject, body); err != nil {
			log.Printf("access_key scheduler: notify user %d (%s) failed: %v", userID, user.Email, err)
		}
	}
}
