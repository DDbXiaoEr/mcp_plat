// Copyright (C) 2026 Zhaoquan Wang
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package database

// Author: deepseek-v4-pro / opencode

import (
	"fmt"
	"log"

	"mcp_plat-console/auditstore"
	"mcp_plat-console/config"
	"mcp_plat-console/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
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
	case "mysql":
		my := cfg.MySQL
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			my.User, my.Password, my.Host, my.Port, my.DBName)
		dialector = mysql.Open(dsn)
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
