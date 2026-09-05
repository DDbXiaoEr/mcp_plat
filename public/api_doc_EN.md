# API Documentation (English)

Base URL: `http://localhost:8080`

Common response format:
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

> **i18n**: the `message` field is bilingual. The client declares its preference with the `Accept-Language` request header:
> - When it carries `en` / `en-*` (e.g. `Accept-Language: en-US`), `message` is returned in English (approximate dictionary translation; unmatched concatenated fragments keep their original text);
> - When absent or Chinese (`zh` / `zh-*`), the original Chinese text is returned.
> Business `data` is user data and is not translated; the few backend-generated fallback labels inside `data` (e.g. LDAP test result `message`, or overview-chart labels like `未知服务器`/`其他`/`未分组`) are localized with the language.

---

## 1. Authentication Module

### 1.1 Login

```
POST /api/auth/login
```

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| username | string | yes | username |
| password | string | yes | password |

**Success response:**
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

### 1.2 Get Personal Profile

```
GET /api/auth/profile
```

**Request headers:**
| Field | Description |
|-------|-------------|
| Authorization | Bearer {token} |

**Success response:**
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

### 1.3 Update Personal Profile (email)

```
PUT /api/auth/profile
```

**Request headers:**
| Field | Description |
|-------|-------------|
| Authorization | Bearer {token} |

> Only regular users (non-admin) may call this; admin accounts get a 403. After changing the email, an LDAP user's email may still be overwritten/synced by the LDAP attribute mapping on the next login.

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | yes | new email address; must be well-formed and not used by another account |

**Success response:**
```json
{
  "code": 200,
  "message": "保存成功",
  "data": {
    "id": 1,
    "uid": "2021001",
    "username": "user1",
    "email": "new@example.com",
    "phone": "13800000000",
    "organization": "计算机学院",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

**Failure response (malformed / already-used / empty email):**
```json
{
  "code": 400,
  "message": "邮箱格式不正确"
}
```

### 1.4 Get Auth Method

```
GET /api/auth/method
```

**Success response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "method": "cas",
    "cas": {
      "serverUrl": "https://cas.example.edu.cn",
      "serviceUrl": "http://localhost:5174",
      "version": "3.0",
      "attrMapping": { "name": "displayName", "email": "mail" }
    }
  }
}
```

Possible `method` values: `local` (local database), `ldap`, `cas`. `cas.attrMapping` is the attribute mapping (platform field → CAS-returned attribute name); platform fields can be `uid` / `name` / `email` / `phone` / `organization`, used to create/update the CAS user profile. CAS has no separate test-mapping endpoint.

### 1.5 Get Platform Info (login page)

```
GET /api/auth/platform
```

No auth required; lets the login page fetch public info such as platform name, logo and login background image.

**Success response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "name": "某某大学",
    "logoUrl": "https://www.example.edu.cn/logo.png",
    "siteUrl": "https://www.example.edu.cn",
    "loginBackground": "https://www.example.edu.cn/bg.jpg"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| name | string | platform name |
| logoUrl | string | logo image URL |
| siteUrl | string | jump link |
| loginBackground | string | login-page background image URL (empty → default gradient background) |

### 1.6 CAS Login Validation

```
POST /api/auth/cas/validate
```

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| ticket | string | yes | ticket returned by CAS |
| serviceUrl | string | yes | callback URL; must match the one used when redirecting to CAS |

**Success response:**
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

## 2. AccessKey Management

All endpoints below require the auth header: `Authorization: Bearer {token}`

### 2.1 List AccessKeys

```
GET /api/access-keys
```

**Success response:**
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

### 2.2 Create AccessKey

```
POST /api/access-keys
```

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | yes | AccessKey name |
| servers | string | no | server-permission config (JSON string) |
| expired_at | string | no | expiry time (ISO 8601) |

**Success response:**
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

### 2.3 Update AccessKey

```
PUT /api/access-keys/:id
```

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | no | AccessKey name |
| enabled | bool | no | whether enabled |
| servers | string | no | server-permission config (JSON string) |
| expired_at | string | no | expiry time (ISO 8601) |

**Success response:**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

