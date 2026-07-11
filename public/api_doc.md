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
