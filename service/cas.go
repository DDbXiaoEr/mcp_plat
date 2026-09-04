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
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CasConfig struct {
	ServerUrl   string            `json:"serverUrl"`
	ServiceUrl  string            `json:"serviceUrl"`
	Version     string            `json:"version"`
	AttrMapping map[string]string `json:"attrMapping"`
}

type casServiceResponse struct {
	XMLName     xml.Name        `xml:"serviceResponse"`
	AuthSuccess *casAuthSuccess `xml:"authenticationSuccess"`
	AuthFailure *casAuthFailure `xml:"authenticationFailure"`
}

type casAuthSuccess struct {
	User string `xml:"user"`
}

type casAuthFailure struct {
	Code        string `xml:"code,attr"`
	Description string `xml:",innerxml"`
}

func casValidateTicket(serverUrl, serviceUrl, ticket string) (string, map[string]string, error) {
	validateUrl := fmt.Sprintf("%s/serviceValidate", strings.TrimRight(serverUrl, "/"))
	params := url.Values{}
	params.Set("service", serviceUrl)
	params.Set("ticket", ticket)

	fullUrl := fmt.Sprintf("%s?%s", validateUrl, params.Encode())

	resp, err := http.Get(fullUrl)
	if err != nil {
		return "", nil, fmt.Errorf("CAS 验证请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("CAS 读取响应失败: %v", err)
	}

	var casResp casServiceResponse
	if err := xml.Unmarshal(body, &casResp); err != nil {
		return "", nil, fmt.Errorf("CAS 响应解析失败: %v", err)
	}

	if casResp.AuthSuccess != nil {
		return casResp.AuthSuccess.User, parseCasAttributes(body), nil
	}

	if casResp.AuthFailure != nil {
		return "", nil, fmt.Errorf("CAS 认证失败: %s", strings.TrimSpace(casResp.AuthFailure.Description))
	}

	return "", nil, fmt.Errorf("CAS 认证失败")
}

var casControlElements = map[string]bool{
	"user":                   true,
	"proxyGrantingTicket":    true,
	"proxyGrantingTicketIou": true,
}

func parseCasAttributes(body []byte) map[string]string {
	attrs := make(map[string]string)
	dec := xml.NewDecoder(bytes.NewReader(body))

	var inSuccess bool
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if !inSuccess && t.Name.Local == "authenticationSuccess" {
				inSuccess = true
				continue
			}
			if !inSuccess {
				continue
			}

			local := t.Name.Local
			switch {
			case casControlElements[local]:
				skipElement(dec)
			case local == "attributes":
			default:
				if text, has, nested := readElementText(dec); has && !nested {
					if _, exists := attrs[local]; !exists {
						attrs[local] = text
					}
				}
			}
		case xml.EndElement:
			if !inSuccess {
				continue
			}
			if t.Name.Local == "authenticationSuccess" {
				inSuccess = false
				if len(attrs) > 0 {
					return attrs
				}
			}
		}
	}
	return attrs
}

func readElementText(dec *xml.Decoder) (text string, has bool, nested bool) {
	var sb strings.Builder
	for {
		tok, err := dec.Token()
		if err != nil {
			return "", false, false
		}

		switch t := tok.(type) {
		case xml.CharData:
			sb.Write(t)
		case xml.StartElement:
			skipElement(dec)
			nested = true
		case xml.EndElement:
			text = strings.TrimSpace(sb.String())
			return text, text != "", nested
		}
	}
}

func skipElement(dec *xml.Decoder) {
	depth := 1
	for depth > 0 {
		tok, err := dec.Token()
		if err != nil {
			return
		}
		switch tok.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}
	}
}
