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

### 3.1 获取使用历史列表

```
GET /api/history
```

**请求参数（Query String）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| server | string | 否 | 按服务器筛选 |
| access_key_id | uint | 否 | 按 AccessKey ID 筛选 |
| start_date | string | 否 | 开始日期（YYYY-MM-DD） |
| end_date | string | 否 | 结束日期（YYYY-MM-DD） |
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20 |

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "access_key_id": 1,
        "server": "server1",
        "endpoint": "/api/tool1",
        "status": "success",
        "created_at": "2026-01-01T12:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

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
      "department": "教务处",
      "protocol": "SSE",
      "tools": "[\"查询课表\",\"成绩查询\",\"选课信息\"]",
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
| address | string | MCP 服务器地址 |
| department | string | 负责部门 |
| protocol | string | 协议类型（SSE / Streamable HTTP / stdio） |
| tools | string | 工具列表，JSON 字符串数组格式 |

### 4.2 新增 MCP 服务器

```
POST /api/servers
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 服务器名称 |
| address | string | 是 | MCP 服务器地址 |
| department | string | 否 | 负责部门 |
| protocol | string | 否 | 协议类型，默认 SSE |
| tools | string | 否 | 工具列表，JSON 字符串数组格式 |

**响应示例：**
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 1,
    "name": "教务系统 MCP",
    "address": "https://mcp.xauat.edu.cn/jwc",
    "department": "教务处",
    "protocol": "SSE",
    "tools": "[\"查询课表\",\"成绩查询\"]",
    "created_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-01-01T12:00:00Z"
  }
}
```

### 4.3 更新 MCP 服务器

```
PUT /api/servers/:id
```

**请求参数（JSON Body）：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 服务器名称 |
| address | string | 否 | MCP 服务器地址 |
| department | string | 否 | 负责部门 |
| protocol | string | 否 | 协议类型 |
| tools | string | 否 | 工具列表，JSON 字符串数组格式 |

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

## 6. 数据模型说明

### 6.1 角色表（roles）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint (PK) | 自增主键 |
| name | string | 角色名称，唯一 |
| description | string | 角色描述 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 6.2 用户表（users）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint (PK) | 自增主键 |
| username | string | 用户名，唯一 |
| password | string | 加密后的密码 |
| role_id | uint (FK → roles.id) | 绑定的角色 ID，可空 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### 6.3 角色-服务器关联表（role_servers）

| 字段 | 类型 | 说明 |
|------|------|------|
| role_id | uint (FK → roles.id) | 角色 ID |
| server_id | uint (FK → servers.id) | MCP 服务器 ID |

> 联合主键 (role_id, server_id)
