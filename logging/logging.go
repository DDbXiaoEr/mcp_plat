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

package logging

// Author: deepseek-v4-pro / opencode

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
)

type logSettings struct {
	SyslogEnabled  bool   `json:"syslogEnabled"`
	LogPath        string `json:"logPath"`
	LogLevel       string `json:"logLevel"`
	LogPrefix      string `json:"logPrefix"`
	SyslogHost     string `json:"syslogHost"`
	SyslogPort     int    `json:"syslogPort"`
	SyslogProtocol string `json:"syslogProtocol"`
}

type syslogWriter struct {
	mu       sync.Mutex
	network  string
	addr     string
	tag      string
	hostname string
	conn     net.Conn
}

func (w *syslogWriter) connect() error {
	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
	}

	conn, err := net.DialTimeout(w.network, w.addr, 5*time.Second)
	if err != nil {
		return err
	}
	w.conn = conn
	return nil
}

func (w *syslogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	msg := fmt.Sprintf("<14>%s %s %s: %s",
		time.Now().Format(time.Stamp), w.hostname, w.tag, string(p))
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}

	if w.conn != nil {
		if _, err := w.conn.Write([]byte(msg)); err == nil {
			return len(p), nil
		}
	}

	if err := w.connect(); err != nil {
		return 0, err
	}
	if _, err := w.conn.Write([]byte(msg)); err != nil {
		return 0, err
	}
	return len(p), nil
}

func Setup() {
	var settings logSettings
	if err := service.GetSetting("log", &settings); err != nil {
		return
	}

	if !settings.SyslogEnabled || settings.SyslogHost == "" {
		return
	}

	network := "tcp"
	if settings.SyslogProtocol == "udp" {
		network = "udp"
	}
	port := settings.SyslogPort
	if port <= 0 {
		port = 514
	}
	tag := settings.LogPrefix
	if tag == "" {
		tag = "mcp-plat"
	}
	hostname, _ := os.Hostname()

	w := &syslogWriter{
		network:  network,
		addr:     net.JoinHostPort(settings.SyslogHost, strconv.Itoa(port)),
		tag:      tag,
		hostname: hostname,
	}

	if err := w.connect(); err != nil {
		log.Printf("syslog: failed to connect %s://%s: %v, fallback to stdout", network, w.addr, err)
		return
	}

	log.Printf("syslog: redirecting logs to %s://%s", network, w.addr)
	log.SetOutput(w)
	gin.DefaultWriter = w
	gin.DefaultErrorWriter = w
}
