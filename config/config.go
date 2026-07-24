package config

import (
	"fmt"
	"log"
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
		fmt.Printf("config: %s not found, using env vars only\n", configPath)
	} else {
		if err := yaml.Unmarshal(data, AppConfig); err != nil {
			log.Fatalf("config: failed to parse %s: %v", configPath, err)
		}
	}

	applyEnvOverrides()
	applySoftDefaults()
	validateRequired()
}

func applyEnvOverrides() {
	if v := os.Getenv("JWT_SECRET"); v != "" {
		AppConfig.JWTSecret = v
	}
	if v := os.Getenv("ACCESS_KEY_SECRET"); v != "" {
		AppConfig.AccessKeySecret = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		AppConfig.ServerPort = v
	}
	if v := os.Getenv("ADMIN_USERNAME"); v != "" {
		AppConfig.Admin.Username = v
	}
	if v := os.Getenv("ADMIN_PASSWORD"); v != "" {
		AppConfig.Admin.Password = v
	}
	if v := os.Getenv("DB_TYPE"); v != "" {
		AppConfig.Database.Type = v
	}
	if v := os.Getenv("DB_POSTGRES_HOST"); v != "" {
		AppConfig.Database.Postgres.Host = v
	}
	if v := os.Getenv("DB_POSTGRES_PORT"); v != "" {
		AppConfig.Database.Postgres.Port = v
	}
	if v := os.Getenv("DB_POSTGRES_USER"); v != "" {
		AppConfig.Database.Postgres.User = v
	}
	if v := os.Getenv("DB_POSTGRES_PASSWORD"); v != "" {
		AppConfig.Database.Postgres.Password = v
	}
	if v := os.Getenv("DB_POSTGRES_DBNAME"); v != "" {
		AppConfig.Database.Postgres.DBName = v
	}
	if v := os.Getenv("DB_SQLITE_PATH"); v != "" {
		AppConfig.Database.SQLite.Path = v
	}
}

func applySoftDefaults() {
	if AppConfig.Database.Type == "" {
		AppConfig.Database.Type = "sqlite"
	}
	if AppConfig.Database.SQLite.Path == "" {
		AppConfig.Database.SQLite.Path = "data.db"
	}
	if AppConfig.ServerPort == "" {
		AppConfig.ServerPort = "8080"
	}
}

func validateRequired() {
	if AppConfig.JWTSecret == "" {
		log.Fatal("config: JWT_SECRET is required, set via config.yaml or JWT_SECRET env var")
	}
	if len(AppConfig.JWTSecret) < 32 {
		log.Fatal("config: JWT_SECRET must be at least 32 characters")
	}
	if AppConfig.AccessKeySecret == "" {
		log.Fatal("config: ACCESS_KEY_SECRET is required, set via config.yaml or ACCESS_KEY_SECRET env var")
	}
	if len(AppConfig.AccessKeySecret) < 32 {
		log.Fatal("config: ACCESS_KEY_SECRET must be at least 32 characters")
	}
	if AppConfig.Admin.Username == "" {
		log.Fatal("config: ADMIN_USERNAME is required, set via config.yaml or ADMIN_USERNAME env var")
	}
	if AppConfig.Admin.Password == "" {
		log.Fatal("config: ADMIN_PASSWORD is required, set via config.yaml or ADMIN_PASSWORD env var")
	}
}

func (c *Config) AdminUsername() string { return c.Admin.Username }
func (c *Config) AdminPassword() string { return c.Admin.Password }
