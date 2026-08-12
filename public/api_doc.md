# API 文档

Base URL: `http://localhost:8080`

通用响应格式：
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

---

## 1. 认证模块

### 1.1 登录

```
POST /api/auth/login
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**成功响应：**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOi...",
    "username": "admin",
    "role": "admin"
  }
}
```

### 1.2 获取个人信息

```
GET /api/auth/profile
```

**请求头：**
| 字段 | 说明 |
|------|------|
| Authorization | Bearer {token} |

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "uid": "2021001",
    "username": "user1",
    "email": "user1@example.com",
    "phone": "13800000000",
    "organization": "计算机学院",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

### 1.3 获取认证方式

```
GET /api/auth/method
```

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "method": "cas",
    "cas": { "serverUrl": "https://cas.example.edu.cn", "serviceUrl": "http://localhost:5174", "version": "3.0" }
  }
}
```

`method` 可能的值：`local`（本地数据库）、`ldap`、`cas`。

### 1.4 CAS 登录验证

```
POST /api/auth/cas/validate
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ticket | string | 是 | CAS 返回的 ticket |
| serviceUrl | string | 是 | 回调地址，需与 CAS redirect 时一致 |

**成功响应：**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOi...",
    "username": "2021001",
    "role": "user"
  }
}
```

---

## 2. AccessKey 管理

以下接口均需携带认证头：`Authorization: Bearer {token}`

### 2.1 获取 AccessKey 列表

```
GET /api/access-keys
```

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "name": "my-key",
      "key": "ak-abc123...",
      "enabled": true,
      "expired_at": "2027-01-01T00:00:00Z",
      "servers": "[{\"name\":\"server1\",\"tools\":[\"tool1\",\"tool2\"]}]",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

### 2.2 创建 AccessKey

```
POST /api/access-keys
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | AccessKey 名称 |
| servers | string | 否 | 服务器权限配置（JSON 字符串） |
| expired_at | string | 否 | 过期时间（ISO 8601 格式） |

**成功响应：**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 1,
    "name": "my-key",
    "key": "ak-abc123...",
    "enabled": true,
    "expired_at": null,
    "servers": ""
  }
}
```

### 2.3 更新 AccessKey

```
PUT /api/access-keys/:id
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | AccessKey 名称 |
| enabled | bool | 否 | 是否启用 |
| servers | string | 否 | 服务器权限配置（JSON 字符串） |
| expired_at | string | 否 | 过期时间（ISO 8601 格式） |

**成功响应：**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

### 2.4 删除 AccessKey

```
DELETE /api/access-keys/:id
```

**成功响应：**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

## 3. 使用历史

> ✅ 已实现：基于审计日志（audit_logs）。存储后端可配置（关系库 Postgres/MySQL 或 ClickHouse），见 `config.yaml` 的 `audit_log_db`。

### 3.1 获取使用历史列表

```
GET /api/history
```

**请求参数（Query String）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| access_key | string | 否 | 按 AccessKey 筛选（传完整 key，服务端哈希后匹配） |
| server_id | string | 否 | 按 MCP 服务器筛选 |
| tool_name | string | 否 | 按工具名筛选 |
| start | string | 否 | 开始日期（YYYY-MM-DD） |
| end | string | 否 | 结束日期（YYYY-MM-DD） |
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20（最大 100） |

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "access_key_id": 3,
        "access_key_name": "我的测试 Key",
        "user_id": 1,
        "server_id": "b8a1...",
        "tool_name": "get_weather",
        "success": true,
        "message": "ok",
        "created_at": "2026-01-01T12:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

> 说明：`access_key_id` 为 AccessKey 在数据库中的 ID（不落明文 key，减少数据量）；`access_key_name` 为服务端根据当前用户持有的 Key 反查的名称，Key 已删除时可能为空。

---

## 4. MCP 服务器管理

> 所有接口均需认证，在请求头中携带 `Authorization: Bearer <token>`

### 4.1 获取 MCP 服务器列表

```
GET /api/servers
```

