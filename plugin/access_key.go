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

package plugin

// Author: deepseek-v4-pro / opencode

import (
	"fmt"
	"strings"

	"mcp_plat-console/config"

	"github.com/golang-jwt/jwt/v5"
)

type AccessKeyClaims struct {
	KeyID  uint   `json:"key_id"`
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func ParseAccessKey(key string) (*AccessKeyClaims, error) {
	return ParseAccessKeyWithSecret(key, []byte(config.AppConfig.AccessKeySecret))
}

func ParseAccessKeyWithSecret(key string, secret []byte) (*AccessKeyClaims, error) {
	if !strings.HasPrefix(key, "ak-") {
		return nil, fmt.Errorf("access key 格式无效")
	}

	tokenStr := strings.TrimPrefix(key, "ak-")

	claims := &AccessKeyClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("access key 签名方法无效")
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("access key 解析失败: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("access key 无效")
	}

	return claims, nil
}
