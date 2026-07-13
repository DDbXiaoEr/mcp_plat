package database

import (
	"fmt"
	"log"

	"mcp_plat-console/config"
	"mcp_plat-console/model"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	cfg := config.AppConfig

	var dialector gorm.Dialector

	switch cfg.Database.Type {
	case "postgres":
		pg := cfg.Database.Postgres
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			pg.Host, pg.User, pg.Password, pg.DBName, pg.Port)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.SQLite.Path)
	default:
		log.Fatalf("unsupported database type: %s", cfg.Database.Type)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = DB.AutoMigrate(
		&model.User{},
		&model.AccessKey{},
		&model.UsageHistory{},
		&model.MCPServer{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