**请求参数：** 无

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "教务系统 MCP",
      "address": "https://mcp.xauat.edu.cn/jwc",
      "service_address": "[\"192.168.1.100:8081\",\"192.168.1.101:8081\"]",
      "department": "教务处",
      "protocol": "SSE",
      "protocol_version": "2026-07-28",
      "tools": "[{\"name\":\"查询课表\",\"description\":\"查询学期课程安排\"},{\"name\":\"成绩查询\",\"description\":\"查询考试成绩\"}]",
      "description": "教务系统 MCP 服务，提供课表查询与成绩查询能力",
      "status": "published",
      "created_at": "2026-01-01T12:00:00Z",
      "updated_at": "2026-01-01T12:00:00Z"
    }
  ]
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 服务器 ID |
| name | string | 服务器名称 |
| address | string | MCP 服务器地址（API 网关路径或完整 URL） |
| service_address | string | MCP 服务地址，JSON 数组字符串（如 `["192.168.1.100:8081","192.168.1.101:8081"]`），支持多个地址，后端发布时作为 API 网关 upstream 节点 |
| department | string | 负责部门 |
| protocol | string | 协议类型（SSE / Streamable HTTP / stdio） |
| protocol_version | string | MCP 协议版本（2025-03-26 / 2025-06-18 / 2026-07-28），默认 2026-07-28 |
| tools | string | 工具列表，JSON 字符串数组格式，每个元素为 `{"name":"...","description":"..."}` |
| description | string | 服务器描述信息 |
| status | string | 发布状态，`published`（已发布）或 `unpublished`（未发布） |

### 4.2 新增 MCP 服务器

```
POST /api/servers
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 服务器名称 |
| address | string | 是 | MCP 服务器地址（API 网关路径或完整 URL） |
| service_address | string | 否 | MCP 服务地址，JSON 数组字符串（如 `["192.168.1.100:8081","192.168.1.101:8081"]`），支持多个地址 |
| department | string | 否 | 负责部门 |
| protocol | string | 否 | 协议类型，默认 SSE |
| protocol_version | string | 否 | MCP 协议版本，默认 2026-07-28 |
| tools | string | 否 | 工具列表，JSON 字符串数组格式 |
| description | string | 否 | 服务器描述信息 |

**响应示例：**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 1,
    "name": "教务系统 MCP",
    "address": "https://mcp.xauat.edu.cn/jwc",
    "service_address": "192.168.1.100:8081",
    "department": "教务处",
    "protocol": "SSE",
    "protocol_version": "2026-07-28",
    "tools": "[\"查询课表\",\"成绩查询\"]",
    "created_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-01-01T12:00:00Z"
  }
}
```

### 4.5 发布服务器到 API 网关

```
POST /api/servers/publish
```

> 需管理员权限。

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| server_ids | []string | 是 | 要发布的 MCP 服务器 ID 列表 |
| enable_auth | bool | 否 | 是否启用 Access Key 认证（默认 false），启用后路由会下发 accesskey_verify 插件 |
| accesskey_header | string | 否 | 自定义 Access Key 的 HTTP Header 名称，未传时优先使用网关设置中配置的值，均未配置时默认 `X-Access-Key` |

**请求示例：**
```json
{
  "server_ids": ["uuid-1", "uuid-2"],
  "enable_auth": true
}
```

**成功响应：**
```json
{
  "code": 200,
  "message": "发布成功"
}
```

**说明：** 后端根据已配置的 API 网关设置（`api_gateway` 分组），将选中的 MCP 服务器路由注册到 API 网关中。需先配置网关的 Admin API 地址和 Key。若启用认证，还需配置 `authGrpcAddr` 指向 accesskey-auth-server 的 gRPC 地址。

### 4.6 获取 MCP 服务器工具列表

```
POST /api/servers/fetch-tools
```

> 所有接口均需认证，在请求头中携带 `Authorization: Bearer <token>`

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| address | string | 是 | MCP 服务器地址 |
| protocol | string | 否 | 协议类型（SSE / Streamable HTTP / stdio），默认 Streamable HTTP |
| protocol_version | string | 否 | MCP 协议版本，默认 2026-07-28 |

**响应示例：**
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "tools": [
      { "name": "查询课表", "description": "查询学期课程安排" },
      { "name": "成绩查询", "description": "查询考试成绩" },
      { "name": "选课信息", "description": "查询选课信息" }
    ]
  }
}
```

**说明：** 后端根据地址和协议连接目标 MCP 服务器，调用 `tools/list` 协议方法，返回可用工具名称及描述列表。连接失败或地址无效时返回 4xx/5xx。

### 4.3 更新 MCP 服务器

```
PUT /api/servers/:id
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 服务器名称 |
| address | string | 否 | MCP 服务器地址 |
| service_address | string | 否 | MCP 服务地址，JSON 数组字符串，支持多个地址 |
| department | string | 否 | 负责部门 |
| protocol | string | 否 | 协议类型 |
| protocol_version | string | 否 | MCP 协议版本 |
| tools | string | 否 | 工具列表，JSON 字符串数组格式 |
| description | string | 否 | 服务器描述信息 |

**响应示例：**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

### 4.4 删除 MCP 服务器

```
DELETE /api/servers/:id
```

**请求参数：** 无（ID 通过 URL 路径传递）

**响应示例：**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

## 5. RBAC 管理

> 所有接口均需认证，在请求头中携带 `Authorization: Bearer <token>`，且仅管理员可调用。

### 5.1 角色管理