### 2.4 Delete AccessKey

```
DELETE /api/access-keys/:id
```

**Success response:**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

## 3. Usage History

> ✅ Implemented: based on audit logs (`audit_logs`). The storage backend is configurable (relational Postgres/MySQL or ClickHouse), see `audit_log_db` in `config.yaml`.

### 3.1 List Usage History

```
GET /api/history
```

**Request params (Query String):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| access_key | string | no | filter by AccessKey (pass the full key; matched after server-side hashing) |
| server_id | string | no | filter by MCP server |
| tool_name | string | no | filter by tool name |
| start | string | no | start date (YYYY-MM-DD) |
| end | string | no | end date (YYYY-MM-DD) |
| page | int | no | page number, default 1 |
| page_size | int | no | page size, default 20 (max 100) |

**Success response:**
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
        "client_ip": "10.0.0.8",
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

> Note: `access_key_id` is the AccessKey's database ID (the plaintext key is not stored, reducing data volume); `access_key_name` is resolved server-side from the keys the current user holds and may be empty once the key is deleted. `client_ip` is the caller's source IP (collected by the APISIX `audit_log` plugin; behind a gateway it may be the gateway/reverse-proxy egress address; empty for legacy records or when the audit plugin is disabled).

---

## 4. MCP Server Management

> All endpoints require auth; send `Authorization: Bearer <token>` in the request header.

### 4.1 List MCP Servers

```
GET /api/servers
```

**Request params:** none

**Response example:**
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

| Field | Type | Description |
|-------|------|-------------|
| id | uint | server ID |
| name | string | server name |
| address | string | MCP server address (API-gateway path or full URL) |
| service_address | string | MCP service addresses as a JSON-array string (e.g. `["192.168.1.100:8081","192.168.1.101:8081"]`); supports multiple addresses, used as API-gateway upstream nodes when published |
| department | string | responsible department |
| protocol | string | protocol type (SSE / Streamable HTTP / stdio) |
| protocol_version | string | MCP protocol version (2025-03-26 / 2025-06-18 / 2026-07-28), default 2026-07-28 |
| tools | string | tool list, JSON-array-string format; each element is `{"name":"...","description":"..."}` |
| description | string | server description |
| status | string | publish status: `published`, `maintenance` or `unpublished` |
| auth_enabled | bool | whether Access-Key auth is enabled (recorded from `enable_auth` when publishing; always `false` for Kong-gateway publishing; default `true` for unpublished servers) |

### 4.2 Create MCP Server

```
POST /api/servers
```

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | yes | server name |
| address | string | yes | MCP server address (API-gateway path or full URL) |
| service_address | string | no | MCP service addresses, JSON-array string (e.g. `["192.168.1.100:8081","192.168.1.101:8081"]`); supports multiple addresses |
| department | string | no | responsible department |
| protocol | string | no | protocol type, default SSE |
| protocol_version | string | no | MCP protocol version, default 2026-07-28 |
| tools | string | no | tool list, JSON-array-string format |
| description | string | no | server description |

**Response example:**
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

### 4.5 Publish Servers to the API Gateway

```
POST /api/servers/publish
```

> Admin permission required.

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| server_ids | []string | yes | list of MCP server IDs to publish |
| enable_auth | bool | no | whether to enable Access-Key auth (default false); when enabled the route gets the accesskey_verify plugin |
| enable_audit_log | bool | no | whether to enable audit logging (default false); when enabled the route gets the audit_log plugin |
| accesskey_header | string | no | custom HTTP header name for the Access Key; when omitted, the value configured in gateway settings is preferred; if none is configured, default `X-Access-Key` |

**Request example:**
```json
{
  "server_ids": ["uuid-1", "uuid-2"],
  "enable_auth": true
}
```

**Success response:**
```json
{
  "code": 200,
  "message": "发布成功"
}
```

