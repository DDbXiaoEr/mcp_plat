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
	"regexp"
	"strings"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/model"
	"mcp_plat-console/session"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginOutput struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	SessionID string `json:"-"`
}

type CASValidateInput struct {
	Ticket     string `json:"ticket" binding:"required"`
	ServiceUrl string `json:"serviceUrl" binding:"required"`
}

type authSettings struct {
	Method string     `json:"method"`
	Cas    CasConfig  `json:"cas"`
	Ldap   LdapConfig `json:"ldap"`
}

type RoleAssignRule struct {
	RoleID  uint   `json:"roleId"`
	Pattern string `json:"pattern"`
}

type userOpsSettings struct {
	FilterAttribute string           `json:"filterAttribute"`
	RoleRules       []RoleAssignRule `json:"roleRules"`
	MaxAccessKeys   int              `json:"maxAccessKeys"`
	AccessKeyCron   string           `json:"accessKeyCron"`
}

func getAutoAssignRoleID(attrs map[string]string) *uint {
	var ops userOpsSettings
	if err := GetSetting("user_ops", &ops); err != nil {
		return nil
	}
	value := attrs[ops.FilterAttribute]
	for _, rule := range ops.RoleRules {
		if rule.Pattern == "" {
			continue
		}
		matched, err := regexp.MatchString(rule.Pattern, value)
		if err == nil && matched {
			roleID := rule.RoleID
			return &roleID
		}
	}
	return nil
}

func getMaxAccessKeys() int {
	var ops userOpsSettings
	if err := GetSetting("user_ops", &ops); err == nil && ops.MaxAccessKeys > 0 {
		return ops.MaxAccessKeys
	}
	return 0
}

func getAccessKeyCron() string {
	var ops userOpsSettings
	if err := GetSetting("user_ops", &ops); err == nil {
		return ops.AccessKeyCron
	}
	return ""
}

func Login(input LoginInput) (*LoginOutput, error) {
	cfg := config.AppConfig

	if cfg.AdminEnabled() && input.Username == cfg.Admin.Username && input.Password == cfg.Admin.Password {
		return issueSession(session.Data{
			UserID:   0,
			Username: cfg.Admin.Username,
			Name:     cfg.Admin.Username,
			Role:     "admin",
		})
	}

	var authCfg authSettings
	authErr := GetSetting("auth", &authCfg)

	if !cfg.AdminEnabled() {
		// 服务节点：仅使用平台配置的认证方式，未配置时直接返回登录错误
		if authErr != nil || authCfg.Method == "" {
			return nil, errors.New("用户名或密码错误")
		}

		switch authCfg.Method {
		case "local":
			output, err := loginLocal(input)
			if err != nil {
				return nil, err
			}
			if output != nil {
				return output, nil
			}
		case "ldap":
			return loginLDAP(authCfg, input)
		}

		return nil, errors.New("用户名或密码错误")
	}

	output, err := loginLocal(input)
	if err != nil {
		return nil, err
	}
	if output != nil {
		return output, nil
	}

	if authCfg.Method == "ldap" {
		return loginLDAP(authCfg, input)
	}

	return nil, errors.New("用户名或密码错误")
}

// loginLocal 使用平台本地用户密码认证；用户不存在或密码不匹配返回 (nil, nil)。
func loginLocal(input LoginInput) (*LoginOutput, error) {
	var user model.User
	err := database.DB.Where("username = ?", input.Username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status == 0 {
		return nil, errors.New("该用户已被禁用，请联系管理员")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, nil
	}

	return issueUserSession(user)
}

func loginLDAP(authCfg authSettings, input LoginInput) (*LoginOutput, error) {
	var user model.User
	userErr := database.DB.Where("username = ?", input.Username).First(&user).Error
	userNotFound := errors.Is(userErr, gorm.ErrRecordNotFound)

	if userErr != nil && !userNotFound {
		return nil, errors.New("用户名或密码错误")
	}

	if !userNotFound && user.Status == 0 {
		return nil, errors.New("该用户已被禁用，请联系管理员")
	}

	ldapAttrs, err := ldapAuthenticate(authCfg.Ldap, input.Username, input.Password)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	if userNotFound {
		uid := ldapAttrs["uid"]
		if uid == "" {
			uid = input.Username
		}
		user = model.User{
			UID:      input.Username,
			Username: input.Username,
			Password: string(hashedPassword),
			RoleID: getAutoAssignRoleID(map[string]string{
				"username":     input.Username,
				"uid":          uid,
				"name":         ldapAttrs["name"],
				"email":        ldapAttrs["email"],
				"phone":        ldapAttrs["phone"],
				"organization": ldapAttrs["organization"],
			}),
			Name:         ldapAttrs["name"],
			Email:        ldapAttrs["email"],
			Phone:        ldapAttrs["phone"],
			Organization: ldapAttrs["organization"],
		}
		if err := database.DB.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}
	} else {
		updates := map[string]interface{}{"password": string(hashedPassword)}
		if v, ok := ldapAttrs["name"]; ok && v != "" {
			updates["name"] = v
			user.Name = v
		}
		if v, ok := ldapAttrs["uid"]; ok && v != "" {
			updates["uid"] = v
			user.UID = v
		} else if user.UID != input.Username {
			updates["uid"] = input.Username
			user.UID = input.Username
		}
		if v, ok := ldapAttrs["email"]; ok && v != "" {
			updates["email"] = v
		}
		if v, ok := ldapAttrs["phone"]; ok && v != "" {
			updates["phone"] = v
		}
		if v, ok := ldapAttrs["organization"]; ok && v != "" {
			updates["organization"] = v
		}
		database.DB.Model(&user).Updates(updates)
	}

	return issueUserSession(user)
}

