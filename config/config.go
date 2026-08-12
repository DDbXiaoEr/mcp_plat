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
	Type       string           `yaml:"type"`
	Postgres   PostgresConfig   `yaml:"postgres"`
	MySQL      MySQLConfig      `yaml:"mysql"`
	ClickHouse ClickHouseConfig `yaml:"clickhouse"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// ClickHouseConfig 审计日志存储引擎（仅 audit_log_db 支持）
type ClickHouseConfig struct {
	Addr     string `yaml:"addr"` // 例如 127.0.0.1:9000
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`       // ClickHouse 库名，默认 default
	Table    string `yaml:"table"`    // 表名，默认 audit_logs
	Engine   string `yaml:"engine"`   // merge_tree（单机）| replicated_merge_tree（集群）
	Cluster  string `yaml:"cluster"`  // 集群名，replicated_merge_tree 必填
	TTLDays  int    `yaml:"ttl_days"` // 数据保留天数，默认 90
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
		AppConfig.Database.Type = "postgres"
	}
	if AppConfig.AuditLogDB.Type == "" {
		AppConfig.AuditLogDB.Type = "postgres"
	}
	if AppConfig.AuditLogDB.Type == "clickhouse" {
		ck := &AppConfig.AuditLogDB.ClickHouse
		if ck.Addr == "" {
			ck.Addr = "127.0.0.1:9000"
		}
		if ck.DB == "" {
			ck.DB = "default"
		}
		if ck.Table == "" {
			ck.Table = "audit_logs"
		}
		if ck.TTLDays <= 0 {
			ck.TTLDays = 90
		}
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

func SetServerPort(port string) {
	if port != "" {
		AppConfig.ServerPort = port
	}
}

func (c *Config) AdminUsername() string { return c.Admin.Username }
func (c *Config) AdminPassword() string { return c.Admin.Password }
