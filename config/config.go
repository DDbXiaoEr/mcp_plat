package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database        DatabaseConfig `yaml:"database"`
	JWTSecret       string         `yaml:"jwt_secret"`
	AccessKeySecret string         `yaml:"access_key_secret"`
	ServerPort      string         `yaml:"server_port"`
	Admin           AdminConfig    `yaml:"admin"`
}

type DatabaseConfig struct {
	Type     string       `yaml:"type"`
	SQLite   SQLiteConfig `yaml:"sqlite"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type SQLiteConfig struct {
	Path string `yaml:"path"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type AdminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var AppConfig *Config

func Load() {
	AppConfig = &Config{}

	configPath := "config.yaml"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		configPath = p
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("config: %s not found, using defaults\n", configPath)
		setDefaults()
		return
	}

	if err := yaml.Unmarshal(data, AppConfig); err != nil {
		fmt.Printf("config: failed to parse %s: %v, using defaults\n", configPath, err)
		setDefaults()
		return
	}

	if AppConfig.Database.Type == "" {
		AppConfig.Database.Type = "sqlite"
	}
	if AppConfig.Database.SQLite.Path == "" {
		AppConfig.Database.SQLite.Path = "data.db"
	}
	if AppConfig.JWTSecret == "" {
		AppConfig.JWTSecret = "mcp-platform-secret-key"
	}
	if AppConfig.AccessKeySecret == "" {
		// 请修改为自定义密钥
		AppConfig.AccessKeySecret = "a8k3x9m2p7q1r6w4v5y0b3n8t2h7j1k5"
	}
	if AppConfig.ServerPort == "" {
		AppConfig.ServerPort = "8080"
	}
	if AppConfig.Admin.Username == "" {
		AppConfig.Admin.Username = "admin"
	}
	if AppConfig.Admin.Password == "" {
		AppConfig.Admin.Password = "admin123"
	}
}

func setDefaults() {
	AppConfig = &Config{
		Database: DatabaseConfig{
			Type: "sqlite",
			SQLite: SQLiteConfig{
				Path: "data.db",
			},
		},
		JWTSecret:       "mcp-platform-secret-key",
		AccessKeySecret: "a8k3x9m2p7q1r6w4v5y0b3n8t2h7j1k5", // 请修改为自定义密钥
		ServerPort:      "8080",
		Admin: AdminConfig{
			Username: "admin",
			Password: "admin123",
		},
	}
}

func (c *Config) AdminUsername() string { return c.Admin.Username }
func (c *Config) AdminPassword() string { return c.Admin.Password }
