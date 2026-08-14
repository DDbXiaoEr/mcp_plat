# mcp_plat-console

MCP 服务平台 · 管理控制台

## 项目功能特点

`mcp_plat-console` 是 MCP（Model Context Protocol）服务平台的管理后台，面向**平台管理员**，提供 MCP 服务全生命周期管理、用户权限管控、API 网关集成等核心能力。

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

- **后端**: Go 1.25 + Gin + GORM + PostgreSQL/MySQL + JWT + gRPC
- **审计日志存储**: 关系库（PostgreSQL/MySQL）或 ClickHouse（可选，两库分离）
- **前端**: Vue 3 + Vite（`web/` 目录）

## 快速开始

### 1. 配置

```bash
# 编辑 config.yaml 填写数据库、JWT 等配置
# 主库支持 PostgreSQL 和 MySQL 两种数据库
```

`config.yaml` 关键配置：

```yaml
# 主库（用户、角色、AccessKey、服务器、系统设置等）
database:
  type: postgres           # 或 mysql
  postgres:
    host: localhost
    port: 5432
    user: root
    password: 123456
    dbname: mcp_platform

# 审计日志库（与 audit-log-server 共用），与主库分离
# type 可选：postgres / mysql / clickhouse
audit_log_db:
  type: postgres
  postgres:
    host: localhost
    port: 5432
    user: root
    password: 123456
    dbname: mcp_history_db
  # 切换 ClickHouse（高量级）时：
  # type: clickhouse
  # clickhouse:
  #   addr: 127.0.0.1:9000
  #   user: default
  #   password: ""
  #   db: mcp_platform
  #   table: audit_logs
  #   engine: merge_tree            # 单机：merge_tree；集群：replicated_merge_tree
  #   cluster: ""                   # replicated_merge_tree 时必填
  #   ttl_days: 90                  # 数据保留天数

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
| `make build-server` | 纯后端二进制（本地 debug） |
| `make build-web` | 仅构建前端（`web/dist/`） |
| `make build-embed` | 构建前端 + 内嵌单二进制 |
| `make build-local` | 本地全量：server + accesskey-auth-server + audit-log-server |
| `make build-local-release` | 本地发布版：embed + accesskey-auth-server + audit-log-server |
| `make build-linux-amd64-release` | 交叉编译 linux amd64（server / embed / gRPC 服务） |
| `make build-linux-arm64-release` | 交叉编译 linux arm64（embed / gRPC 服务） |
| `make build-accesskey-auth-server` | 构建 AccessKey gRPC 鉴权服务 |
| `make build-audit-log-server` | 构建审计日志 gRPC 服务 |
| `make build-apisix-plugin` | 构建 APISIX Go 插件运行器（linux/$(ARCH)） |
| `make build-tools` | 构建工具集（datagen / accesskey-test 等） |

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
- LDAP 目录服务认证，支持属性映射（登录时自动同步用户信息）
- CAS 单点登录（前端自动检测 URL 中的 ticket 参数）
- 登录方式可通过系统设置切换，支持 LDAP 连接与映射在线测试
- JWT Bearer Token 鉴权，区分普通用户和管理员权限

### 平台概况

- 数据看板：平台用户数、MCP 服务器数、工具数、当天 AI 调用次数
- AI 调用历史趋势折线图（最近 7 / 30 天），支持手动与自动刷新

### MCP 服务管理

- MCP Server 的创建、编辑、删除
- 服务信息：名称、网关路径、后端地址（支持多个 IP:Port）、所属部门、协议（SSE / Streamable HTTP / stdio）、协议版本（2025-03-26 / 2025-06-18 / 2026-07-28）
- **工具发现**：连接 MCP Server 自动调用 `tools/list` 协议拉取可用工具列表，支持下拉选择或自定义地址，多地址时随机选取
- **一键发布**：批量选择 MCP Server 自动注册为 APISIX 网关路由，多地址发布为轮询负载均衡 upstream，支持勾选 AccessKey 鉴权插件与自定义 Header
- 兼容 MCP 2025-06-18 规范，支持 Bearer Token 鉴权

### AccessKey 管理

- 用户可创建多个 API Key，用于调用 MCP 服务（数量上限可在系统设置中配置）
- 每个 Key 可指定允许访问的 MCP Server 及其工具（细粒度权限范围）
- 支持启用/禁用、到期日期设置；定时任务自动禁用过期 Key
- 配合 APISIX 网关插件实现运行时的 AccessKey 校验

### RBAC 权限

- **角色管理**：创建/编辑/删除角色，角色关联可访问的 MCP Server
- **用户管理**：管理员可创建、编辑、删除用户，分配用户到角色
- 支持穿梭框批量为用户分配角色（支持按账号搜索）

### 系统设置

系统配置按分组管理，通过页面统一编辑：

| 分组 | 配置项 |
|------|--------|
| `platform` | 平台名称、Logo URL、站点地址 |
| `api_gateway` | APISIX Admin API 地址/Key、发布域名、gRPC 鉴权地址、AccessKey Header |
| `auth` | 登录方式切换（local / ldap / cas）、LDAP/CAS/OAuth 连接参数 |
| `log` | syslog 转发（支持 TCP/UDP） |
| `smtp` | 邮件服务配置 |
| `user_ops` | 用户默认角色、AccessKey 数量上限、过期检查定时表达式 |
| `network_security` | 网络安全策略（IP 白名单） |
| `audit_log` | 审计日志开关与 audit-log-server 的 gRPC 地址 |

### APISIX 网关集成

- **Go 插件运行器**（`cmd/apisix-runner`）：在 APISIX 侧运行的自定义插件，实现请求级 AccessKey 校验
- **gRPC 鉴权服务**（`cmd/accesskey-auth-server`）：独立的 AccessKey 验证服务，供网关插件调用
- **gRPC 审计服务**（`cmd/audit-log-server`）：独立的审计日志采集服务，负责写入审计库
- MCP 服务发布时自动向 APISIX Admin API 注册路由和上游

### Kong 网关接入（测试中）

- 在系统设置中切换网关类型为 **Kong**，即可通过 Kong Admin API 发布 MCP 服务
- 仅创建/更新 **上游（Upstream）、Service、路由（Route）**，**不下发任何插件**
- 因此 Kong 场景下 Access Key 认证、审计日志、路径重写等能力不生效；网关路径使用服务器配置的 URI 路径（address）
- 本地 Docker 测试：Admin API 默认 `http://127.0.0.1:8001`，未启用 RBAC 时 Admin API Key 可留空

