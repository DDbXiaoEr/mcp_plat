# API 开发进度

> 每次实现或修改 API 后，及时更新此文件中的进度状态。

## 进度状态

| 模块 | 接口 | 状态 | 备注 |
|------|------|------|------|
| 认证 | POST /api/auth/login | ✅ 已完成 | |
| 认证 | GET /api/auth/profile | ✅ 已完成 | |
| AccessKey | GET /api/access-keys | ✅ 已完成 | |
| AccessKey | POST /api/access-keys | ✅ 已完成 | |
| AccessKey | PUT /api/access-keys/:id | ✅ 已完成 | |
| AccessKey | DELETE /api/access-keys/:id | ✅ 已完成 | |
| 使用历史 | GET /api/history | ⏳ 计划中 | 存储方案未确定 |
| MCP 服务器 | GET /api/servers | ✅ 已完成 | |
| MCP 服务器 | POST /api/servers | ✅ 已完成 | |
| MCP 服务器 | PUT /api/servers/:id | ✅ 已完成 | |
| MCP 服务器 | DELETE /api/servers/:id | ✅ 已完成 | |
| MCP 服务器 | POST /api/servers/fetch-tools | ✅ 已完成 | 前端已接入 |
| RBAC | GET /api/roles | ✅ 已完成 | |
| RBAC | POST /api/roles | ✅ 已完成 | |
| RBAC | PUT /api/roles/:id | ✅ 已完成 | |
| RBAC | DELETE /api/roles/:id | ✅ 已完成 | |
| RBAC | GET /api/roles/:id/users | ✅ 已完成 | |
| RBAC | PUT /api/roles/:id/users | ✅ 已完成 | |
| RBAC | GET /api/users | ✅ 已完成 | |
| RBAC | POST /api/users | ✅ 已完成 | |
| RBAC | PUT /api/users/:id | ✅ 已完成 | |
| RBAC | DELETE /api/users/:id | ✅ 已完成 | |

- ✅ 已完成
- 🚧 开发中
- ⏳ 待开发
- ⏳ 计划中
