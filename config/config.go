package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	ServerPort    string
	AdminUsername string
	AdminPassword string
}

var AppConfig *Config

func Load() {
	_ = godotenv.Load()

	AppConfig = &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", "123456"),
		DBName:        getEnv("DB_NAME", "mcp_platform"),
		JWTSecret:     getEnv("JWT_SECRET", "mcp-platform-secret-key"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
	}
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost + " user=" + c.DBUser + " password=" + c.DBPassword + " dbname=" + c.DBName + " port=" + c.DBPort + " sslmode=disable TimeZone=Asia/Shanghai"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
