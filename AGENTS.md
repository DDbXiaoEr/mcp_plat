# 项目规范

> **项目结构速查见 `PROJECT_STRUCTURE.md`**，包含完整目录树、API 路由表、构建模式说明。每次会话开始时优先读取该文件了解项目全貌，无需重新探索整个代码库。

## 技术栈

- **后端**: Go + Gin + GORM + PostgreSQL + JWT
- **前端**: Vue 3 + Vite（`web/` 目录）

## API 文档规范（重要）

**每个 API 请求参数的结构必须写入 `public/api_doc.md` 文件中。**

- 新增接口时，必须在 `public/api_doc.md` 中添加对应的接口文档
- 修改接口参数（增删改字段、修改类型等）时，必须同步更新 `public/api_doc.md`
- 文档格式参照已有的文档结构：请求路径、请求方法、请求参数（含类型、必填、说明）、响应示例
- 此为强制要求，不得遗漏

## 后端开发规范

- Go 代码放在项目根目录
- 分层架构：`handler` → `service` → `model`
- 所有 HTTP 接口统一返回 JSON 格式：`{"code": 200, "message": "success", "data": {}}`
- 数据库操作使用 GORM，所有模型定义在 `model/` 目录
- 鉴权使用 JWT，中间件在 `middleware/` 目录
- 环境变量配置模板在 `.env.example`

## 前端开发规范

- 见 `web/AGENTS.md`
