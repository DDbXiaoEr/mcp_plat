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

// zhEnMessages is the Chinese -> English phrase dictionary. Keys may be whole
// sentences or fragments of fmt.Errorf templates; at translation time the
// longest registered phrase is replaced first, which keeps interpolated values
// (server names, error details, counts) intact.
var zhEnMessages = map[string]string{
	// handlers / middleware
	"参数错误":   "Invalid parameters",
	"查询失败":   "Query failed",
	"创建失败":   "Failed to create",
	"创建成功":   "Created successfully",
	"更新成功":   "Updated successfully",
	"删除成功":   "Deleted successfully",
	"获取成功":   "Fetched successfully",
	"发布成功":   "Published successfully",
	"维护设置成功": "Maintenance settings updated",
	"分配成功":   "Users assigned successfully",
	"保存成功":   "Saved successfully",
	"登录成功":   "Login successful",
	"登录请求过于频繁，请稍后再试": "Too many login attempts, please try again later",
	"管理员账号不支持修改邮箱":   "Email cannot be changed for admin accounts",
	"收件人邮箱不能为空":      "Recipient email cannot be empty",
	"测试邮件发送成功":       "Test email sent successfully",
	"未提供认证令牌":        "Missing authentication token",
	"认证格式错误":         "Invalid authentication format",
	"令牌无效或已过期":       "Invalid or expired token",
	"仅管理员可操作":        "Admin privileges required",

	// auth / cas
	"用户名或密码错误":       "Incorrect username or password",
	"该用户已被禁用，请联系管理员": "This account has been disabled, please contact the administrator",
	"密码加密失败":         "Failed to encrypt password",
	"创建用户失败":         "Failed to create user",
	"CAS 配置未找到":      "CAS configuration not found",
	"CAS 认证未启用":      "CAS authentication is not enabled",
	"CAS 服务地址未配置":    "CAS service URL is not configured",
	"CAS 验证请求失败":     "CAS validation request failed",
	"CAS 读取响应失败":     "Failed to read CAS response",
	"CAS 响应解析失败":     "Failed to parse CAS response",
	"CAS 认证失败":       "CAS authentication failed",
	"用户不存在":          "User not found",
	"批量创建成功":         "Users created successfully",
	"用户列表不能为空":       "The user list is empty",
	"用户名和密码不能为空":     "Username and password are required",
	"生成学号/工号失败，请重试":  "Failed to generate a student/staff ID, please try again",
	"邮箱不能为空":         "Email cannot be empty",
	"邮箱格式不正确":        "Invalid email format",
	"邮箱长度超出限制":       "Email exceeds the maximum length",
	"该邮箱已被其他账号使用":    "This email is already in use by another account",
	"保存失败，请稍后重试":     "Save failed, please try again later",

	// settings
	"不支持的设置项":  "Unsupported setting",
	"设置项不存在":   "Setting not found",
	"设置内容格式错误": "Invalid setting content format",

	// access key
	"AccessKey 不存在": "AccessKey not found",
	"该 AccessKey 已过期，无法启用，请删除后重新创建": "This AccessKey has expired and cannot be enabled, please delete it and create a new one",

	// roles / rbac users
	"角色不存在": "Role not found",
	"该角色下存在用户，请先解除用户绑定": "Users are still assigned to this role, please unassign them first",
	"用户名已存在": "Username already exists",

	// mail / smtp
	"SMTP 服务器地址未配置":  "SMTP server address is not configured",
	"缺少发件人地址":        "Sender address is missing",
	"发件人地址不合法":       "Invalid sender address",
	"收件人地址不合法":       "Invalid recipient address",
	"初始化 SMTP 客户端失败": "Failed to initialize SMTP client",
	"邮件发送失败":         "Failed to send email",

	// mcp server fetch tools
	"MCP 服务器不存在":            "MCP server not found",
	"请先在系统设置中配置 API 网关默认域名": "Please configure the API gateway default domain in system settings first",
	"stdio 协议暂不支持远程获取工具列表":  "stdio protocol does not support fetching remote tool lists yet",
	"不支持的协议类型":              "Unsupported protocol type",
	"URL 解析失败":              "Failed to parse URL",
	"仅支持 http/https 协议":     "Only http/https protocols are supported",
	"DNS 解析失败":              "DNS resolution failed",
	"不允许连接到内网地址":            "Connecting to internal addresses is not allowed",
	"拒绝连接到内网地址":             "Refused to connect to internal address",
	"重定向次数过多":               "Too many redirects",
	"重定向目标 DNS 解析失败":        "DNS resolution of the redirect target failed",
	"拒绝重定向到内网地址":            "Refused to redirect to an internal address",
	"初始化 MCP 会话失败":          "Failed to initialize MCP session",
	"MCP 初始化返回错误":           "MCP initialization returned an error",
	"MCP 服务器返回错误":           "MCP server returned an error",
	"MCP 服务器返回错误状态":         "MCP server returned error status",
	"解析工具列表失败":              "Failed to parse tool list",
	"构建请求失败":                "Failed to build request",
	"连接 MCP 服务器失败":          "Failed to connect to MCP server",
	"解析响应失败":                "Failed to parse response",
	"读取 SSE 响应失败":           "Failed to read SSE response",
	"未从 SSE 响应中读取到数据":       "No data read from SSE response",
	"连接 SSE 服务器失败":          "Failed to connect to SSE server",
	"SSE 端点不允许使用绝对 URL":     "SSE endpoint does not allow absolute URLs",
	"解析 SSE 端点地址失败":         "Failed to parse SSE endpoint URL",
	"解析 SSE 端点相对路径失败":       "Failed to parse SSE endpoint relative path",
	"读取 SSE 流失败":            "Failed to read SSE stream",
	"未能获取 SSE 消息端点":         "Failed to obtain the SSE message endpoint",
	"通过 SSE 端点获取工具列表失败":     "Failed to fetch tool list via SSE endpoint",

	// publish / maintenance / gateway
	"请先在系统设置中配置 API 网关参数":         "Please configure the API gateway parameters in system settings first",
	"API 网关 Admin API 地址未配置":      "API gateway Admin API address is not configured",
	"API 网关 Admin API Key 未配置":    "API gateway Admin API key is not configured",
	"查询服务器失败":                     "Failed to query servers",
	"未找到指定的服务器":                   "Specified server not found",
	"未配置后端服务地址":                   "Backend service address not configured",
	"创建上游失败":                      "Failed to create upstream",
	"设置上游节点失败":                    "Failed to set upstream nodes",
	"创建 Service 失败":               "Failed to create Service",
	"创建路由失败":                      "Failed to create route",
	"启用认证但未配置 gRPC 校验地址":          "Authentication is enabled but no gRPC verification address is configured",
	"启用审计日志但未配置 gRPC 地址，跳过审计日志插件": "Audit logging is enabled but no gRPC address is configured, skipping the audit log plugin",
	"保存路由配置失败":                    "Failed to save route configuration",
	"部分发布失败":                      "Some servers failed to publish",
	"更新发布状态失败":                    "Failed to update publish status",
	"设置维护失败":                      "Failed to enable maintenance",
	"更新维护状态失败":                    "Failed to update maintenance status",
	"取消维护失败":                      "Failed to disable maintenance",
	"更新状态失败":                      "Failed to update status",
	"操作部分失败":                      "Some operations failed",
	"请求 Kong Admin API 失败":        "Request to Kong Admin API failed",
	"请求 APISIX Admin API 失败":      "Request to APISIX Admin API failed",
	"Kong 返回错误":                   "Kong returned an error",
	"APISIX 返回错误":                 "APISIX returned an error",

	// ldap
	"LDAP 服务未配置":   "LDAP service is not configured",
	"LDAP 连接失败":    "LDAP connection failed",
	"LDAP 管理员绑定失败": "LDAP admin bind failed",
	"LDAP 搜索失败":    "LDAP search failed",
	"账号异常":         "Unexpected account state",
	"未找到用户":        "No user found for",

	// overview chart fallback labels
	"未知服务器": "Unknown server",
	"其他":    "Other",
	"未分组":   "Ungrouped",
}

// zhEnPatterns holds whole-message templates whose interpolated values cannot be
// translated by simple phrase substitution. Each regexp matches a fully
// formatted message and rewrites it with $1..$n captures.
var zhEnPatterns = []struct {
	re   string
	tmpl string
}{
	{`^未找到用户 '(.*)'$`, "No user found for '$1'"},
	{`^搜索到 (\d+) 个匹配结果，请精确用户过滤器$`, "Found $1 matching entries; please make the user filter more precise"},
	{`^列表中存在重复的用户名 '(.*)'$`, "Duplicate username in the list: '$1'"},
	{`^用户名 '(.*)' 已存在$`, "Username '$1' already exists"},
	{`^学号/工号 '(.*)' 已被使用$`, "Student/Staff ID '$1' is already in use"},
	{`^邮箱 '(.*)' 已被其他账号使用$`, "Email '$1' is already used by another account"},
}
