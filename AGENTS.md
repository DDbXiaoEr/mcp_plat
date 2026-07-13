# 项目规范

> **项目结构速查见 `PROJECT_STRUCTURE.md`**，包含完整目录树、API 路由表、构建模式说明。每次会话开始时优先读取该文件了解项目全貌，无需重新探索整个代码库。

## 技术栈

- **后端**: Go + Gin + GORM + PostgreSQL + JWT
- **前端**: Vue 3 + Vite（`web/` 目录）

## 公共文件规范（强制）

- **`public/api_doc.md`** — 保存 API 的请求/响应字段说明，新增或修改接口时必须同步更新
- **`public/api-devstatus.md`** — 保存 API 开发进度状态，实现或修改接口后必须及时更新

## 后端开发规范

- Go 代码放在项目根目录
- 分层架构：`handler` → `service` → `model`
- 所有 HTTP 接口统一返回 JSON 格式：`{"code": 200, "message": "success", "data": {}}`
- 数据库操作使用 GORM，所有模型定义在 `model/` 目录
- 鉴权使用 JWT，中间件在 `middleware/` 目录
- 环境变量配置模板在 `.env.example`

## 前端开发规范

- 见 `web/AGENTS.md`
