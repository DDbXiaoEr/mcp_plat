# 项目结构

```
mcp_plat-console/
├── main.go                  # 入口（!embed 标签，纯后端构建）
├── main_embed.go            # 入口（embed 标签，内嵌前端单二进制构建）
├── version.go               # Version / BuildTime 变量
├── Makefile                 # 构建脚本（build / build-server / build-embed / build-web / dev / clean）
├── go.mod / go.sum
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
│   └── datagen/
│       └── main.go           # 测试数据生成工具
│
├── config/
│   └── config.go             # 配置加载（DB / JWT / 端口 / 管理员账号）
│
├── database/
│   └── database.go           # GORM 初始化 + AutoMigrate
│
├── logging/
│   └── logging.go            # 日志初始化（启动时读 log 设置，启用则重定向到 Syslog，否则标准输出）
│
├── model/
│   ├── user.go               # User 模型
│   ├── access_key.go         # AccessKey 模型
│   ├── usage_history.go      # UsageHistory 模型
│   ├── mcp_server.go         # MCPServer 模型
│   ├── role.go               # Role 模型
│   ├── role_server.go        # RoleServer 关联模型
│   └── setting.go            # Setting 键值模型（系统设置）
│
├── handler/
│   ├── auth.go               # POST /api/auth/login, GET /api/auth/profile
│   ├── access_key.go         # CRUD /api/access-keys
│   ├── history.go            # GET /api/history
│   ├── mcp_server.go         # CRUD /api/servers
│   ├── role.go               # CRUD /api/roles（仅管理员）
│   ├── rbac_user.go          # CRUD /api/users（仅管理员）
│   └── setting.go            # GET/PUT /api/settings（仅管理员）
│
├── service/
│   ├── auth.go               # 登录业务逻辑 + JWT 生成
│   ├── access_key.go         # AccessKey 业务逻辑
│   ├── cas.go                # CAS 认证业务逻辑
│   ├── ldap.go               # LDAP 认证业务逻辑
│   ├── history.go            # 使用历史业务逻辑
│   ├── mcp_server.go         # MCPServer 业务逻辑
│   ├── mcp_server_test.go    # FetchTools 单元测试（mock + 可选真实服务器）
│   ├── role.go               # Role 业务逻辑
│   ├── rbac_user.go          # RBACUser 业务逻辑
│   └── setting.go            # 系统设置业务逻辑（键值 JSON 存取）
│
├── middleware/
│   └── auth.go               # JWT Bearer Token 鉴权中间件 + AdminRequired
│
├── plugin/
│   ├── access_key.go         # AccessKey grpc 插件注册
│   ├── accesskey_client.go   # AccessKey grpc 客户端
│   ├── accesskey_grpc.pb.go  # gRPC 生成代码
│   └── accesskey.pb.go       # protobuf 生成代码
│
├── router/
│   └── router.go             # Gin 路由注册（仅 main.go 使用，main_embed.go 自行注册）
│
├── public/
│   ├── api_doc.md            # API 接口文档（强制维护）
│   └── api-devstatus.md      # API 开发进度状态（强制维护）
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
        ├── main.js           # Vue 入口
        ├── App.vue           # 根组件（按登录状态条件渲染）
        ├── api.js            # HTTP API 封装
        ├── styles/
        │   └── global.css    # 全局 CSS 变量
        ├── stores/
        │   ├── auth.js       # 鉴权状态（登录/登出/角色）
        │   ├── nav.js        # 侧边栏导航状态（无 Vue Router）
        │   └── settings.js   # 平台设置共享态（名称/logo/跳转链接）
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
| POST | `/api/auth/login` | 否 | - | handler/auth.go → Login |
| GET | `/api/auth/method` | 否 | - | handler/auth.go → GetAuthMethod |
| POST | `/api/auth/cas/validate` | 否 | - | handler/auth.go → CASValidate |
| GET | `/api/auth/profile` | 是 | - | handler/auth.go → Profile |
| GET | `/api/access-keys` | 是 | - | handler/access_key.go → List |
| POST | `/api/access-keys` | 是 | - | handler/access_key.go → Create |
| PUT | `/api/access-keys/:id` | 是 | - | handler/access_key.go → Update |
| DELETE | `/api/access-keys/:id` | 是 | - | handler/access_key.go → Delete |
| GET | `/api/history` | 是 | - | handler/history.go → List（⏳ 计划中） |
| GET | `/api/servers` | 是 | - | handler/mcp_server.go → List |
| POST | `/api/servers` | 是 | 是 | handler/mcp_server.go → Create |
| PUT | `/api/servers/:id` | 是 | 是 | handler/mcp_server.go → Update |
| DELETE | `/api/servers/:id` | 是 | 是 | handler/mcp_server.go → Delete |
| POST | `/api/servers/fetch-tools` | 是 | 是 | handler/mcp_server.go → FetchTools |
| POST | `/api/servers/publish` | 是 | 是 | handler/mcp_server.go → Publish |
| GET | `/api/roles` | 是 | 是 | handler/role.go → List |
| POST | `/api/roles` | 是 | 是 | handler/role.go → Create |
| PUT | `/api/roles/:id` | 是 | 是 | handler/role.go → Update |
| DELETE | `/api/roles/:id` | 是 | 是 | handler/role.go → Delete |
| GET | `/api/roles/:id/users` | 是 | 是 | handler/role.go → GetUsers |
| PUT | `/api/roles/:id/users` | 是 | 是 | handler/role.go → AssignUsers |
| GET | `/api/users` | 是 | 是 | handler/rbac_user.go → List |
| POST | `/api/users` | 是 | 是 | handler/rbac_user.go → Create |
| PUT | `/api/users/:id` | 是 | 是 | handler/rbac_user.go → Update |
| DELETE | `/api/users/:id` | 是 | 是 | handler/rbac_user.go → Delete |
| GET | `/api/settings` | 是 | 是 | handler/setting.go → Get |
| PUT | `/api/settings/:key` | 是 | 是 | handler/setting.go → Save |
| GET | `/api/settings/gateway-status` | 是 | - | handler/setting.go → GatewayStatus |

## 内嵌版本前端路由（重要）

`main_embed.go` 的 `NoRoute` 处理 SPA 回退逻辑：
- 非 API 路径且对应文件不存在时，回退到 `index.html`
- 静态资源路径（`/assets/...`、`/favicon.svg`）直接返回对应文件

## 前端导航说明

前端不使用 Vue Router，通过 `stores/nav.js` 的响应式 `active` 状态切换组件：
- `App.vue` 使用 `<component :is="...">` 动态渲染
- 导航项按角色定义（admin / user）
