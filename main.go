//go:build !embed

package main

import (
	"fmt"
	"log"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/router"
)

func main() {
	config.Load()

	fmt.Printf("mcp_plat-console version %s, built at %s\n", Version, BuildTime)

	database.Init()

	r := router.Setup()

	log.Printf("server starting on :%s", config.AppConfig.ServerPort)
	if err := r.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