#### 5.1.1 获取角色列表

```
GET /api/roles
```

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "超级管理员",
      "description": "拥有全部权限",
      "server_ids": [1, 2, 3, 4, 5, 6, 7, 8, 9],
      "user_count": 1,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "name": "MCP管理员",
      "description": "管理MCP服务器与工具",
      "server_ids": [1, 2, 7],
      "user_count": 1,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 角色 ID |
| name | string | 角色名称 |
| description | string | 角色描述 |
| server_ids | []uint | 可使用的 MCP 服务器 ID 列表 |
| user_count | int | 绑定该角色的用户数 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |

#### 5.1.2 创建角色

```
POST /api/roles
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 角色名称 |
| description | string | 否 | 角色描述 |
| server_ids | []uint | 否 | 可使用的 MCP 服务器 ID 列表 |

**成功响应：**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 4,
    "name": "审计员",
    "description": "查看所有日志",
    "server_ids": [1, 4, 6],
    "user_count": 0,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

#### 5.1.3 更新角色

```
PUT /api/roles/:id
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 角色名称 |
| description | string | 否 | 角色描述 |
| server_ids | []uint | 否 | 可使用的 MCP 服务器 ID 列表 |

**成功响应：**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

#### 5.1.4 删除角色

```
DELETE /api/roles/:id
```

> 注意：有关联用户的角色不可删除，需先解除用户绑定。

**成功响应：**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

### 5.2 用户管理

#### 5.2.1 获取用户列表

```
GET /api/users
```

**请求参数（Query String）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| role_id | uint | 否 | 按角色筛选 |

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "username": "admin",
      "role_id": 1,
      "role_name": "超级管理员",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "username": "zhangsan",
      "role_id": 2,
      "role_name": "MCP管理员",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    },
    {
      "id": 5,
      "username": "zhaoliu",
      "role_id": null,
      "role_name": "",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 用户 ID |
| username | string | 用户名 |
| role_id | uint/null | 绑定的角色 ID，null 表示未分配 |
| role_name | string | 绑定的角色名称，未分配时为空字符串 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |

#### 5.2.2 创建用户

```
POST /api/users
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 登录密码 |
| role_id | uint | 否 | 绑定的角色 ID |

**成功响应：**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 8,
    "username": "newuser",
    "role_id": null,
    "role_name": "",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

#### 5.2.3 更新用户

```
PUT /api/users/:id
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 否 | 用户名 |
| password | string | 否 | 新密码 |
| role_id | uint | 否 | 绑定的角色 ID（传 null 或 0 表示取消角色绑定） |

**成功响应：**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

#### 5.2.4 删除用户

```
DELETE /api/users/:id
```

**成功响应：**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

### 5.3 角色-用户批量分配

```
PUT /api/roles/:id/users
```

> 传入目标用户 ID 列表，将该角色全量分配给这些用户。不在列表中的用户将被解除与该角色的绑定。

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_ids | []uint | 是 | 分配到该角色的用户 ID 列表 |

**请求示例：**
```json
{
  "user_ids": [1, 2, 5, 8]
}
```

**成功响应：**
```json
{
  "code": 200,
  "message": "分配成功"
}
```

---

### 5.4 获取角色已分配用户

```
GET /api/roles/:id/users
```

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "username": "admin",
      "role_id": 1,
      "role_name": "超级管理员",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

---

## 6. 系统设置

> 以下接口均需管理员权限（`Authorization: Bearer <token>`，且用户为管理员）。

### 6.1 获取全部设置

```
GET /api/settings
```

**成功响应：**

`data` 为按设置项分组的对象，仅包含已保存过的分组；每组内容为保存时的原始 JSON。

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "log": {
      "syslogEnabled": false,
      "logPath": "/var/log/mcp-plat",
      "logLevel": "info",
      "logPrefix": "mcp-plat",
      "syslogHost": "",
      "syslogPort": 514,
      "syslogProtocol": "tcp"
    },
    "smtp": {
      "enabled": true,
      "host": "smtp.example.edu.cn",
      "port": 465,
      "encryption": "ssl",
      "username": "noreply",
      "password": "******",
      "fromAddress": "noreply@example.edu.cn",
      "fromName": "MCP 服务平台"
    },
    "auth": {
      "method": "cas",
      "cas": { "serverUrl": "", "serviceUrl": "", "version": "3.0" },
      "ldap": { "host": "", "port": 389, "baseDn": "", "bindDn": "", "bindPassword": "", "userFilter": "" },
      "oauth": { "authorizeUrl": "", "tokenUrl": "", "userinfoUrl": "", "clientId": "", "clientSecret": "", "redirectUrl": "", "scope": "" }
    },
    "user_ops": {
      "defaultRoleId": 1,
      "maxAccessKeys": 5
    },
    "platform": {
      "name": "某某大学",
      "logoUrl": "https://www.example.edu.cn/logo.png",
      "siteUrl": "https://www.example.edu.cn"
    },
    "api_gateway": {
      "provider": "apisix",
      "adminUrl": "http://127.0.0.1:9180",
      "adminKey": "edd1c9f034335f136f87ad84b625c8f1",
      "defaultPublishDomain": "mcp.xauat.edu.cn",
      "authGrpcAddr": ":9090",
      "accesskeyHeader": "X-Access-Key"
    },
    "network_security": {
      "allowlist": ["10.0.0.0/8", "172.16.0.0/12", "192.168.1.0/24"]
    }
  }
}
```

### 6.2 保存设置

```
PUT /api/settings/:key
```

**路径参数：**

| 参数 | 说明 |
|------|------|
| key | 设置项分组，支持 `log` / `smtp` / `auth` / `user_ops` / `platform` / `api_gateway` / `network_security` |

**请求参数（JSON Body）：**

任意合法 JSON 对象，整组覆盖保存（后端原样存储，不校验字段）。

> `log` 分组在服务启动时加载：`syslogEnabled` 为 true 且 `syslogHost` 非空时，日志输出重定向到 Syslog 服务器（`syslogProtocol` 支持 tcp/udp，默认端口 514，`logPrefix` 作为 tag）；未启用或连接失败时输出到标准输出。修改后需重启服务生效。

**成功响应：**
```json
{
  "code": 200,
  "message": "保存成功"
}
```

**失败响应：**
```json
{
  "code": 400,
  "message": "不支持的设置项"
}
```

---

### 6.3 获取单个设置项

```
GET /api/settings/:key
```

> 仅需登录，无需管理员权限。

**路径参数：**

| 参数 | 说明 |
|------|------|
| key | 设置项分组，支持 `log` / `smtp` / `auth` / `user_ops` / `platform` / `api_gateway` / `network_security` |

**响应示例（`GET /api/settings/platform`）：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "name": "某某大学",
    "logoUrl": "https://www.example.edu.cn/logo.png",
    "siteUrl": "https://www.example.edu.cn"
  }
}
```

