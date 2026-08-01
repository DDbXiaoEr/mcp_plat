package service

// Author: deepseek-v4-pro / opencode

import (
	"errors"
	"fmt"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/middleware"
	"mcp_plat-console/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginOutput struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
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

type userOpsSettings struct {
	DefaultRoleID *uint `json:"defaultRoleId"`
}

func getDefaultRoleID() *uint {
	var ops userOpsSettings
	if err := GetSetting("user_ops", &ops); err == nil {
		return ops.DefaultRoleID
	}
	return nil
}

func Login(input LoginInput) (*LoginOutput, error) {
	cfg := config.AppConfig

	if input.Username == cfg.Admin.Username && input.Password == cfg.Admin.Password {
		token, err := generateAdminToken()
		if err != nil {
			return nil, err
		}
		return &LoginOutput{
			Token:    token,
			Username: cfg.Admin.Username,
			Role:     "admin",
		}, nil
	}

	var authCfg authSettings
	_ = GetSetting("auth", &authCfg)

	var user model.User
	userErr := database.DB.Where("username = ?", input.Username).First(&user).Error
	userNotFound := errors.Is(userErr, gorm.ErrRecordNotFound)

	if userErr != nil && !userNotFound {
		return nil, errors.New("用户名或密码错误")
	}

	if !userNotFound {
		if user.Status == 0 {
			return nil, errors.New("该用户已被禁用，请联系管理员")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err == nil {
			token, err := generateToken(user)
			if err != nil {
				return nil, err
			}
			return &LoginOutput{
				Token:    token,
				Username: user.Username,
				Role:     "user",
			}, nil
		}
	}

	if authCfg.Method == "ldap" {
		ldapAttrs, err := ldapAuthenticate(authCfg.Ldap, input.Username, input.Password)
		if err != nil {
			return nil, errors.New("用户名或密码错误")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("密码加密失败: %v", err)
		}

		if userNotFound {
			user = model.User{
				UID:          generateUID(),
				Username:     input.Username,
				Password:     string(hashedPassword),
				RoleID:       getDefaultRoleID(),
				Email:        ldapAttrs["email"],
				Phone:        ldapAttrs["phone"],
				Organization: ldapAttrs["organization"],
			}
			if err := database.DB.Create(&user).Error; err != nil {
				return nil, fmt.Errorf("创建用户失败: %v", err)
			}
		} else {
			updates := map[string]interface{}{"password": string(hashedPassword)}
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

		token, err := generateToken(user)
		if err != nil {
			return nil, err
		}
		return &LoginOutput{
			Token:    token,
			Username: user.Username,
			Role:     "user",
		}, nil
	}

	return nil, errors.New("用户名或密码错误")
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

	username, err := casValidateTicket(authCfg.Cas.ServerUrl, input.ServiceUrl, input.Ticket)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		user = model.User{
			UID:      generateUID(),
			Username: username,
			Password: "",
			RoleID:   getDefaultRoleID(),
		}
		if err := database.DB.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}
	}

	if user.Status == 0 {
		return nil, errors.New("该用户已被禁用，请联系管理员")
	}

	token, err := generateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token:    token,
		Username: user.Username,
		Role:     "user",
	}, nil
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

func generateAdminToken() (string, error) {
	claims := &middleware.Claims{
		UserID:   0,
		Username: config.AppConfig.Admin.Username,
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func generateToken(user model.User) (string, error) {
	claims := &middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}
