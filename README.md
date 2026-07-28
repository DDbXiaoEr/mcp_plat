# mcp_plat-console

MCP 服务平台 · 管理控制台

## 项目功能特点

`mcp_plat-console` 是 MCP（Model Context Protocol）服务平台的管理后台，面向**平台管理员**，提供 MCP 服务全生命周期管理、用户权限管控、API 网关集成等核心能力。与面向终端用户的 [mcp_plat_portal](https://github.com/anomalyco/mcp_plat-portal) 配合使用，共同构成完整的 MCP 服务平台。

- **全生命周期管理** — 从 MCP Server 注册、工具发现到发布上线，覆盖服务全流程
- **灵活认证体系** — 支持本地账号、LDAP、CAS 多种登录方式，满足不同组织的基础设施环境
- **精细化权限** — 基于角色的访问控制（RBAC），AccessKey 粒度的服务/工具级权限划分
- **MCP 工具级管控** — AccessKey 可精确控制到每个 MCP Server 的具体工具（tool）粒度，通过 `tools/list` 协议自动发现工具列表，前端勾选式控制工具权限范围，配合 APISIX 自定义插件实现运行时工具级鉴权
- **网关原生集成** — 深度集成 APISIX API 网关，一键发布路由，自定义 Go 插件实现运行时鉴权
- **单二进制部署** — 前端内嵌编译，一个可执行文件即可完成部署，简化运维

## 适用人群

| 角色 | 职责 |
|------|------|
| **超级管理员** | 管理 MCP Server、用户、角色权限、系统设置 |
| **普通用户** | 管理个人 AccessKey，查看使用历史 |

## 技术栈

- **后端**: Go 1.25 + Gin + GORM + PostgreSQL/SQLite + JWT
- **前端**: Vue 3 + Vite（`web/` 目录）

## 快速开始

### 1. 配置

```bash
# 编辑 config.yaml 填写数据库、JWT 等配置
# 支持 PostgreSQL 和 SQLite 两种数据库
```

`config.yaml` 关键配置：

```yaml
database:
  type: postgres           # 或 sqlite
  postgres:
    host: localhost
    port: 5432
    user: root
    password: 123456
    dbname: mcp_platform

jwt_secret: "at-least-32-character-random-string"
access_key_secret: "another-at-least-32-character-random-string"
server_port: 8080

admin:
  username: admin
  password: admin123
```

### 2. 启动

```bash
# 纯后端模式
make run

# 开发模式（同时启动后端 :8080 + 前端 dev server :5174）
make dev
```

前端开发服务器自动将 `/api` 请求代理到后端 `:8080`。

### 3. 构建

| 命令 | 说明 |
|------|------|
| `make build-server` | 纯后端二进制 |
| `make build-web` | 仅构建前端（`web/dist/`） |
| `make build` | 构建前端 + 内嵌单二进制 |
| `make build-embed-all` | 交叉编译 linux/windows amd64/arm64 |
| `make build-accesskey-auth-server` | 构建 gRPC 鉴权服务 |
| `make build-apisix-runner` | 构建 APISIX 插件运行器 |

构建时自动注入 Git 版本信息：

```
mcp_plat-console version v1.0.0 (commit abc1234), built at 2026-07-11_06:42:21
```

### 4. 默认管理员账号

| 用户名 | 密码 |
|--------|------|
| admin | admin123 |

## 功能特性

### 认证系统

- 本地数据库登录（bcrypt 密码哈希）
- LDAP 目录服务认证
- CAS 单点登录（前端自动检测 URL 中的 ticket 参数）
- 登录方式可通过系统设置切换
- JWT Bearer Token 鉴权，区分普通用户和管理员权限

### MCP 服务管理

- MCP Server 的创建、编辑、删除
- 服务信息：名称、网关路径、后端地址（IP:Port）、所属部门、协议（SSE / Streamable HTTP / stdio）
- **工具发现**：连接 MCP Server 自动调用 `tools/list` 协议拉取可用工具列表
- **一键发布**：将 MCP Server 自动注册为 APISIX 网关路由，支持勾选 AccessKey 鉴权插件

### AccessKey 管理

- 用户可创建多个 API Key，用于调用 MCP 服务
- 每个 Key 可指定允许访问的 MCP Server 及其工具（细粒度权限范围）
- 支持启用/禁用、到期日期设置
- 配合 APISIX 网关插件实现运行时的 AccessKey 校验

### RBAC 权限

- **角色管理**：创建/编辑/删除角色，角色关联可访问的 MCP Server
- **用户管理**：管理员可创建、编辑、删除用户，分配用户到角色
- 支持批量为用户分配角色

### 系统设置

系统配置按分组管理，通过页面统一编辑：

| 分组 | 配置项 |
|------|--------|
| `platform` | 平台名称、Logo URL、站点地址 |
| `api_gateway` | APISIX Admin API 地址/Key、发布域名、gRPC 鉴权地址 |
| `auth` | 登录方式切换（local / ldap / cas） |
| `log` | syslog 转发（支持 TCP/UDP） |
| `smtp` | 邮件服务配置 |
| `network_security` | 网络安全策略 |

### APISIX 网关集成

- **Go 插件运行器**（`cmd/apisix-runner`）：在 APISIX 侧运行的自定义插件，实现请求级 AccessKey 校验
- **gRPC 鉴权服务**（`cmd/accesskey-auth-server`）：独立的 AccessKey 验证服务，供网关插件调用
- MCP 服务发布时自动向 APISIX Admin API 注册路由和上游

### 使用历史

- 记录每次 MCP 工具调用的使用日志（用户、AccessKey、MCP 服务、接口、状态）
- 支持按 MCP 服务、时间范围筛选和分页查询

### 双构建模式

- **纯后端**（`!embed` 标签）：仅编译 API 服务，前端单独部署
- **单二进制**（`embed` 标签）：前端 SPA 内嵌到 Go 二进制中，一个文件即可部署

## API

| Method | Path | 鉴权 | 说明 |
|--------|------|:---:|------|
| POST | `/api/auth/login` | - | 登录 |
| GET | `/api/auth/profile` | Bearer | 个人信息 |
| GET | `/api/access-keys` | Bearer | AccessKey 列表 |
| POST | `/api/access-keys` | Bearer | 创建 AccessKey |
| PUT | `/api/access-keys/:id` | Bearer | 更新 AccessKey |
| DELETE | `/api/access-keys/:id` | Bearer | 删除 AccessKey |
| GET | `/api/history` | Bearer | 使用历史 |
| GET | `/api/servers` | Admin | MCP Server 列表 |
| POST | `/api/servers` | Admin | 创建 MCP Server |
| PUT | `/api/servers/:id` | Admin | 更新 MCP Server |
| DELETE | `/api/servers/:id` | Admin | 删除 MCP Server |
| POST | `/api/servers/:id/fetch-tools` | Admin | 拉取 tools/list |
| POST | `/api/servers/:id/publish` | Admin | 发布到 API 网关 |
| GET | `/api/roles` | Admin | 角色列表 |
| POST | `/api/roles` | Admin | 创建角色 |
| PUT | `/api/roles/:id` | Admin | 更新角色 |
| DELETE | `/api/roles/:id` | Admin | 删除角色 |
| PUT | `/api/roles/:id/users` | Admin | 分配用户到角色 |
| GET | `/api/users` | Admin | 用户列表 |
| POST | `/api/users` | Admin | 创建用户 |
| PUT | `/api/users/:id` | Admin | 更新用户 |
| DELETE | `/api/users/:id` | Admin | 删除用户 |
| GET | `/api/settings` | Bearer | 获取系统设置 |
| PUT | `/api/settings` | Admin | 更新系统设置 |

> 详细请求/响应字段见 `public/api_doc.md`，开发进度见 `public/api-devstatus.md`

