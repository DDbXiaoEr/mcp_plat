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

package middleware

// Author: deepseek-v4-pro / opencode

import (
	"mcp_plat-console/i18n"

	"github.com/gin-gonic/gin"
)

// Locale resolves the client language from the Accept-Language header and
// stores it on the context so responses (and any data localization) can pick it
// up. Register it globally, before the auth middlewares, in router.Setup.
func Locale() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("locale", i18n.Lang(c.GetHeader("Accept-Language")))
		c.Next()
	}
}
