# 项目结构

```
mcp_plat-console/
├── main.go                  # 入口（!embed 标签，纯后端构建）
├── main_embed.go            # 入口（embed 标签，内嵌前端单二进制构建）
├── version.go               # Version / BuildTime 变量
├── Makefile                 # 构建脚本（build / build-server / build-embed / build-web / dev / clean）
├── go.mod / go.sum
├── config.yaml.example      # 主程序配置示例（复制为 config.yaml 使用）
├── config-generator.html    # 独立配置生成页面（HTML+纯JS，直接浏览器打开）：生成主程序/认证/审计/APISIX 配置 + K8s 清单（Namespace/ConfigMap/外部服务 Endpoints），支持中英文切换
├── AGENTS.md                # 项目整体规范
├── PROJECT_STRUCTURE.md     # 本文件 —— 项目结构速查
├── README.md
│
├── cmd/
│   ├── accesskey-auth-server/
│   │   └── main.go           # AccessKey 认证 gRPC 服务
│   ├── accesskey-test/
│   │   └── main.go           # AccessKey 测试工具
│   ├── apisix-runner/
│   │   ├── main.go           # APISIX runner 入口
│   │   └── plugins/
│   │       ├── accesskey_verify.go   # AccessKey 校验插件
│   │       └── apisix_route_example.md
│   ├── audit-log-server/
│   │   └── main.go           # 审计日志 gRPC 服务
│   ├── auditgen/
│   │   └── main.go           # 审计数据生成器（基于现有数据生成近 N 天审计日志）
│   └── datagen/
│       └── main.go           # 测试数据生成工具
│
├── config/
│   └── config.go             # 配置加载（DB / JWT / 端口 / 管理员账号）
│
├── database/
│   └── database.go           # GORM 初始化 + AutoMigrate（AuditStore 审计存储抽象入口）
│
├── auditstore/               # 审计日志存储抽象层（可配置：关系库 / ClickHouse）
│   ├── store.go              # Store 接口 + 工厂 New()
│   ├── gorm_store.go         # 关系库实现（Postgres/MySQL）+ 定期清理
│   └── clickhouse_store.go   # ClickHouse 实现（MergeTree / ReplicatedMergeTree + TTL）
│
├── logging/
│   └── logging.go            # 日志初始化（默认标准输出 + 本地文件 ./logs/mcp_plat.log 并轮转；启用 Syslog 后二者失效）
│
├── model/
│   ├── user.go               # User 模型
│   ├── access_key.go         # AccessKey 模型
│   ├── usage_history.go      # UsageHistory 模型
│   ├── mcp_server.go         # MCPServer 模型
│   ├── role.go               # Role 模型
│   ├── role_server.go        # RoleServer 关联模型
│   ├── audit_log.go          # AuditLog 模型（AccessKey 以 ID 存储，可落 ClickHouse）
│   └── setting.go            # Setting 键值模型（系统设置）
│
├── handler/
│   ├── auth.go               # POST /api/auth/login, GET /api/auth/profile
│   ├── access_key.go         # CRUD /api/access-keys
│   ├── health.go             # GET /healthz（存活探针）, GET /readyz（就绪探针，校验 DB）
│   ├── history.go            # GET /api/history
│   ├── mcp_server.go         # CRUD /api/servers
│   ├── overview.go           # GET /api/overview/stats（平台概况统计，仅管理员）, GET /api/overview/call-trend（AI 调用历史趋势，仅管理员）
│   ├── role.go               # CRUD /api/roles（仅管理员）
│   ├── rbac_user.go          # CRUD /api/users（仅管理员）
│   └── setting.go            # GET/PUT /api/settings（仅管理员）
│
├── i18n/                     # 轻量消息翻译包（词典式，无第三方依赖）
│   ├── i18n.go               # 语言检测 Lang + Translate（整句/正则模板/短语最长优先替换）
│   ├── zh_en.go              # 中文→英文翻译表（handler/middleware/service 全部直达用户消息）
│   └── i18n_test.go
│
├── resp/                     # 统一 JSON 响应助手（所有 c.JSON 出口，按请求语言翻译 message）
│   ├── resp.go               # OK / Fail / JSON / Locale
│   └── resp_test.go
│
├── service/
│   ├── auth.go               # 登录业务逻辑 + JWT 生成
│   ├── access_key.go         # AccessKey 业务逻辑
│   ├── cas.go                # CAS 认证业务逻辑
│   ├── ldap.go               # LDAP 认证业务逻辑
│   ├── history.go            # 使用历史业务逻辑
│   ├── kong.go               # Kong 网关发布（仅上游/Service/路由，不下发插件）
│   ├── mail.go               # 邮件发送（SMTP，支持 none/ssl/starttls）
│   ├── mcp_server.go         # MCPServer 业务逻辑
│   ├── mcp_server_test.go    # FetchTools 单元测试（mock + 可选真实服务器）
│   ├── notify.go             # 邮件通知统一封装（SendNotificationMail 等，SMTP 配置来自系统设置）
│   ├── overview.go           # 平台概况统计业务逻辑（用户/服务器/工具/当天调用数）
│   ├── role.go               # Role 业务逻辑
│   ├── rbac_user.go          # RBACUser 业务逻辑
│   └── setting.go            # 系统设置业务逻辑（键值 JSON 存取）
│
├── middleware/
│   ├── auth.go               # JWT Bearer Token 鉴权中间件 + AdminRequired
│   └── locale.go             # 语言中间件（解析 Accept-Language → c.Set("locale")）
│
├── plugin/
│   ├── access_key.go         # AccessKey grpc 插件注册
│   ├── accesskey_client.go   # AccessKey grpc 客户端
│   ├── accesskey_grpc.pb.go  # gRPC 生成代码
│   ├── accesskey.pb.go       # protobuf 生成代码
│   ├── accesskey.proto       # AccessKey proto 定义
│   ├── auditlog_client.go    # AuditLog gRPC 客户端
│   ├── auditlog_grpc.pb.go   # gRPC 生成代码
│   ├── auditlog.pb.go        # protobuf 生成代码
│   ├── auditlog.proto        # AuditLog proto 定义
│   └── grpc_lb.go            # gRPC 多后端负载均衡（最少连接）
│
├── router/
│   └── router.go             # Gin 路由注册（仅 main.go 使用，main_embed.go 自行注册）
│
├── public/
│   ├── api_doc.md            # API 接口文档（强制维护）
│   └── api-devstatus.md      # API 开发进度状态（强制维护）
│
├── kubernetes/               # Kubernetes 部署清单（namespace / configmap / 各服务）
│   ├── 00-namespace.yaml
│   ├── 01-external-services.yaml
│   ├── 10-configmaps.yaml
│   ├── 40-accesskey-auth-server.yaml
│   └── 50-audit-log-server.yaml
│
├── testmcp/                  # 测试 MCP 服务合集（card / network / library / academic）
│   ├── cmd/
│   │   ├── card/main.go      # 卡片服务入口
│   │   ├── network/main.go   # 网络服务入口
│   │   ├── library/main.go   # 图书馆服务入口
│   │   └── academic/main.go  # 学术服务入口
│   ├── internal/
│   │   ├── card/service.go   # 卡片服务逻辑
│   │   ├── network/service.go # 网络服务逻辑
│   │   ├── library/service.go # 图书馆服务逻辑
│   │   └── academic/service.go # 学术服务逻辑
│   ├── configs/              # 各服务配置文件
│   ├── bin/                  # 编译产物
│   ├── go.mod / go.sum
│   ├── Makefile
│   └── test_accounts.json    # 测试账号
│
└── web/                      # Vue 3 前端（SPA）
    ├── AGENTS.md             # 前端规范
    ├── PROJECT_STRUCTURE.md  # 前端项目结构
    ├── Makefile              # 前端构建脚本（install / build / dev / preview / clean）
    ├── index.html            # HTML 入口
    ├── package.json
    ├── vite.config.js        # Vite 配置（dev proxy /api → localhost:8080）
    ├── dist/                 # 构建产物（embed 打包用）
    │   ├── index.html
    │   ├── favicon.svg
    │   └── assets/
    └── src/
        ├── main.js           # Vue 入口（挂载 vue-i18n）
        ├── App.vue           # 根组件（按登录状态条件渲染）
        ├── api.js            # HTTP API 封装
        ├── i18n.js           # vue-i18n 实例（zh-CN/en-US）+ roleLabel 助手
        ├── locales/          # 国际化词条：zh-CN/（中文原文）+ en-US/（英文），按视图拆域，index.js 合并
        ├── styles/
        │   └── global.css    # 全局 CSS 变量
        ├── stores/
        │   ├── auth.js       # 鉴权状态（登录/登出/角色 code）
        │   ├── locale.js     # 语言状态（localStorage mcp-console-locale 持久化）
        │   ├── nav.js        # 侧边栏导航状态（无 Vue Router，菜单存 labelKey）
        │   ├── settings.js   # 平台设置共享态（名称/logo/跳转链接）
        │   └── theme.js      # 亮暗主题状态（localStorage 持久化）
        └── components/
            ├── LoginView.vue
            ├── TheHeader.vue
            ├── TheSidebar.vue
            ├── WelcomeView.vue
            ├── OverviewView.vue
            ├── ProfileView.vue
            ├── AccessKeyView.vue
            ├── AccessKeyDrawer.vue
            ├── ServersView.vue
            ├── HistoryView.vue
            ├── QuickAccessView.vue
            ├── RbacView.vue
            └── SettingsView.vue
```

