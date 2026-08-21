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
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"mcp_plat-console/service"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultLogPath   = "./logs"
	defaultLogPrefix = "mcp_plat"
)

type LogSettings struct {
	SyslogEnabled  bool   `json:"syslogEnabled"`
	LogPath        string `json:"logPath"`
	LogLevel       string `json:"logLevel"`
	LogPrefix      string `json:"logPrefix"`
	MaxSize        int    `json:"maxSize"`    // 单文件大小上限（MB）
	MaxBackups     int    `json:"maxBackups"` // 保留的历史日志文件数
	MaxAge         int    `json:"maxAge"`     // 保留天数
	Compress       bool   `json:"compress"`   // 历史日志是否 gzip 压缩
	SyslogHost     string `json:"syslogHost"`
	SyslogPort     int    `json:"syslogPort"`
	SyslogProtocol string `json:"syslogProtocol"`
}

var (
	mu            sync.Mutex
	currentWriter io.Writer        = os.Stdout
	fileLogger    *lumberjack.Logger
)

// proxyWriter 将标准库 log 与 Gin 的输出统一代理到 currentWriter，
// 使保存设置后（Reconfigure）日志目标热切换时对所有输出方即时生效。
type proxyWriter struct{}

func (proxyWriter) Write(p []byte) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	return currentWriter.Write(p)
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

	msg := "<14>" + time.Now().Format(time.Stamp) + " " + w.hostname + " " + w.tag + ": " + string(p)
	if msg[len(msg)-1] != '\n' {
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

func loadSettings() LogSettings {
	settings := LogSettings{
		LogPath:    defaultLogPath,
		LogPrefix:  defaultLogPrefix,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
	}
	if err := service.GetSetting("log", &settings); err != nil {
		return settings
	}
	if settings.LogPath == "" {
		settings.LogPath = defaultLogPath
	}
	if settings.LogPrefix == "" {
		settings.LogPrefix = defaultLogPrefix
	}
	if settings.MaxSize <= 0 {
		settings.MaxSize = 100
	}
	if settings.MaxBackups <= 0 {
		settings.MaxBackups = 10
	}
	if settings.MaxAge <= 0 {
		settings.MaxAge = 30
	}
	return settings
}

// Setup 启动时初始化日志输出。默认输出到标准输出和文件（./logs/mcp_plat.log，带轮转），
// 启用 Syslog 后二者均失效，仅输出到 Syslog 服务器。
func Setup() {
	pw := proxyWriter{}
	log.SetOutput(pw)
	gin.DefaultWriter = pw
	gin.DefaultErrorWriter = pw

	apply(loadSettings())
}

// Reconfigure 重新读取 log 设置并应用，供前端保存日志设置后热生效。
func Reconfigure() {
	apply(loadSettings())
}

func apply(settings LogSettings) {
	w, fl, filePath, syslogErr := buildOutputs(settings)

	if settings.SyslogEnabled && settings.SyslogHost != "" && syslogErr != nil {
		w, fl, filePath, _ = buildOutputs(LogSettings{
			LogPath:    settings.LogPath,
			LogPrefix:  settings.LogPrefix,
			MaxSize:    settings.MaxSize,
			MaxBackups: settings.MaxBackups,
			MaxAge:     settings.MaxAge,
			Compress:   settings.Compress,
		})
	}

	mu.Lock()
	if fileLogger != nil {
		fileLogger.Close()
	}
	fileLogger = fl
	currentWriter = w
	mu.Unlock()

	switch {
	case settings.SyslogEnabled && settings.SyslogHost != "" && syslogErr == nil:
		log.Printf("syslog: redirecting logs to %s://%s", settings.SyslogProtocol, settings.SyslogHost)
	case settings.SyslogEnabled && settings.SyslogHost != "":
		log.Printf("syslog: failed to connect, fallback to stdout+file: %v", syslogErr)
	default:
		log.Printf("logging: stdout + %s enabled", filePath)
	}
}

// buildOutputs 根据设置构造日志目标；Syslog 连接失败时返回 error（由调用方回退）。
func buildOutputs(settings LogSettings) (io.Writer, *lumberjack.Logger, string, error) {
	if settings.SyslogEnabled && settings.SyslogHost != "" {
		w, err := newSyslogWriter(settings)
		if err != nil {
			return nil, nil, "", err
		}
		return w, nil, "", nil
	}

	writers := []io.Writer{os.Stdout}
	filePath := ""
	if settings.LogPath != "" {
		if err := os.MkdirAll(settings.LogPath, 0o755); err != nil {
			log.Printf("logging: failed to create log dir %s: %v", settings.LogPath, err)
		} else {
			filePath = filepath.Join(settings.LogPath, settings.LogPrefix+".log")
			fl := &lumberjack.Logger{
				Filename:   filePath,
				MaxSize:    settings.MaxSize,
				MaxBackups: settings.MaxBackups,
				MaxAge:     settings.MaxAge,
				Compress:   settings.Compress,
				LocalTime:  true,
			}
			writers = append(writers, fl)
			return io.MultiWriter(writers...), fl, filePath, nil
		}
	}
	return io.MultiWriter(writers...), nil, filePath, nil
}

func newSyslogWriter(settings LogSettings) (*syslogWriter, error) {
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
		return nil, err
	}
	return w, nil
}
