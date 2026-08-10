package database

// Author: deepseek-v4-pro / opencode

import (
	"fmt"
	"log"

	"mcp_plat-console/auditstore"
	"mcp_plat-console/config"
	"mcp_plat-console/model"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var AuditStore auditstore.Store

func Init() {
	cfg := config.AppConfig

	DB = openDB(cfg.Database)
	err := DB.AutoMigrate(
		&model.User{},
		&model.AccessKey{},
		&model.UsageHistory{},
		&model.MCPServer{},
		&model.Role{},
		&model.RoleServer{},
		&model.Setting{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if cfg.AuditLogDB.Type != "" {
		AuditStore, err = auditstore.New(cfg.AuditLogDB)
		if err != nil {
			log.Fatalf("failed to init audit store: %v", err)
		}
	}
}

func openDB(cfg config.DatabaseConfig) *gorm.DB {
	var dialector gorm.Dialector

	switch cfg.Type {
	case "postgres":
		pg := cfg.Postgres
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			pg.Host, pg.User, pg.Password, pg.DBName, pg.Port)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.SQLite.Path)
	default:
		log.Fatalf("unsupported database type: %s", cfg.Type)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}
