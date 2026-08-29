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
	"log"

	"mcp_plat-console/model"
)

// 邮件通知统一封装。后续确定需要发邮件的场景（如登录异常、AccessKey 到期提醒、
// 新用户创建、管理员操作告警等）时，直接调用本文件的函数即可，例如：
//
//	if err := service.SendNotificationMail(user.Email, "主题", "<p>正文</p>"); err != nil {
//	    log.Printf("notify failed: %v", err)
//	}
//
// 若不想阻塞当前流程（如请求处理中），可用 `go service.SendNotificationMail(...)` 异步发送。

// SendNotificationMail 发送一封通知邮件。
// SMTP 配置从系统设置读取：未配置或未启用时返回明确错误，不发送。
func SendNotificationMail(to, subject, body string) error {
	cfg, err := loadSmtpConfig()
	if err != nil {
		return err
	}
	if err := SendMail(cfg, to, subject, body); err != nil {
		log.Printf("notify: send mail to %s failed: %v", to, err)
		return err
	}
	log.Printf("notify: mail sent to %s subject=%q", to, subject)
	return nil
}

// SendNotificationMailToMany 给多个收件人批量发送通知。
// SMTP 未配置/未启用时返回错误；单个收件人发送失败会记录日志并继续，
// 全部成功返回 nil，否则返回首个失败错误。
func SendNotificationMailToMany(tos []string, subject, body string) error {
	cfg, err := loadSmtpConfig()
	if err != nil {
		return err
	}
	var firstErr error
	for _, to := range tos {
		if to == "" {
			continue
		}
		if err := SendMail(cfg, to, subject, body); err != nil {
			log.Printf("notify: send mail to %s failed: %v", to, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		log.Printf("notify: mail sent to %s subject=%q", to, subject)
	}
	return firstErr
}

// SendNotificationMailToUser 给指定用户发送通知邮件。
// 用户未填写邮箱时跳过发送并返回 nil（不视为错误）。
func SendNotificationMailToUser(user model.User, subject, body string) error {
	if user.Email == "" {
		return nil
	}
	return SendNotificationMail(user.Email, subject, body)
}

// loadSmtpConfig 读取系统设置中保存的 SMTP 配置，并校验是否已启用。
func loadSmtpConfig() (SmtpSetting, error) {
	var cfg SmtpSetting
	if err := GetSetting("smtp", &cfg); err != nil {
		return cfg, errors.New("邮件通知未配置，请在系统设置中配置 SMTP")
	}
	if !cfg.Enabled {
		return cfg, errors.New("邮件通知未启用，请在系统设置中开启邮件通知开关")
	}
	return cfg, nil
}