**Notes:** the backend registers routes for the selected MCP servers to the configured API gateway (settings group `api_gateway`, field `provider`). The gateway Admin API address must be configured first (APISIX also needs its Key). After a successful publish the server is marked `published` and `auth_enabled` is recorded from `enable_auth` (always `false` for Kong).
- `provider=apisix` (default): creates upstream + route, optionally pushing plugins such as `accesskey_verify` / `audit_log` / `proxy-rewrite`. If auth is enabled, `authGrpcAddrs` must also be configured (an array; multiple addresses are supported and load-balanced by connection count) pointing to the accesskey-auth-server gRPC addresses.
- `provider=kong`: only creates upstream + Service + Route and **pushes no plugins**, so auth, audit and path-rewrite do not take effect under Kong; the gateway path uses the URI path configured on the server (`address`). The Admin API Key is only needed when Kong RBAC is enabled.

### 4.6 Fetch MCP Server Tools

```
POST /api/servers/fetch-tools
```

> Requires auth; send `Authorization: Bearer <token>` in the request header.

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| address | string | yes | MCP server address |
| protocol | string | no | protocol type (SSE / Streamable HTTP / stdio), default Streamable HTTP |
| protocol_version | string | no | MCP protocol version, default 2026-07-28 |

**Response example:**
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

**Notes:** the backend connects to the target MCP server by address and protocol, calls the `tools/list` protocol method and returns the available tool names and descriptions. A connection failure or invalid address returns a 4xx/5xx.

### 4.7 Set Server Maintenance Status

```
POST /api/servers/maintenance
```

> Admin permission required.

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| server_ids | []string | no | MCP server IDs to enter maintenance; their gateway routes return 503 directly |
| restore_ids | []string | no | MCP server IDs to leave maintenance; their gateway routes restore the original config |

**Request example:**
```json
{
  "server_ids": ["uuid-1"],
  "restore_ids": ["uuid-2"]
}
```

**Success response:**
```json
{
  "code": 200,
  "message": "维护设置成功"
}
```

**Notes:** the backend changes the selected MCP server routes so they return 503 directly, based on the configured API gateway settings (`api_gateway` group, `provider` field):
- `provider=apisix` (default): the route `plugins` are fully replaced with the `mocking` plugin (`response_status=503`; older APISIX falls back to the `mock` plugin `response_code=503`). Leaving maintenance restores the original route config saved at publish time (including auth/audit plugins).
- `provider=kong`: the `request-termination` plugin is pushed to the route (`status_code=503`); leaving maintenance removes it.
- On entering maintenance the server `status` is set to `maintenance`; on restore it goes back to `published`.

### 4.3 Update MCP Server

```
PUT /api/servers/:id
```

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | no | server name |
| address | string | no | MCP server address |
| service_address | string | no | MCP service addresses, JSON-array string; supports multiple addresses |
| department | string | no | responsible department |
| protocol | string | no | protocol type |
| protocol_version | string | no | MCP protocol version |
| tools | string | no | tool list, JSON-array-string format |
| description | string | no | server description |

**Response example:**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

### 4.4 Delete MCP Server

```
DELETE /api/servers/:id
```

**Request params:** none (ID passed in the URL path)

**Response example:**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

## 5. RBAC Management

> All endpoints require auth — `Authorization: Bearer <token>` in the request header — and only admins may call them.

### 5.1 Role Management

#### 5.1.1 List Roles

```
GET /api/roles
```

**Success response:**
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

| Field | Type | Description |
|-------|------|-------------|
| id | uint | role ID |
| name | string | role name |
| description | string | role description |
| server_ids | []uint | accessible MCP server ID list |
| user_count | int | number of users bound to this role |
| created_at | string | creation time |
| updated_at | string | update time |

#### 5.1.2 Create Role

```
POST /api/roles
```

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | yes | role name |
| description | string | no | role description |
| server_ids | []uint | no | accessible MCP server ID list |

**Success response:**
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

#### 5.1.3 Update Role

```
PUT /api/roles/:id
```

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | no | role name |
| description | string | no | role description |
| server_ids | []uint | no | accessible MCP server ID list |

**Success response:**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

#### 5.1.4 Delete Role

```
DELETE /api/roles/:id
```

> Note: a role that still has bound users cannot be deleted; unbind users first.

**Success response:**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

### 5.2 User Management

#### 5.2.1 List Users

```
GET /api/users
```