---

### 6.4 检查 API 网关配置状态

```
GET /api/settings/gateway-status
```

> 无需管理员权限，仅需登录。

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "configured": true
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| configured | bool | API 网关 Admin API 地址和 Key 是否均已配置 |

---

## 7. 平台概况统计

> 仅管理员可调用，需携带 `Authorization: Bearer <token>`。

### 7.1 获取平台概况统计

```
GET /api/overview/stats
```

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "users": 1286,
    "servers": 24,
    "tools": 68,
    "today_calls": 3472
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| users | int | 平台用户总数 |
| servers | int | 已注册的 MCP 服务器数 |
| tools | int | 全部服务器的工具数合计 |
| today_calls | int | 当天（自然日）AI 调用次数，来自审计日志表（audit_logs），未配置审计库时为 0 |

---

### 7.2 获取 AI 调用历史趋势

```
GET /api/overview/call-trend
```

> 仅管理员可调用，需携带 `Authorization: Bearer <token>`。

**请求参数（Query String）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| days | int | 否 | 统计最近 N 天（含今天），默认 30，取值范围 1~365 |

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "dates": ["07-11", "07-12", "07-13"],
    "counts": [1286, 1520, 3472]
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| dates | []string | 日期标签数组，`MM-DD` 格式，按时间升序，长度恒等于请求的 days（含今天，无调用记录的天数为 0） |
| counts | []int64 | 每日 AI 调用次数，与 dates 一一对应，统计 audit_logs 全量记录数（含成功与失败），未配置审计库时全为 0 |

---

## 8. 数据模型说明

### 7.1 角色表（roles）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint (PK) | 自增主键 |
| name | string | 角色名称，唯一 |
| description | string | 角色描述 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 7.2 用户表（users）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint (PK) | 自增主键 |
| username | string | 用户名，唯一 |
| password | string | 加密后的密码 |
| role_id | uint (FK → roles.id) | 绑定的角色 ID，可空 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 7.3 角色-服务器关联表（role_servers）

| 字段 | 类型 | 说明 |
|------|------|------|
| role_id | uint (FK → roles.id) | 角色 ID |
| server_id | uint (FK → servers.id) | MCP 服务器 ID |

> 联合主键 (role_id, server_id)

### 7.4 系统设置表（settings）

| 字段 | 类型 | 说明 |
|------|------|------|
| key | string (PK) | 设置项分组（log / smtp / auth / user_ops / platform / api_gateway / network_security） |
| value | text | 该分组的 JSON 内容 |
| updated_at | datetime | 更新时间 |