func CASLogin(input CASValidateInput) (*LoginOutput, error) {
	var authCfg authSettings
	if err := GetSetting("auth", &authCfg); err != nil {
		return nil, errors.New("CAS 配置未找到")
	}

	if authCfg.Method != "cas" {
		return nil, errors.New("CAS 认证未启用")
	}

	if authCfg.Cas.ServerUrl == "" {
		return nil, errors.New("CAS 服务地址未配置")
	}

	username, casAttrs, err := casValidateTicket(authCfg.Cas.ServerUrl, input.ServiceUrl, input.Ticket)
	if err != nil {
		return nil, err
	}

	attrs := map[string]string{"username": username}
	for platformField, casAttr := range authCfg.Cas.AttrMapping {
		if v, ok := casAttrs[casAttr]; ok && v != "" {
			attrs[platformField] = v
		}
	}

	var user model.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		uid := attrs["uid"]
		if uid == "" {
			uid = username
		}
		user = model.User{
			UID:      uid,
			Username: username,
			Password: "",
			RoleID: getAutoAssignRoleID(map[string]string{
				"username":     username,
				"uid":          uid,
				"name":         attrs["name"],
				"email":        attrs["email"],
				"phone":        attrs["phone"],
				"organization": attrs["organization"],
			}),
			Name:         attrs["name"],
			Email:        attrs["email"],
			Phone:        attrs["phone"],
			Organization: attrs["organization"],
		}
		if err := database.DB.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}
	} else {
		if user.Status == 0 {
			return nil, errors.New("该用户已被禁用，请联系管理员")
		}
		updates := map[string]interface{}{}
		if v, ok := attrs["name"]; ok && v != "" {
			updates["name"] = v
			user.Name = v
		}
		if v, ok := attrs["uid"]; ok && v != "" {
			updates["uid"] = v
			user.UID = v
		} else if user.UID != username {
			updates["uid"] = username
			user.UID = username
		}
		if v, ok := attrs["email"]; ok && v != "" {
			updates["email"] = v
		}
		if v, ok := attrs["phone"]; ok && v != "" {
			updates["phone"] = v
		}
		if v, ok := attrs["organization"]; ok && v != "" {
			updates["organization"] = v
		}
		if len(updates) > 0 {
			database.DB.Model(&user).Updates(updates)
		}
	}

	return issueUserSession(user)
}

type AuthMethodOutput struct {
	Method string    `json:"method"`
	Cas    CasConfig `json:"cas"`
}

func GetAuthMethod() (*AuthMethodOutput, error) {
	var authCfg authSettings
	if err := GetSetting("auth", &authCfg); err != nil {
		return &AuthMethodOutput{Method: "local"}, nil
	}
	if authCfg.Method == "" {
		authCfg.Method = "local"
	}

	return &AuthMethodOutput{
		Method: authCfg.Method,
		Cas:    authCfg.Cas,
	}, nil
}

func GetProfile(userID uint) (*model.User, error) {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

type UpdateProfileInput struct {
	Email string `json:"email"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func UpdateProfile(userID uint, input UpdateProfileInput) (*model.User, error) {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	email := strings.TrimSpace(input.Email)
	if email == "" {
		return nil, errors.New("邮箱不能为空")
	}
	if !emailRegex.MatchString(email) {
		return nil, errors.New("邮箱格式不正确")
	}
	if len(email) > 128 {
		return nil, errors.New("邮箱长度超出限制")
	}

	var dup model.User
	if err := database.DB.Where("email = ? AND id != ?", email, userID).First(&dup).Error; err == nil {
		return nil, errors.New("该邮箱已被其他账号使用")
	}

	if err := database.DB.Model(&user).Update("email", email).Error; err != nil {
		return nil, errors.New("保存失败，请稍后重试")
	}
	user.Email = email

	return &user, nil
}

func issueUserSession(user model.User) (*LoginOutput, error) {
	name := user.Name
	if name == "" {
		name = user.Username
	}
	return issueSession(session.Data{
		UserID:   user.ID,
		Username: user.Username,
		Name:     name,
		Role:     "user",
	})
}

func issueSession(data session.Data) (*LoginOutput, error) {
	id, err := session.Create(data)
	if err != nil {
		return nil, errors.New("创建会话失败")
	}
	return &LoginOutput{
		Username:  data.Username,
		Name:      data.Name,
		Role:      data.Role,
		SessionID: id,
	}, nil
}