**Request params (Query String):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| role_id | uint | no | filter by role |
| q | string | no | fuzzy search by student/staff number (uid) |

**Success response:**
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

| Field | Type | Description |
|-------|------|-------------|
| id | uint | user ID |
| username | string | username |
| role_id | uint/null | bound role ID; null = unassigned |
| role_name | string | bound role name; empty string when unassigned |
| created_at | string | creation time |
| updated_at | string | update time |

#### 5.2.2 Create User

```
POST /api/users
```

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| username | string | yes | username |
| password | string | yes | login password |
| role_id | uint | no | bound role ID |

**Success response:**
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

#### 5.2.3 Update User

```
PUT /api/users/:id
```

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| username | string | no | username |
| password | string | no | new password |
| role_id | uint | no | bound role ID (pass null or 0 to unbind the role) |

**Success response:**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

#### 5.2.4 Delete User

```
DELETE /api/users/:id
```

**Success response:**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

#### 5.2.5 Batch Create Users

```
POST /api/users/batch
```

> Creates multiple users in one request. If any username/student-number/email in the list is duplicated or conflicts with an existing account, the whole batch creates nothing (transaction rollback) and the specific conflicts are returned so you can fix and resubmit.

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| role_id | uint/null | no | role ID uniformly bound for the batch; null = none |
| users | []object | yes | users to create; cannot be empty |

**users[]. fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| username | string | yes | username (login account), unique |
| password | string | yes | login password |
| name | string | no | full name |
| email | string | no | email (must pass format validation if non-empty; must not duplicate an existing account) |
| uid | string | no | student/staff number, unique; auto-generated when left empty |
| phone | string | no | phone |
| organization | string | no | unit/organization |

**Request example:**
```json
{
  "role_id": 2,
  "users": [
    {
      "username": "newuser1",
      "password": "pass123",
      "name": "张三",
      "email": "zhangsan@example.edu.cn",
      "uid": "20260001",
      "phone": "13800000001",
      "organization": "计算机学院"
    },
    {
      "username": "newuser2",
      "password": "pass456",
      "email": "lisi@example.edu.cn"
    }
  ]
}
```

**Success response:**
```json
{
  "code": 200,
  "message": "批量创建成功",
  "data": [
    {
      "id": 8,
      "uid": "20260001",
      "username": "newuser1",
      "role_id": 2,
      "role_name": "普通用户",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

> `data` is the array of created users; each item has the same shape as a single item in "List Users". On conflict failures the `message` looks like `用户名 'xxx' 已存在` / `邮箱 'xxx' 已被其他账号使用` / `学号/工号 'xxx' 已被使用` / `列表中存在重复的用户名 'xxx'`.

---

### 5.3 Batch Assign Role to Users

```
PUT /api/roles/:id/users
```

> Pass the target user ID list; the role is fully (re)assigned to these users. Users not in the list are unbound from this role.

**Request params (JSON Body):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| user_ids | []uint | yes | user ID list to assign to this role |

**Request example:**
```json
{
  "user_ids": [1, 2, 5, 8]
}
```

**Success response:**
```json
{
  "code": 200,
  "message": "分配成功"
}
```

---

### 5.4 List Users of a Role

```
GET /api/roles/:id/users
```

**Success response:**
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

## 6. System Settings

> The endpoints below require admin permission (`Authorization: Bearer <token>` and the user is an admin).

### 6.1 Get All Settings

```
GET /api/settings
```

**Success response:**

`data` is an object grouped by settings item; it only contains groups that have been saved. Each group's content is the raw JSON saved at the time.

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "log": {
      "syslogEnabled": false,
      "logPath": "./logs",
      "logLevel": "info",
      "logPrefix": "mcp_plat",
      "maxSize": 100,
      "maxBackups": 10,
      "maxAge": 30,
      "compress": false,
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
      "cas": { "serverUrl": "", "serviceUrl": "", "version": "3.0", "attrMapping": {} },
      "ldap": { "host": "", "port": 389, "baseDn": "", "bindDn": "", "bindPassword": "", "userFilter": "", "attrMapping": {} },
      "oauth": { "authorizeUrl": "", "tokenUrl": "", "userinfoUrl": "", "clientId": "", "clientSecret": "", "redirectUrl": "", "scope": "" }
    },
    "user_ops": {
      "filterAttribute": "organization",
      "roleRules": [
        { "roleId": 2, "pattern": "^计算机学院" },
        { "roleId": 3, "pattern": "^理学院" }
      ],
      "maxAccessKeys": 5,
      "accessKeyCron": "0 */6 * * *"
    },
    "platform": {
      "name": "某某大学",
      "logoUrl": "https://www.example.edu.cn/logo.png",
      "siteUrl": "https://www.example.edu.cn",
      "loginBackground": "https://www.example.edu.cn/bg.jpg"
    },
    "api_gateway": {
      "provider": "apisix",
      "adminUrl": "http://127.0.0.1:9180",
      "adminKey": "edd1c9f034335f136f87ad84b625c8f1",
      "defaultPublishDomain": "mcp.xauat.edu.cn",
      "authGrpcAddrs": [":9090", "127.0.0.1:9091"],
      "accesskeyHeader": "X-Access-Key"
    },
    "audit_log": {
      "enabled": true,
      "grpcAddrs": [":9091", "127.0.0.1:9191"]
    },
    "network_security": {
      "allowlist": ["10.0.0.0/8", "172.16.0.0/12", "192.168.1.0/24"]
    }
  }
}
```

