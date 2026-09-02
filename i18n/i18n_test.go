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

package i18n

// Author: deepseek-v4-pro / opencode

import "testing"

func TestLang(t *testing.T) {
	cases := map[string]string{
		"":                             "zh",
		"zh-CN,zh;q=0.9":               "zh",
		"zh":                           "zh",
		"en-US,en;q=0.9,zh-CN;q=0.8":   "en",
		"en":                           "en",
		"fr-FR":                        "zh",
		"en-US,fr-FR;q=0.9":            "en",
		"  zh-CN , zh ;q=0.9 ":         "zh",
		"fr-FR, en-US;q=0.9, zh;q=0.8": "zh",
	}
	for in, want := range cases {
		if got := Lang(in); got != want {
			t.Errorf("Lang(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTranslateKeepsChineseByDefault(t *testing.T) {
	msg := "用户名或密码错误"
	if got := Translate("zh", msg); got != msg {
		t.Errorf("zh translate should be identity, got %q", got)
	}
	if got := Translate("", msg); got != msg {
		t.Errorf("empty lang should be identity, got %q", got)
	}
}

func TestTranslateEnglishDoesNotTouchPureEnglish(t *testing.T) {
	for _, msg := range []string{"success", "ok", "database unavailable", "server x: boom"} {
		if got := Translate("en", msg); got != msg {
			t.Errorf("pure english %q should stay, got %q", msg, got)
		}
	}
}

func TestTranslateExactSentence(t *testing.T) {
	cases := []struct{ zh, en string }{
		{"参数错误", "Invalid parameters"},
		{"登录成功", "Login successful"},
		{"该用户已被禁用，请联系管理员", "This account has been disabled, please contact the administrator"},
		{"仅管理员可操作", "Admin privileges required"},
	}
	for _, c := range cases {
		if got := Translate("en", c.zh); got != c.en {
			t.Errorf("Translate(%q) = %q, want %q", c.zh, got, c.en)
		}
	}
}

func TestTranslateTemplateFragments(t *testing.T) {
	cases := []struct{ zh, en string }{
		{"CAS 认证失败: 响应超时", "CAS authentication failed: 响应超时"},
		{"my-server: 创建上游失败: dial tcp: refused", "my-server: Failed to create upstream: dial tcp: refused"},
		{"部分发布失败: a: 创建路由失败: err; b: 创建上游失败: err2", "Some servers failed to publish: a: Failed to create route: err; b: Failed to create upstream: err2"},
		{"AccessKey 不存在", "AccessKey not found"},
	}
	for _, c := range cases {
		if got := Translate("en", c.zh); got != c.en {
			t.Errorf("Translate(%q) = %q, want %q", c.zh, got, c.en)
		}
	}
}

func TestTranslateWholeMessagePatterns(t *testing.T) {
	cases := []struct{ zh, en string }{
		{"未找到用户 'zhangsan'", "No user found for 'zhangsan'"},
		{"搜索到 3 个匹配结果，请精确用户过滤器", "Found 3 matching entries; please make the user filter more precise"},
	}
	for _, c := range cases {
		if got := Translate("en", c.zh); got != c.en {
			t.Errorf("Translate(%q) = %q, want %q", c.zh, got, c.en)
		}
	}
}

func TestTranslateLongestMatchWins(t *testing.T) {
	// "MCP 服务器返回错误状态" is longer and must win over "MCP 服务器返回错误".
	msg := "MCP 服务器返回错误状态 500"
	want := "MCP server returned error status 500"
	if got := Translate("en", msg); got != want {
		t.Errorf("Translate(%q) = %q, want %q", msg, got, want)
	}
}

func TestTranslateUnknownPhraseKept(t *testing.T) {
	msg := "自定服务器出现未知错误"
	if got := Translate("en", msg); got != msg {
		t.Errorf("unknown chinese phrase should be kept verbatim, got %q", got)
	}
}