## 构建模式

| 命令 | 构建标签 | 入口文件 | 说明 |
|------|---------|---------|------|
| `make build-server` / `go build .` | `!embed`（默认） | `main.go` | 纯后端，无前端 |
| `make build-embed` / `go build -tags embed .` | `embed` | `main_embed.go` | 内嵌 `web/dist/` 单二进制 |

## API 路由总览

| Method | Path | 鉴权 | 管理员 | Handler |
|--------|------|:---:|:---:|---------|
| GET | `/healthz` | 否 | - | handler/health.go → Healthz |
| GET | `/readyz` | 否 | - | handler/health.go → Readyz |
| POST | `/api/auth/login` | 否 | - | handler/auth.go → Login |
| GET | `/api/auth/method` | 否 | - | handler/auth.go → GetAuthMethod |
| GET | `/api/auth/platform` | 否 | - | handler/auth.go → GetPlatform |
| POST | `/api/auth/cas/validate` | 否 | - | handler/auth.go → CASValidate |
| GET | `/api/auth/profile` | 是 | - | handler/auth.go → Profile |
| PUT | `/api/auth/profile` | 是 | - | handler/auth.go → UpdateProfile |
| GET | `/api/access-keys` | 是 | - | handler/access_key.go → List |
| POST | `/api/access-keys` | 是 | - | handler/access_key.go → Create |
| PUT | `/api/access-keys/:id` | 是 | - | handler/access_key.go → Update |
| DELETE | `/api/access-keys/:id` | 是 | - | handler/access_key.go → Delete |
| GET | `/api/history` | 是 | - | handler/history.go → List（⏳ 计划中） |
| GET | `/api/overview/stats` | 是 | 是 | handler/overview.go → Stats |
| GET | `/api/overview/call-trend` | 是 | 是 | handler/overview.go → CallTrend |
| GET | `/api/servers` | 是 | - | handler/mcp_server.go → List |
| POST | `/api/servers` | 是 | 是 | handler/mcp_server.go → Create |
| PUT | `/api/servers/:id` | 是 | 是 | handler/mcp_server.go → Update |
| DELETE | `/api/servers/:id` | 是 | 是 | handler/mcp_server.go → Delete |
| POST | `/api/servers/fetch-tools` | 是 | 是 | handler/mcp_server.go → FetchTools |
| POST | `/api/servers/publish` | 是 | 是 | handler/mcp_server.go → Publish |
| POST | `/api/servers/maintenance` | 是 | 是 | handler/mcp_server.go → Maintenance |
| GET | `/api/roles` | 是 | 是 | handler/role.go → List |
| POST | `/api/roles` | 是 | 是 | handler/role.go → Create |
| PUT | `/api/roles/:id` | 是 | 是 | handler/role.go → Update |
| DELETE | `/api/roles/:id` | 是 | 是 | handler/role.go → Delete |
| GET | `/api/roles/:id/users` | 是 | 是 | handler/role.go → GetUsers |
| PUT | `/api/roles/:id/users` | 是 | 是 | handler/role.go → AssignUsers |
| GET | `/api/users` | 是 | 是 | handler/rbac_user.go → List |
| POST | `/api/users` | 是 | 是 | handler/rbac_user.go → Create |
| POST | `/api/users/batch` | 是 | 是 | handler/rbac_user.go → BatchCreate |
| PUT | `/api/users/:id` | 是 | 是 | handler/rbac_user.go → Update |
| DELETE | `/api/users/:id` | 是 | 是 | handler/rbac_user.go → Delete |
| GET | `/api/settings` | 是 | 是 | handler/setting.go → Get |
| PUT | `/api/settings/:key` | 是 | 是 | handler/setting.go → Save |
| GET | `/api/settings/:key` | 是 | - | handler/setting.go → GetByKey |
| GET | `/api/settings/gateway-status` | 是 | - | handler/setting.go → GatewayStatus |
| POST | `/api/settings/test-smtp` | 是 | 是 | handler/setting.go → TestSmtp |

> **i18n**：所有接口的 `message` 支持中/英双语，按请求头 `Accept-Language` 返回（`en*` → 英文近似翻译，其余 → 中文原文）；由 `router.Setup()` 全局注册的 `middleware.Locale()` + `resp` 响应助手实现，词条表见 `i18n/`。

## 内嵌版本前端路由（重要）

`main_embed.go` 的 `NoRoute` 处理 SPA 回退逻辑：
- 非 API 路径且对应文件不存在时，回退到 `index.html`
- 静态资源路径（`/assets/...`、`/favicon.svg`）直接返回对应文件

## 前端导航说明

前端不使用 Vue Router，通过 `stores/nav.js` 的响应式 `active` 状态切换组件：
- `App.vue` 使用 `<component :is="...">` 动态渲染
- 导航项按角色定义（admin / user）
- 用户菜单含「快速接入」页（QuickAccessView），用于为常见 MCP 客户端（如 CherryStudio）生成带 AccessKey 的接入配置；端点 URL 依据网关设置（`defaultPublishDomain` + provider）由前端拼接
