package service

// Author: deepseek-v4-pro / opencode

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CasConfig struct {
	ServerUrl  string `json:"serverUrl"`
	ServiceUrl string `json:"serviceUrl"`
	Version    string `json:"version"`
}

type casServiceResponse struct {
	XMLName      xml.Name           `xml:"serviceResponse"`
	AuthSuccess  *casAuthSuccess    `xml:"authenticationSuccess"`
	AuthFailure  *casAuthFailure    `xml:"authenticationFailure"`
}

type casAuthSuccess struct {
	User string `xml:"user"`
}

type casAuthFailure struct {
	Code        string `xml:"code,attr"`
	Description string `xml:",innerxml"`
}

func casValidateTicket(serverUrl, serviceUrl, ticket string) (string, error) {
	validateUrl := fmt.Sprintf("%s/serviceValidate", strings.TrimRight(serverUrl, "/"))
	params := url.Values{}
	params.Set("service", serviceUrl)
	params.Set("ticket", ticket)

	fullUrl := fmt.Sprintf("%s?%s", validateUrl, params.Encode())

	resp, err := http.Get(fullUrl)
	if err != nil {
		return "", fmt.Errorf("CAS 验证请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("CAS 读取响应失败: %v", err)
	}

	var casResp casServiceResponse
	if err := xml.Unmarshal(body, &casResp); err != nil {
		return "", fmt.Errorf("CAS 响应解析失败: %v", err)
	}

	if casResp.AuthSuccess != nil {
		return casResp.AuthSuccess.User, nil
	}

	if casResp.AuthFailure != nil {
		return "", fmt.Errorf("CAS 认证失败: %s", strings.TrimSpace(casResp.AuthFailure.Description))
	}

	return "", fmt.Errorf("CAS 认证失败")
}
