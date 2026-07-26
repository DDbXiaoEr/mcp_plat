package service

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

type LdapConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	BaseDn       string `json:"baseDn"`
	BindDn       string `json:"bindDn"`
	BindPassword string `json:"bindPassword"`
	UserFilter   string `json:"userFilter"`
}

func ldapAuthenticate(cfg LdapConfig, username, password string) error {
	if cfg.Host == "" || cfg.Port == 0 {
		return fmt.Errorf("LDAP 服务未配置")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", addr))
	if err != nil {
		return fmt.Errorf("LDAP 连接失败: %v", err)
	}
	defer conn.Close()

	if cfg.BindDn != "" && cfg.BindPassword != "" {
		if err := conn.Bind(cfg.BindDn, cfg.BindPassword); err != nil {
			return fmt.Errorf("LDAP 管理员绑定失败: %v", err)
		}
	}

	filter := fmt.Sprintf("(uid=%s)", username)
	if cfg.UserFilter != "" {
		filter = strings.ReplaceAll(cfg.UserFilter, "%s", username)
	}

	searchReq := ldap.NewSearchRequest(
		cfg.BaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"dn"},
		nil,
	)

	result, err := conn.Search(searchReq)
	if err != nil {
		return fmt.Errorf("LDAP 搜索失败: %v", err)
	}

	if len(result.Entries) == 0 {
		return fmt.Errorf("用户名或密码错误")
	}

	if len(result.Entries) > 1 {
		return fmt.Errorf("账号异常")
	}

	userDn := result.Entries[0].DN

	if cfg.BindDn == "" || cfg.BindPassword == "" {
		newConn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", addr))
		if err != nil {
			return fmt.Errorf("LDAP 连接失败: %v", err)
		}
		defer newConn.Close()

		if err := newConn.Bind(userDn, password); err != nil {
			return fmt.Errorf("用户名或密码错误")
		}
		return nil
	}

	if err := conn.Bind(userDn, password); err != nil {
		return fmt.Errorf("用户名或密码错误")
	}

	return nil
}
