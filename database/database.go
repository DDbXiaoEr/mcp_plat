package database

import (
	"log"

	"mcp_plat-console/config"
	"mcp_plat-console/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	dsn := config.AppConfig.DSN()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = DB.AutoMigrate(
		&model.User{},
		&model.AccessKey{},
		&model.UsageHistory{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
