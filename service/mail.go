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

package service

// Author: deepseek-v4-pro / opencode

import (
	"errors"
	"fmt"
	"time"

	mail "github.com/wneessen/go-mail"
)

type SmtpSetting struct {
	Enabled     bool   `json:"enabled"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Encryption  string `json:"encryption"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	FromAddress string `json:"fromAddress"`
	FromName    string `json:"fromName"`
}

func SendMail(cfg SmtpSetting, to, subject, body string) error {
	if to == "" {
		return errors.New("收件人邮箱不能为空")
	}
	if cfg.Host == "" {
		return errors.New("SMTP 服务器地址未配置")
	}
	if cfg.Port == 0 {
		cfg.Port = 465
	}

	from := cfg.FromAddress
	if from == "" {
		from = cfg.Username
	}
	if from == "" {
		return errors.New("缺少发件人地址")
	}

	opts := []mail.Option{
		mail.WithPort(cfg.Port),
		mail.WithTimeout(15 * time.Second),
	}
	if cfg.Username != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(cfg.Username),
			mail.WithPassword(cfg.Password),
		)
	}
	switch cfg.Encryption {
	case "ssl":
		opts = append(opts, mail.WithSSL())
	case "starttls":
		opts = append(opts, mail.WithTLSPolicy(mail.TLSMandatory))
	default:
		opts = append(opts, mail.WithTLSPolicy(mail.NoTLS))
	}

	m := mail.NewMsg()
	if err := m.From(fmt.Sprintf("%s <%s>", fromNameOr(cfg.FromName, cfg.FromAddress), from)); err != nil {
		return fmt.Errorf("发件人地址不合法: %w", err)
	}
	if err := m.To(to); err != nil {
		return fmt.Errorf("收件人地址不合法: %w", err)
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body)

	c, err := mail.NewClient(cfg.Host, opts...)
	if err != nil {
		return fmt.Errorf("初始化 SMTP 客户端失败: %w", err)
	}
	if err := c.DialAndSend(m); err != nil {
		return fmt.Errorf("邮件发送失败: %w", err)
	}
	return nil
}

func fromNameOr(name, fallback string) string {
	if name != "" {
		return name
	}
	return fallback
}

func TestSmtp(cfg SmtpSetting, to string) error {
	name := GetPlatform().Name
	return SendMail(cfg, to,
		fmt.Sprintf("%s测试邮件", name),
		fmt.Sprintf("<p>这是一封来自「%s」的测试邮件，收到此邮件说明 SMTP 配置正常。</p>", name))
}
