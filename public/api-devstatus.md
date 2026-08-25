# API 开发进度

> 每次实现或修改 API 后，及时更新此文件中的进度状态。

## 进度状态

| 模块 | 接口 | 状态 | 备注 |
|------|------|------|------|
| 健康检查 | GET /healthz | ✅ 已完成 | 存活探针，无鉴权，供 Docker / K8s 使用 |
| 健康检查 | GET /readyz | ✅ 已完成 | 就绪探针，校验数据库连通性，无鉴权 |
| 认证 | POST /api/auth/login | ✅ 已完成 | |
| 认证 | GET /api/auth/method | ✅ 已完成 | 前端登录页按返回方式渲染 |
| 认证 | POST /api/auth/cas/validate | ✅ 已完成 | CAS ticket 验证 |
| 认证 | GET /api/auth/profile | ✅ 已完成 | |
| AccessKey | GET /api/access-keys | ✅ 已完成 | |
| AccessKey | POST /api/access-keys | ✅ 已完成 | |
| AccessKey | PUT /api/access-keys/:id | ✅ 已完成 | |
| AccessKey | DELETE /api/access-keys/:id | ✅ 已完成 | |
| 使用历史 | GET /api/history | ✅ 已完成 | 审计日志存储可配置（关系库 / ClickHouse），AccessKey 以 ID 存储减少数据量 |
| 使用历史 | GET /api/audit-logs | ✅ 已完成 | 审计日志列表查询 |
| MCP 服务器 | POST /api/servers/publish | ✅ 已完成 | 新增 accesskey_header 参数；发布成功后自动标记 status=published 并记录 auth_enabled（Kong 恒 false）；支持 provider=kong 时仅发布上游/Service/路由（不下发插件）；accesskey_verify 与 audit_log 插件支持下发多个后端 gRPC 地址（grpc_addrs），按连接数负载均衡 |
| MCP 服务器 | POST /api/servers/maintenance | ✅ 已完成 | 设置/取消维护状态；APISIX 用 mocking 插件（response_status=503，旧版回退 mock/response_code）返回 503，Kong 用 request-termination 插件；取消维护恢复发布时保存的原始路由（GatewayRoute）；进入维护置 status=maintenance，取消恢复 published |
| MCP 服务器 | GET /api/servers | ✅ 已完成 | 新增 description、status、auth_enabled 字段 |
| MCP 服务器 | POST /api/servers | ✅ 已完成 | 新增 description 字段 |
| MCP 服务器 | PUT /api/servers/:id | ✅ 已完成 | 新增 description 字段 |
| MCP 服务器 | DELETE /api/servers/:id | ✅ 已完成 | |
| MCP 服务器 | POST /api/servers/fetch-tools | ✅ 已完成 | 前端已接入，返回 `[{name, description}]` 格式 |
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
| 平台概况 | GET /api/overview/stats | ✅ 已完成 | 统计卡片已接入（用户/服务器/工具/当天调用数） |
| 平台概况 | GET /api/overview/call-trend | ✅ 已完成 | AI 调用历史趋势已接入（近 7/30 天，audit_logs 按天分组统计）；含服务器调用量占比、用户组用户数占比饼图数据，及按天×服务器堆叠柱状图数据（Top 8 + 其他） |
| 系统设置 | GET /api/settings | ✅ 已完成 | 前端已接入 |
| 系统设置 | PUT /api/settings/:key | ✅ 已完成 | key: log / smtp / auth / user_ops / platform / api_gateway / network_security / audit_log |
| 系统设置 | PUT /api/settings/user_ops | ✅ 已完成 | 新用户角色自动分配改为全局单一过滤属性 filterAttribute + 按角色配置正则规则（roleRules：roleId + pattern），登录时按属性依次匹配，命中即分配，均不匹配则暂不分配等待管理员手动处理；移除原 defaultRoleId |
| 系统设置 | PUT /api/settings/log | ✅ 已完成 | 保存后热生效；未启用 Syslog 时日志输出到标准输出 + 本地文件（默认 ./logs/mcp_plat.log），支持大小/数量/天数轮转与 gzip 压缩；启用 Syslog 后二者失效 |
| 系统设置 | GET /api/settings/:key | ✅ 已完成 | 仅需登录 |
| 系统设置 | GET /api/settings/gateway-status | ✅ 已完成 | 网关配置状态检测，不泄露管理配置；返回 provider / adminUrl / defaultPublishDomain / accesskeyHeader |
| 系统设置 | PUT/GET /api/settings/quick_access | ✅ 已完成 | 快速接入设置（客户端启停 + http/https 接入协议），普通用户可读 |
| 系统设置 | POST /api/settings/test-ldap | ✅ 已完成 | LDAP 连接与属性映射测试 |
| 快速接入（前端） | — | ✅ 已完成 | 用户侧边栏新增「快速接入」页，支持 CherryStudio 客户端生成带 AccessKey 的 mcpServers 配置，并一键通过 `cherrystudio://mcp/install?servers=` 深链打开应用导入；未开启 Key 验证的服务器（auth_enabled=false）无需选择 AccessKey；客户端启停与接入协议由管理员运营设置控制；仅展示已发布服务器并按 Key 权限过滤 |

- ✅ 已完成
- 🚧 开发中
- ⏳ 待开发
- ⏳ 计划中
