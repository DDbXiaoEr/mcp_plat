package config

// Author: deepseek-v4-pro / opencode

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database        DatabaseConfig `yaml:"database"`
	AuditLogDB      DatabaseConfig `yaml:"audit_log_db"`
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
		log.Fatalf("config: %s not found", configPath)
	}

	if err := yaml.Unmarshal(data, AppConfig); err != nil {
		log.Fatalf("config: failed to parse %s: %v", configPath, err)
	}

	applySoftDefaults()
	validateRequired()
}

func applySoftDefaults() {
	if AppConfig.Database.Type == "" {
		AppConfig.Database.Type = "sqlite"
	}
	if AppConfig.Database.SQLite.Path == "" {
		AppConfig.Database.SQLite.Path = "data.db"
	}
	if AppConfig.AuditLogDB.Type == "" {
		AppConfig.AuditLogDB.Type = "sqlite"
	}
	if AppConfig.AuditLogDB.SQLite.Path == "" {
		AppConfig.AuditLogDB.SQLite.Path = "audit_log.db"
	}
	if AppConfig.ServerPort == "" {
		AppConfig.ServerPort = "8080"
	}
}

func validateRequired() {
	if AppConfig.JWTSecret == "" {
		log.Fatal("config: jwt_secret is required in config.yaml")
	}
	if len(AppConfig.JWTSecret) < 32 {
		log.Fatal("config: jwt_secret must be at least 32 characters")
	}
	if AppConfig.AccessKeySecret == "" {
		log.Fatal("config: access_key_secret is required in config.yaml")
	}
	if len(AppConfig.AccessKeySecret) < 32 {
		log.Fatal("config: access_key_secret must be at least 32 characters")
	}
	if AppConfig.Admin.Username == "" {
		log.Fatal("config: admin.username is required in config.yaml")
	}
	if AppConfig.Admin.Password == "" {
		log.Fatal("config: admin.password is required in config.yaml")
	}
}

func (c *Config) AdminUsername() string { return c.Admin.Username }
func (c *Config) AdminPassword() string { return c.Admin.Password }