### 6.2 Save Settings

```
PUT /api/settings/:key
```

**Path params:**

| Param | Description |
|-------|-------------|
| key | settings group; supports `log` / `smtp` / `auth` / `user_ops` / `platform` / `api_gateway` / `network_security` / `audit_log` / `quick_access` |

**Request params (JSON Body):**

Any valid JSON object; the whole group is overwritten and saved (stored verbatim, fields are not validated).

> The `log` group takes effect immediately on save (no restart needed). When Syslog is disabled (`syslogEnabled=false` or `syslogHost` empty), logs go to both stdout and the local file `{logPath}/{logPrefix}.log` (default `./logs/mcp_plat.log`) with rotation by `maxSize` (per-file size cap in MB, default 100), `maxBackups` (number of retained history files, default 10), `maxAge` (retention days, default 30) and `compress` (gzip or not). When Syslog is enabled, stdout and file logging are both disabled and logs go only to the Syslog server (`syslogProtocol` supports tcp/udp, default port 514, `logPrefix` is used as the tag); on connection failure it falls back to stdout + file.

**Success response:**
```json
{
  "code": 200,
  "message": "保存成功"
}
```

**Failure response:**
```json
{
  "code": 400,
  "message": "不支持的设置项"
}
```

---

### 6.3 Get a Single Settings Item

```
GET /api/settings/:key
```

> Requires login only; no admin permission needed.

**Path params:**

| Param | Description |
|-------|-------------|
| key | settings group; supports `log` / `smtp` / `auth` / `user_ops` / `platform` / `api_gateway` / `network_security` / `audit_log` / `quick_access` |

