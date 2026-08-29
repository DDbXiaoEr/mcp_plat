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

package router

// Author: deepseek-v4-pro / opencode

import (
	"mcp_plat-console/handler"
	"mcp_plat-console/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	authHandler := handler.NewAuthHandler()
	accessKeyHandler := handler.NewAccessKeyHandler()
	historyHandler := handler.NewHistoryHandler()
	auditLogHandler := handler.NewAuditLogHandler()
	serverHandler := handler.NewMCPServerHandler()
	overviewHandler := handler.NewOverviewHandler()
	roleHandler := handler.NewRoleHandler()
	userHandler := handler.NewRBACUserHandler()
	settingHandler := handler.NewSettingHandler()

	r.GET("/healthz", handler.Healthz)
	r.GET("/readyz", handler.Readyz)

	r.POST("/api/auth/login", authHandler.Login)
	r.GET("/api/auth/method", authHandler.GetAuthMethod)
	r.GET("/api/auth/platform", authHandler.GetPlatform)
	r.POST("/api/auth/cas/validate", authHandler.CASValidate)

	auth := r.Group("/api")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/auth/profile", authHandler.Profile)

		auth.GET("/access-keys", accessKeyHandler.List)
		auth.POST("/access-keys", accessKeyHandler.Create)
		auth.PUT("/access-keys/:id", accessKeyHandler.Update)
		auth.DELETE("/access-keys/:id", accessKeyHandler.Delete)

		auth.GET("/history", historyHandler.List)

		auth.GET("/audit-logs", auditLogHandler.List)

		auth.GET("/servers", serverHandler.List)
		auth.GET("/settings/gateway-status", settingHandler.GatewayStatus)
		auth.GET("/settings/:key", settingHandler.GetByKey)
	}

	admin := r.Group("/api")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		admin.POST("/servers", serverHandler.Create)
		admin.POST("/servers/fetch-tools", serverHandler.FetchTools)
		admin.PUT("/servers/:id", serverHandler.Update)
		admin.DELETE("/servers/:id", serverHandler.Delete)
		admin.POST("/servers/publish", serverHandler.Publish)
		admin.POST("/servers/maintenance", serverHandler.Maintenance)

		admin.GET("/roles", roleHandler.List)
		admin.POST("/roles", roleHandler.Create)
		admin.PUT("/roles/:id", roleHandler.Update)
		admin.DELETE("/roles/:id", roleHandler.Delete)
		admin.GET("/roles/:id/users", roleHandler.GetUsers)
		admin.PUT("/roles/:id/users", roleHandler.AssignUsers)

		admin.GET("/users", userHandler.List)
		admin.POST("/users", userHandler.Create)
		admin.PUT("/users/:id", userHandler.Update)
		admin.DELETE("/users/:id", userHandler.Delete)

		admin.GET("/overview/stats", overviewHandler.Stats)
		admin.GET("/overview/call-trend", overviewHandler.CallTrend)

		admin.GET("/settings", settingHandler.Get)
		admin.PUT("/settings/:key", settingHandler.Save)
		admin.POST("/settings/test-ldap", settingHandler.TestLdapMapping)
		admin.POST("/settings/test-smtp", settingHandler.TestSmtp)
	}

	return r
}