### 审计日志

- 审计日志存储抽象层（`auditstore/`），支持关系库（PostgreSQL/MySQL）与 ClickHouse 两种实现
- **两库分离**：主业务库与审计库独立配置（`audit_log_db`），互不影响
- ClickHouse 支持 MergeTree / ReplicatedMergeTree 引擎与 TTL 数据保留，关系库支持定期清理
- 审计日志以 AccessKey ID 存储（不落明文 Key），减少数据量

### 使用历史

- 记录每次 MCP 工具调用的使用日志（用户、AccessKey、MCP 服务、接口、状态）
- 支持按 AccessKey、MCP 服务、工具名、时间范围筛选和分页查询，支持页号跳转

### 双构建模式

- **纯后端**（`!embed` 标签）：仅编译 API 服务，前端单独部署
- **单二进制**（`embed` 标签）：前端 SPA 内嵌到 Go 二进制中，一个文件即可部署

## API

| Method | Path | 鉴权 | 说明 |
|--------|------|:---:|------|
| POST | `/api/auth/login` | - | 登录 |
| GET | `/api/auth/method` | - | 获取认证方式（local / ldap / cas） |
| POST | `/api/auth/cas/validate` | - | CAS 登录验证 |
| GET | `/api/auth/profile` | Bearer | 个人信息 |
| GET | `/api/access-keys` | Bearer | AccessKey 列表 |
| POST | `/api/access-keys` | Bearer | 创建 AccessKey |
| PUT | `/api/access-keys/:id` | Bearer | 更新 AccessKey |
| DELETE | `/api/access-keys/:id` | Bearer | 删除 AccessKey |
| GET | `/api/history` | Bearer | 使用历史（审计日志） |
| GET | `/api/audit-logs` | Bearer | 审计日志列表 |
| GET | `/api/servers` | Bearer | MCP Server 列表 |
| POST | `/api/servers` | Admin | 创建 MCP Server |
| PUT | `/api/servers/:id` | Admin | 更新 MCP Server |
| DELETE | `/api/servers/:id` | Admin | 删除 MCP Server |
| POST | `/api/servers/fetch-tools` | Admin | 拉取 tools/list |
| POST | `/api/servers/publish` | Admin | 批量发布到 API 网关 |
| GET | `/api/roles` | Admin | 角色列表 |
| POST | `/api/roles` | Admin | 创建角色 |
| PUT | `/api/roles/:id` | Admin | 更新角色 |
| DELETE | `/api/roles/:id` | Admin | 删除角色 |
| GET | `/api/roles/:id/users` | Admin | 获取角色已分配用户 |
| PUT | `/api/roles/:id/users` | Admin | 批量分配用户到角色 |
| GET | `/api/users` | Admin | 用户列表 |
| POST | `/api/users` | Admin | 创建用户 |
| PUT | `/api/users/:id` | Admin | 更新用户 |
| DELETE | `/api/users/:id` | Admin | 删除用户 |
| GET | `/api/overview/stats` | Admin | 平台概况统计 |
| GET | `/api/overview/call-trend` | Admin | AI 调用历史趋势 |
| GET | `/api/settings` | Admin | 获取全部系统设置 |
| PUT | `/api/settings/:key` | Admin | 保存系统设置 |
| POST | `/api/settings/test-ldap` | Admin | 测试 LDAP 连接与属性映射 |
| GET | `/api/settings/:key` | Bearer | 获取单个设置项 |
| GET | `/api/settings/gateway-status` | Bearer | 检查 API 网关配置状态 |

> 详细请求/响应字段见 `public/api_doc.md`，开发进度见 `public/api-devstatus.md`

## 许可证

本项目采用 [GNU Affero General Public License v3.0](LICENSE)（AGPL-3.0）开源协议。

Copyright (C) 2026 Zhaoquan Wang

