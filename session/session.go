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

package session

// Author: deepseek-v4-pro / opencode

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"mcp_plat-console/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const keyPrefix = "mcp:sess:"

var (
	ErrNotFound = errors.New("session not found")
	client      *redis.Client
)

type Data struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func Init() {
	cfg := config.AppConfig.Redis
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("session: failed to connect redis %s: %v", cfg.Addr, err)
	}
	log.Printf("session: redis connected (%s)", cfg.Addr)
}

func Ping() error {
	if client == nil {
		return errors.New("session store not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}

func Create(data Data) (string, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	ttl := ttl()
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		id, err := newID()
		if err != nil {
			return "", err
		}
		ok, err := client.SetNX(ctx, keyPrefix+id, raw, ttl).Result()
		if err != nil {
			return "", err
		}
		if ok {
			return id, nil
		}
	}
	return "", errors.New("failed to allocate session id")
}

func Get(id string) (*Data, error) {
	if id == "" {
		return nil, ErrNotFound
	}
	ctx := context.Background()
	raw, err := client.Get(ctx, keyPrefix+id).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var data Data
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func Touch(id string) error {
	if id == "" {
		return ErrNotFound
	}
	ok, err := client.Expire(context.Background(), keyPrefix+id, ttl()).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func Delete(id string) error {
	if id == "" {
		return nil
	}
	return client.Del(context.Background(), keyPrefix+id).Err()
}

func IDFromRequest(c *gin.Context) string {
	cookie, err := c.Cookie(cookieName())
	if err != nil {
		return ""
	}
	return cookie
}

func SetCookie(c *gin.Context, id string) {
	http.SetCookie(c.Writer, newCookie(id, int(ttl().Seconds())))
}

func ClearCookie(c *gin.Context) {
	http.SetCookie(c.Writer, newCookie("", -1))
}

func newCookie(value string, maxAge int) *http.Cookie {
	cfg := config.AppConfig.Session
	secure := cfg.Secure
	sameSite := sameSiteMode(cfg.SameSite)
	if sameSite == http.SameSiteNoneMode {
		secure = true
	}
	return &http.Cookie{
		Name:     cookieName(),
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	}
}

func cookieName() string {
	name := config.AppConfig.Session.CookieName
	if name == "" {
		return "mcp_session"
	}
	return name
}

func ttl() time.Duration {
	sec := config.AppConfig.Session.MaxAge
	if sec <= 0 {
		sec = 86400
	}
	return time.Duration(sec) * time.Second
}

func sameSiteMode(v string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

func newID() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
