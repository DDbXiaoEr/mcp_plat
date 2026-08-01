package service

// Author: deepseek-v4-pro / opencode

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

type LdapConfig struct {
	Host         string            `json:"host"`
	Port         int               `json:"port"`
	BaseDn       string            `json:"baseDn"`
	BindDn       string            `json:"bindDn"`
	BindPassword string            `json:"bindPassword"`
	UserFilter   string            `json:"userFilter"`
	AttrMapping  map[string]string `json:"attrMapping"`
}

func ldapAuthenticate(cfg LdapConfig, username, password string) (map[string]string, error) {
	if cfg.Host == "" || cfg.Port == 0 {
		return nil, fmt.Errorf("LDAP 服务未配置")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", addr))
	if err != nil {
		return nil, fmt.Errorf("LDAP 连接失败: %v", err)
	}
	defer conn.Close()

	if cfg.BindDn != "" && cfg.BindPassword != "" {
		if err := conn.Bind(cfg.BindDn, cfg.BindPassword); err != nil {
			return nil, fmt.Errorf("LDAP 管理员绑定失败: %v", err)
		}
	}

	filter := fmt.Sprintf("(uid=%s)", username)
	if cfg.UserFilter != "" {
		filter = strings.ReplaceAll(cfg.UserFilter, "%s", username)
	}

	searchAttrs := []string{"dn"}
	for _, ldapAttr := range cfg.AttrMapping {
		searchAttrs = append(searchAttrs, ldapAttr)
	}

	searchReq := ldap.NewSearchRequest(
		cfg.BaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		searchAttrs,
		nil,
	)

	result, err := conn.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("LDAP 搜索失败: %v", err)
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	if len(result.Entries) > 1 {
		return nil, fmt.Errorf("账号异常")
	}

	userDn := result.Entries[0].DN

	attrs := make(map[string]string)
	entry := result.Entries[0]
	for platformField, ldapAttr := range cfg.AttrMapping {
		for _, attr := range entry.Attributes {
			if strings.EqualFold(attr.Name, ldapAttr) && len(attr.Values) > 0 {
				attrs[platformField] = attr.Values[0]
				break
			}
		}
	}

	if cfg.BindDn == "" || cfg.BindPassword == "" {
		newConn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", addr))
		if err != nil {
			return nil, fmt.Errorf("LDAP 连接失败: %v", err)
		}
		defer newConn.Close()

		if err := newConn.Bind(userDn, password); err != nil {
			return nil, fmt.Errorf("用户名或密码错误")
		}
		return attrs, nil
	}

	if err := conn.Bind(userDn, password); err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	return attrs, nil
}

type TestLdapInput struct {
	LdapConfig LdapConfig `json:"ldap"`
	Username   string     `json:"username"`
}

type LdapAttributeInfo struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type TestLdapOutput struct {
	Success    bool                `json:"success"`
	Message    string              `json:"message"`
	UserDn     string              `json:"userDn,omitempty"`
	Attributes []LdapAttributeInfo `json:"attributes,omitempty"`
	Mapped     map[string]string   `json:"mapped,omitempty"`
}

func TestLdapMapping(input TestLdapInput) *TestLdapOutput {
	cfg := input.LdapConfig

	if cfg.Host == "" || cfg.Port == 0 {
		return &TestLdapOutput{Success: false, Message: "LDAP 服务未配置"}
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", addr))
	if err != nil {
		return &TestLdapOutput{Success: false, Message: fmt.Sprintf("LDAP 连接失败: %v", err)}
	}
	defer conn.Close()

	if cfg.BindDn != "" && cfg.BindPassword != "" {
		if err := conn.Bind(cfg.BindDn, cfg.BindPassword); err != nil {
			return &TestLdapOutput{Success: false, Message: fmt.Sprintf("LDAP 管理员绑定失败: %v", err)}
		}
	}

	filter := fmt.Sprintf("(uid=%s)", input.Username)
	if cfg.UserFilter != "" {
		filter = strings.ReplaceAll(cfg.UserFilter, "%s", input.Username)
	}

	searchReq := ldap.NewSearchRequest(
		cfg.BaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"*"},
		nil,
	)

	result, err := conn.Search(searchReq)
	if err != nil {
		return &TestLdapOutput{Success: false, Message: fmt.Sprintf("LDAP 搜索失败: %v", err)}
	}

	if len(result.Entries) == 0 {
		return &TestLdapOutput{Success: false, Message: fmt.Sprintf("未找到用户 '%s'", input.Username)}
	}

	if len(result.Entries) > 1 {
		return &TestLdapOutput{Success: false, Message: fmt.Sprintf("搜索到 %d 个匹配结果，请精确用户过滤器", len(result.Entries))}
	}

	entry := result.Entries[0]

	attrs := make([]LdapAttributeInfo, 0, len(entry.Attributes))
	for _, attr := range entry.Attributes {
		attrs = append(attrs, LdapAttributeInfo{
			Name:   attr.Name,
			Values: attr.Values,
		})
	}

	mapped := make(map[string]string)
	for platformField, ldapAttr := range cfg.AttrMapping {
		for _, attr := range entry.Attributes {
			if strings.EqualFold(attr.Name, ldapAttr) && len(attr.Values) > 0 {
				mapped[platformField] = attr.Values[0]
				break
			}
		}
	}

	return &TestLdapOutput{
		Success:    true,
		Message:    "查询成功",
		UserDn:     entry.DN,
		Attributes: attrs,
		Mapped:     mapped,
	}
}
