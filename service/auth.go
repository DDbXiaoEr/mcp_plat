package service

import (
	"errors"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/middleware"
	"mcp_plat-console/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

func Login(input LoginInput) (*LoginOutput, error) {
	cfg := config.AppConfig

	if input.Username == cfg.AdminUsername && input.Password == cfg.AdminPassword {
		token, err := generateAdminToken()
		if err != nil {
			return nil, err
		}
		return &LoginOutput{
			Token:    token,
			Username: cfg.AdminUsername,
			Role:     "admin",
		}, nil
	}

	var user model.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
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
		Username: config.AppConfig.AdminUsername,
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
