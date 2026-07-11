# mcp_plat-console

MCP 服务平台 · 管理控制台

## 技术栈

- **后端**: Go + Gin + GORM + PostgreSQL + JWT
- **前端**: Vue 3 + Vite（`web/` 目录）

## 快速开始

### 1. 环境准备

```bash
cp .env.example .env
# 编辑 .env 填写数据库、JWT 等配置
```

### 2. 运行

```bash
# 纯后端模式
make run

# 开发模式（同时启动后端 + 前端 dev server）
make dev
```

### 3. 构建

| 命令 | 说明 |
|------|------|
| `make build-server` | 纯后端，不含前端 |
| `make build-web` | 仅构建前端 |
| `make build` | 构建前端 + 内嵌单二进制（shell 模式） |
| `make build-embed` | 同上，但使用 `embed` 构建标签 |

构建时自动注入 Git 版本信息：

```
mcp_plat-console version v1.0.0 (commit abc1234), built at 2026-07-11_06:42:21
```

### 4. 默认管理员账号

| 用户名 | 密码 |
|--------|------|
| admin | admin123 |

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DB_HOST` | 数据库地址 | localhost |
| `DB_PORT` | 数据库端口 | 5432 |
| `DB_USER` | 数据库用户 | root |
| `DB_PASSWORD` | 数据库密码 | 123456 |
| `DB_NAME` | 数据库名 | mcp_platform |
| `JWT_SECRET` | JWT 签名密钥 | - |
| `SERVER_PORT` | 服务端口 | 8080 |
| `ADMIN_USERNAME` | 管理员用户名 | admin |
| `ADMIN_PASSWORD` | 管理员密码 | admin123 |

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

详见 `public/api_doc.md`

## 项目结构

```
├── main.go             # 入口（纯后端，默认构建标签）
├── main_embed.go       # 入口（内嵌前端，embed 构建标签）
├── version.go          # 版本变量（Git 注入）
├── config/             # 配置加载
├── database/           # GORM 初始化
├── model/              # 数据模型
├── handler/            # HTTP 处理器
├── service/            # 业务逻辑
├── middleware/         # JWT 鉴权中间件
├── router/             # Gin 路由注册
├── web/                # Vue 3 前端
└── public/api_doc.md   # API 文档
```
