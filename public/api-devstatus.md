# API 开发进度

> 每次实现或修改 API 后，及时更新此文件中的进度状态。

## 进度状态

| 模块 | 接口 | 状态 | 备注 |
|------|------|------|------|
| 健康检查 | GET /healthz | ✅ 已完成 | 存活探针，无鉴权，供 Docker / K8s 使用 |
| 健康检查 | GET /readyz | ✅ 已完成 | 就绪探针，校验数据库连通性，无鉴权 |
| 认证 | POST /api/auth/login | ✅ 已完成 | Cookie Session（Redis），登录成功下发 HttpOnly Cookie；`admin.enabled: true` 时管理员账号可登录，`false` 时禁用管理员登录，仅使用平台配置的认证方式 |
| 认证 | POST /api/auth/logout | ✅ 已完成 | 删除 Redis 会话并清除 Cookie |
| 认证 | GET /api/auth/method | ✅ 已完成 | 前端登录页按返回方式渲染 |
| 认证 | GET /api/auth/platform | ✅ 已完成 | 登录页获取平台公开信息（名称/Logo/跳转链接/登录背景图 loginBackground），无鉴权 |
| 认证 | POST /api/auth/cas/validate | ✅ 已完成 | CAS ticket 验证 |
| 认证 | GET /api/auth/profile | ✅ 已完成 | |
| 认证 | PUT /api/auth/profile | ✅ 已完成 | 普通用户个人信息邮箱修改，仅限本人，格式校验 + 唯一性校验，管理员返回 403 |
| AccessKey | GET /api/access-keys | ✅ 已完成 | |
| AccessKey | POST /api/access-keys | ✅ 已完成 | |
| AccessKey | PUT /api/access-keys/:id | ✅ 已完成 | 过期（系统禁用）Key 禁止启用，返回"该 AccessKey 已过期，无法启用，请删除后重新创建"；手动禁用（未过期）可恢复 |
| AccessKey | DELETE /api/access-keys/:id | ✅ 已完成 | |
| 使用历史 | GET /api/history | ✅ 已完成 | 审计日志存储可配置（关系库 / ClickHouse），AccessKey 以 ID 存储减少数据量 |
| 使用历史 | GET /api/audit-logs | ✅ 已完成 | 审计日志列表查询 |
| 使用历史 | client_ip 记录 | ✅ 已完成 | APISIX audit_log 插件采集请求来源 IP（r.SrcIP）并随 gRPC LogAccessRequest 落库（audit_logs.client_ip），历史/审计列表返回 client_ip 字段；ClickHouse 建表含该列并对旧表幂等补列，关系库由 AutoMigrate 自动加列 |
| MCP 服务器 | POST /api/servers/publish | ✅ 已完成 | 新增 accesskey_header 参数；新增 publish_domain 参数，从运营设置「发布域名」列表中选择，写入路由 host/hosts，并记录到服务器的 publish_domain；发布成功后自动标记 status=published 并记录 auth_enabled（Kong 恒 false）；支持 provider=kong 时仅发布上游/Service/路由（不下发插件）；accesskey_verify 与 audit_log 插件支持下发多个后端 gRPC 地址（grpc_addrs），按连接数负载均衡；新增 gateway_id 参数，可将服务器发布到指定网关，并记录到服务器的 gateway_id |
| MCP 服务器 | POST /api/servers/maintenance | ✅ 已完成 | 设置/取消维护状态；APISIX 用 mocking 插件（response_status=503，旧版回退 mock/response_code）返回 503，Kong 用 request-termination 插件；取消维护恢复发布时保存的原始路由（GatewayRoute）；进入维护置 status=maintenance，取消恢复 published；按服务器记录的 gateway_id 定位网关（为空回退默认网关），支持同批服务器分属不同网关 |
| MCP 服务器 | GET /api/servers | ✅ 已完成 | 新增 description、status、auth_enabled 字段 |
| MCP 服务器 | POST /api/servers | ✅ 已完成 | 新增 description 字段 |
| MCP 服务器 | PUT /api/servers/:id | ✅ 已完成 | 新增 description 字段 |
| MCP 服务器 | DELETE /api/servers/:id | ✅ 已完成 | |
| MCP 服务器 | POST /api/servers/fetch-tools | ✅ 已完成 | 前端已接入，返回 `[{name, description}]` 格式 |
| 网关服务器 | GET /api/gateway-servers | ✅ 已完成 | 列表返回网关实例；并发探测各网关 Admin API 连通性（3 秒超时，online/offline + 耗时），默认网关排在最前 |
| 网关服务器 | POST /api/gateway-servers | ✅ 已完成 | 新增网关实例；首个网关自动设为默认；APISIX 需 Admin Key，Kong 可留空 |
| 网关服务器 | PUT /api/gateway-servers/:id | ✅ 已完成 | 更新网关实例；admin_key 留空表示保留原 Key |
| 网关服务器 | DELETE /api/gateway-servers/:id | ✅ 已完成 | 删除网关实例；删除默认网关时自动将最早的剩余网关设为默认 |
| 网关服务器 | POST /api/gateway-servers/:id/default | ✅ 已完成 | 设为默认网关（发布/维护/快速接入统一使用默认网关） |
| 网关服务器 | POST /api/gateway-servers/:id/test | ✅ 已完成 | 测试已保存网关的连通性（5 秒超时） |
| 网关服务器 | POST /api/gateway-servers/check | ✅ 已完成 | 测试未保存的表单配置；编辑时可传 id 复用已保存的 Admin Key |
| 网关服务器 | 旧配置迁移 | ✅ 已完成 | 服务启动时将旧 settings.api_gateway 迁移为默认网关记录并删除旧设置项 |
| RBAC | GET /api/roles | ✅ 已完成 | |
| RBAC | POST /api/roles | ✅ 已完成 | |
| RBAC | PUT /api/roles/:id | ✅ 已完成 | |
| RBAC | DELETE /api/roles/:id | ✅ 已完成 | |
| RBAC | GET /api/roles/:id/users | ✅ 已完成 | |
| RBAC | PUT /api/roles/:id/users | ✅ 已完成 | |
| RBAC | GET /api/users | ✅ 已完成 | 用户管理页新增搜索框，支持按学号/工号（uid）模糊搜索（?q=），可与角色筛选（?role_id=）叠加 |
| RBAC | POST /api/users | ✅ 已完成 | |
| RBAC | POST /api/users/batch | ✅ 已完成 | 批量创建用户：一次提交整批（含姓名/邮箱/学号工号/电话/单位），整批统一角色 role_id；库内冲突或列表内重复则整批回滚，返回具体冲突项 |
| RBAC | PUT /api/users/:id | ✅ 已完成 | |
| RBAC | DELETE /api/users/:id | ✅ 已完成 | |
| 平台概况 | GET /api/overview/stats | ✅ 已完成 | 统计卡片已接入（用户/服务器/工具/当天调用数） |
| 平台概况 | GET /api/overview/call-trend | ✅ 已完成 | AI 调用历史趋势已接入（近 7/30 天，audit_logs 按天分组统计）；含服务器调用量占比、用户组用户数占比饼图数据，及按天×服务器堆叠柱状图数据（Top 8 + 其他） |
| 系统设置 | GET /api/settings | ✅ 已完成 | 前端已接入 |
| 系统设置 | PUT /api/settings/:key | ✅ 已完成 | key: log / smtp / auth / user_ops / platform / publish_domain / network_security / audit_log / quick_access；API 网关配置已从设置项移除，改由 /api/gateway-servers 管理 |
| 系统设置 | PUT /api/settings/user_ops | ✅ 已完成 | 新用户角色自动分配改为全局单一过滤属性 filterAttribute + 按角色配置正则规则（roleRules：roleId + pattern），登录时按属性依次匹配，命中即分配，均不匹配则暂不分配等待管理员手动处理；移除原 defaultRoleId |
| 系统设置 | PUT /api/settings/log | ✅ 已完成 | 保存后热生效；未启用 Syslog 时日志输出到标准输出 + 本地文件（默认 ./logs/mcp_plat.log），支持大小/数量/天数轮转与 gzip 压缩；启用 Syslog 后二者失效 |
| 系统设置 | GET /api/settings/:key | ✅ 已完成 | 仅需登录 |
| 系统设置 | GET /api/settings/gateway-status | ✅ 已完成 | 默认网关状态检测（取自 API 网关服务器管理），不泄露管理配置；返回 provider / adminUrl / accesskeyHeader；发布域名已拆分为独立的 publish_domain 设置项 |
| 系统设置 | PUT/GET /api/settings/publish_domain | ✅ 已完成 | 发布域名设置（domains 数组），支持配置多个域名；发布 MCP 时从中选择，普通用户可读 |
| 系统设置 | PUT/GET /api/settings/quick_access | ✅ 已完成 | 快速接入设置（客户端启停 + http/https 接入协议），普通用户可读 |
| 系统设置 | PUT/GET /api/settings/platform | ✅ 已完成 | 平台设置新增 loginBackground（登录页背景图），登录页毛玻璃卡片 + 背景图自适应填充 |
| 系统设置 | POST /api/settings/test-ldap | ✅ 已完成 | LDAP 连接与属性映射测试 |
| 系统设置 | POST /api/settings/test-smtp | ✅ 已完成 | 邮件通知 SMTP 配置测试（支持 none/ssl/starttls），前端「发送测试邮件」按钮已接入 |
| 邮件通知 | service/notify.go | ✅ 已完成 | 封装 SendNotificationMail / SendNotificationMailToMany / SendNotificationMailToUser，SMTP 配置读取系统设置并校验 enabled 开关 |
| 邮件通知 | AccessKey 到期禁用通知 | ✅ 已完成 | 定时任务 DisableExpiredAccessKeys 禁用过期 Key 后，按用户聚合发邮件，正文仅列出 Key 名称 + 到期时间（HTML 转义），提示登录平台删除；复用 SMTP enabled 开关 |
| MCP 发布（前端） | — | ✅ 已完成 | 发布对话框新增「API 网关服务器」下拉选择（默认选中默认网关），确认页展示目标网关；所选网关随发布请求写入 gateway_id |
| 网关管理（前端） | — | ✅ 已完成 | 管理员侧边栏新增「API 网关服务器」页，卡片式展示网关实例（在线状态/延迟/默认标记/认证地址），支持新增、编辑、删除、设为默认与连通性测试；系统设置页已移除原 API 网关配置 |
| 快速接入（前端） | — | ✅ 已完成 | 用户侧边栏新增「快速接入」页，支持 CherryStudio 客户端生成带 AccessKey 的 mcpServers 配置，并一键通过 `cherrystudio://mcp/install?servers=` 深链打开应用导入；未开启 Key 验证的服务器（auth_enabled=false）无需选择 AccessKey；客户端启停与接入协议由管理员运营设置控制；仅展示已发布服务器并按 Key 权限过滤；接入地址优先使用服务器记录的 publish_domain，未记录时回退运营设置中的首个发布域名 |
| i18n（后端） | 全部接口 | ✅ 已完成 | 新增 `i18n` 包与 `resp` 统一响应助手 + `Locale` 语言中间件（全局注册）。`message` 按请求头 `Accept-Language` 返回中/英文（中文为原文，英文为词典式近似翻译）；LDAP 测试 `data.message` 与概况图表兜底标签随语言本地化；通知邮件/网关维护文案保持中文（本期不含） |
| i18n（前端） | 全站 | ✅ 已完成 | 引入 vue-i18n（zh-CN/en-US，默认跟随浏览器并可切换，localStorage `mcp-console-locale` 记忆）；14 个组件 + stores/api 文案全部词条化；`<html lang>`、`document.title` 随语言更新；日期/数字按当前 locale 格式化；语言切换入口位于登录页与顶栏右侧 |

- ✅ 已完成
- 🚧 开发中
- ⏳ 待开发
- ⏳ 计划中