**Response example (`GET /api/settings/platform`):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "name": "某某大学",
    "logoUrl": "https://www.example.edu.cn/logo.png",
    "siteUrl": "https://www.example.edu.cn",
    "loginBackground": "https://www.example.edu.cn/bg.jpg"
  }
}
```

---

### 6.4 Check API Gateway Config Status

```
GET /api/settings/gateway-status
```

> Requires login only; no admin permission needed.

**Response example:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "configured": true,
    "provider": "kong",
    "adminUrl": "http://127.0.0.1:8001",
    "defaultPublishDomain": "mcp.xauat.edu.cn",
    "accesskeyHeader": "X-Access-Key"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| configured | bool | whether the gateway is configured (APISIX needs URL + Key; Kong needs only URL) |
| provider | string | current gateway type: `apisix` / `kong` / `tyk` |
| adminUrl | string | configured Admin API base URL |
| defaultPublishDomain | string | default publish domain, used to build MCP access URLs (empty string when not configured) |
| accesskeyHeader | string | HTTP header name used for Access-Key auth (empty string when not configured; the frontend falls back to `X-Access-Key`) |

> Note: the endpoint never returns gateway admin credentials (e.g. `adminKey`); it only exposes the fields needed to build access URLs.

---

### 6.5 Quick Access Settings (quick_access)

> Admin permission required to save (`PUT /api/settings/quick_access`); regular users can read it via `GET /api/settings/quick_access` (login only).

**Save example (JSON Body):**
```json
{
  "enabledClients": ["cherrystudio"],
  "scheme": "https"
}
```

| Field | Type | Description |
|-------|------|-------------|
| enabledClients | []string | clients exposed on the user "Quick Access" page; currently supports `cherrystudio`; empty = no client exposed |
| scheme | string | protocol used to build access URLs: `http` / `https`, default `https` |

---

### 6.6 User Operations Settings (user_ops)

> Admin permission required to save (`PUT /api/settings/user_ops`).

**Save example (JSON Body):**
```json
{
  "filterAttribute": "organization",
  "roleRules": [
    { "roleId": 2, "pattern": "^计算机学院" },
    { "roleId": 3, "pattern": "^理学院" }
  ],
  "maxAccessKeys": 5,
  "accessKeyCron": "0 */6 * * *"
}
```

| Field | Type | Description |
|-------|------|-------------|
| filterAttribute | string | the single global filter attribute shared by all role rules; options: `username` / `uid` / `name` / `email` / `phone` / `organization` |
| roleRules | []object | auto role-assignment rules for new users, matched in order; the first match assigns its role |
| roleRules[].roleId | uint | target role ID |
| roleRules[].pattern | string | regular expression matched against the attribute named by `filterAttribute` (Go `regexp` semantics) |
| maxAccessKeys | int | maximum AccessKeys per user |
| accessKeyCron | string | cron expression for the expired-AccessKey scan |

> Note: for a new user (created on first LDAP/CAS login), the value of the attribute named by `filterAttribute` is taken and matched in order against `roleRules`; the first match auto-assigns the role. If no rule matches, the user's `role_id` stays empty (unassigned) until an admin assigns it manually in user management.

---

### 6.7 Email Notification Settings (smtp)

> Admin permission required to save (`PUT /api/settings/smtp`). Base config for the email-notification feature; later concrete scenarios (login alerts, AccessKey-expiry reminders, etc.) send mail based on this config.

**Save example (JSON Body):**
```json
{
  "enabled": true,
  "host": "smtp.example.edu.cn",
  "port": 465,
  "encryption": "ssl",
  "username": "noreply@example.edu.cn",
  "password": "********",
  "fromAddress": "noreply@example.edu.cn",
  "fromName": "MCP 服务平台"
}
```

| Field | Type | Description |
|-------|------|-------------|
| enabled | bool | whether email notification is enabled |
| host | string | SMTP server address |
| port | int | port, default 465 |
| encryption | string | encryption: `none` / `ssl` / `starttls` |
| username | string | SMTP login username |
| password | string | SMTP login password |
| fromAddress | string | sender address (empty → uses username) |
| fromName | string | sender name |

---

### 6.8 Send Test Email

```
POST /api/settings/test-smtp
```

> Admin permission required. Verifies that the SMTP config can send mail; the request carries the full smtp config (same shape as save) and does not require saving first.

**Request params (JSON Body):**
```json
{
  "smtp": {
    "host": "smtp.example.edu.cn",
    "port": 465,
    "encryption": "ssl",
    "username": "noreply@example.edu.cn",
    "password": "********",
    "fromAddress": "noreply@example.edu.cn",
    "fromName": "MCP 服务平台"
  },
  "to": "admin@example.edu.cn"
}
```

| Field | Type | Description |
|-------|------|-------------|
| smtp | object | SMTP config identical to `PUT /api/settings/smtp` |
| to | string | recipient email for the test mail |

**Success response:**
```json
{
  "code": 200,
  "message": "测试邮件发送成功"
}
```

**Failure response (e.g. bad config, cannot connect):**
```json
{
  "code": 400,
  "message": "邮件发送失败: ..."
}
```

---

## 7. Platform Overview Stats

> Admin only; requires `Authorization: Bearer <token>`.

### 7.1 Get Platform Overview Stats

```
GET /api/overview/stats
```

**Success response:**
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

| Field | Type | Description |
|-------|------|-------------|
| users | int | total platform users |
| servers | int | registered MCP server count |
| tools | int | total tools across all servers |
| today_calls | int | AI calls for the current (calendar) day, from the audit-log table (`audit_logs`); 0 when no audit DB is configured |

---

### 7.2 Get AI Call History Trend

```
GET /api/overview/call-trend
```

> Admin only; requires `Authorization: Bearer <token>`.

**Request params (Query String):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| days | int | no | last N days (incl. today), default 30, range 1~365 |

**Success response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "dates": ["07-11", "07-12", "07-13"],
    "counts": [1286, 1520, 3472],
    "server_calls": [
      { "name": "统一身份认证", "count": 3210 },
      { "name": "教务管理系统", "count": 2458 }
    ],
    "server_names": ["统一身份认证", "教务管理系统", "其他"],
    "server_counts_by_day": [
      [1024, 220, 42],
      [1180, 300, 40],
      [1006, 310, 46]
    ],
    "user_groups": [
      { "name": "普通用户", "count": 980 },
      { "name": "管理员", "count": 5 },
      { "name": "未分组", "count": 12 }
    ]
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| dates | []string | date label array, `MM-DD` format, ascending; length always equals the requested days (incl. today; days without calls are 0) |
| counts | []int64 | daily AI call counts, one-to-one with `dates`; counts all `audit_logs` rows (success + failure); all 0 when no audit DB is configured |
| server_calls | []{name,count} | per-server (MCP server) AI calls over the period; `name` is the server name (falls back to 「未知服务器」 when the server_id is not found in the servers table); descending by calls; empty array when no audit DB is configured |
| server_names | []string | stacked-bar series names; takes the top 8 servers by total calls in the period, the rest are merged into 「其他」 (shortened when there are fewer than 8 or no 「其他」); empty array when no audit DB is configured |
| server_counts_by_day | [][]int64 | per-day × per-series call matrix; row count matches `dates`, column count matches `server_names`; `server_counts_by_day[i][j]` is the call count of series j on day i; empty array when no audit DB is configured |
| user_groups | []{name,count} | user-count distribution by role (`roles`); users without a role go into 「未分组」; descending by count |

---

## 8. Health Checks

### 8.1 Liveness Probe

```
GET /healthz
```

No auth; returns 200 while the process is alive; used by Docker / K8s liveness probes.

**Success response:**
```json
{
  "code": 200,
  "message": "ok"
}
```

### 8.2 Readiness Probe

```
GET /readyz
```

No auth; returns 200 only when the process is alive and the database is reachable, otherwise 503; used by Docker / K8s readiness probes.

**Success response:**
```json
{
  "code": 200,
  "message": "ok"
}
```

**Failure response (database unreachable):**
```json
{
  "code": 503,
  "message": "database unavailable"
}
```

---

## 9. Data Model Notes

### 7.1 Role table (roles)

| Field | Type | Description |
|-------|------|-------------|
| id | uint (PK) | auto-increment primary key |
| name | string | role name, unique |
| description | string | role description |
| created_at | datetime | creation time |
| updated_at | datetime | update time |

### 7.2 User table (users)

| Field | Type | Description |
|-------|------|-------------|
| id | uint (PK) | auto-increment primary key |
| username | string | username, unique |
| password | string | encrypted password |
| role_id | uint (FK → roles.id) | bound role ID, nullable |
| created_at | datetime | creation time |
| updated_at | datetime | update time |

### 7.3 Role-Server association table (role_servers)

| Field | Type | Description |
|-------|------|-------------|
| role_id | uint (FK → roles.id) | role ID |
| server_id | uint (FK → servers.id) | MCP server ID |

> Composite primary key (role_id, server_id)

### 7.4 System Settings table (settings)

| Field | Type | Description |
|-------|------|-------------|
| key | string (PK) | settings group (log / smtp / auth / user_ops / platform / api_gateway / network_security / quick_access) |
| value | text | JSON content of the group |
| updated_at | datetime | update time |
