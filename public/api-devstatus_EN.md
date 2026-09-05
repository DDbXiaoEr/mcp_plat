# API Development Status (English)

> After implementing or modifying any API, update the status below promptly.

## Progress Status

| Module | Endpoint | Status | Notes |
|--------|----------|--------|-------|
| Health check | GET /healthz | ✅ Done | liveness probe, no auth, for Docker / K8s |
| Health check | GET /readyz | ✅ Done | readiness probe, checks DB connectivity, no auth |
| Auth | POST /api/auth/login | ✅ Done | |
| Auth | GET /api/auth/method | ✅ Done | frontend login page renders by the returned method |
| Auth | GET /api/auth/platform | ✅ Done | login page fetches public platform info (name / logo / jump link / login background image loginBackground), no auth |
| Auth | POST /api/auth/cas/validate | ✅ Done | CAS ticket validation |
| Auth | GET /api/auth/profile | ✅ Done | |
| Auth | PUT /api/auth/profile | ✅ Done | regular users update their own profile email; format + uniqueness validation; admins get 403 |
| AccessKey | GET /api/access-keys | ✅ Done | |
| AccessKey | POST /api/access-keys | ✅ Done | |
| AccessKey | PUT /api/access-keys/:id | ✅ Done | expired (system-disabled) keys cannot be enabled, returns "该 AccessKey 已过期，无法启用，请删除后重新创建"; manually disabled (not expired) keys can be re-enabled |
| AccessKey | DELETE /api/access-keys/:id | ✅ Done | |
| Usage history | GET /api/history | ✅ Done | audit-log storage is configurable (relational DB / ClickHouse); AccessKey stored by ID to reduce volume |
| Usage history | GET /api/audit-logs | ✅ Done | audit-log list query |
| Usage history | client_ip recording | ✅ Done | the APISIX audit_log plugin collects the source IP (r.SrcIP) and sends it to audit_logs.client_ip via the gRPC LogAccessRequest; history/audit lists return client_ip; the ClickHouse table includes this column and backfills it idempotently on legacy tables, AutoMigrate adds the column on relational DBs |
| MCP server | POST /api/servers/publish | ✅ Done | new accesskey_header param; on success marks status=published and records auth_enabled (always false for Kong); supports provider=kong (upstream/Service/routes only, no plugins); accesskey_verify and audit_log plugins accept multiple backend gRPC addresses (grpc_addrs) with least-connection load balancing |
| MCP server | POST /api/servers/maintenance | ✅ Done | set/cancel maintenance; APISIX returns 503 via the mocking plugin (response_status=503, legacy fallback mock/response_code), Kong uses request-termination; leaving maintenance restores the original route saved at publish (GatewayRoute); entering sets status=maintenance, restoring sets published |
| MCP server | GET /api/servers | ✅ Done | added description, status, auth_enabled fields |
| MCP server | POST /api/servers | ✅ Done | added description field |
| MCP server | PUT /api/servers/:id | ✅ Done | added description field |
| MCP server | DELETE /api/servers/:id | ✅ Done | |
| MCP server | POST /api/servers/fetch-tools | ✅ Done | frontend connected; returns `[{name, description}]` |
| RBAC | GET /api/roles | ✅ Done | |
| RBAC | POST /api/roles | ✅ Done | |
| RBAC | PUT /api/roles/:id | ✅ Done | |
| RBAC | DELETE /api/roles/:id | ✅ Done | |
| RBAC | GET /api/roles/:id/users | ✅ Done | |
| RBAC | PUT /api/roles/:id/users | ✅ Done | |
| RBAC | GET /api/users | ✅ Done | user-management page adds a search box: fuzzy search by student/staff number uid (`?q=`), combinable with role filter (`?role_id=`) |
| RBAC | POST /api/users | ✅ Done | |
| RBAC | POST /api/users/batch | ✅ Done | batch user creation: submits a whole batch at once (name/email/student-staff-number/phone/unit), one uniform role_id; DB conflicts or in-list duplicates roll back the whole batch and return the specific conflicts |
| RBAC | PUT /api/users/:id | ✅ Done | |
| RBAC | DELETE /api/users/:id | ✅ Done | |
| Platform overview | GET /api/overview/stats | ✅ Done | stat cards connected (users/servers/tools/daily calls) |
| Platform overview | GET /api/overview/call-trend | ✅ Done | AI call-trend chart connected (last 7/30 days, audit_logs grouped by day); includes per-server call share and per-user-group pie data, plus per-day × per-server stacked-bar data (Top 8 + Other) |
| System settings | GET /api/settings | ✅ Done | frontend connected |
| System settings | PUT /api/settings/:key | ✅ Done | key: log / smtp / auth / user_ops / platform / api_gateway / network_security / audit_log |
| System settings | PUT /api/settings/user_ops | ✅ Done | new-user auto role assignment switched to a single global filter attribute filterAttribute + per-role regex rules (roleRules: roleId + pattern), matched in order by attribute at login; first match assigns, none match → unassigned until admin handles it; removed the old defaultRoleId |
| System settings | PUT /api/settings/log | ✅ Done | hot reload after save; with Syslog disabled logs go to stdout + local file (default ./logs/mcp_plat.log) with size/count/days rotation and gzip compression; both disabled when Syslog is enabled |
| System settings | GET /api/settings/:key | ✅ Done | login only |
| System settings | GET /api/settings/gateway-status | ✅ Done | gateway config status check; never leaks admin config; returns provider / adminUrl / defaultPublishDomain / accesskeyHeader |
| System settings | PUT/GET /api/settings/quick_access | ✅ Done | quick-access settings (client enable/disable + http/https access protocol); readable by regular users |
| System settings | PUT/GET /api/settings/platform | ✅ Done | platform settings add loginBackground (login-page background image); frosted-glass login card + adaptive background fill |
| System settings | POST /api/settings/test-ldap | ✅ Done | LDAP connection and attribute-mapping test |
| System settings | POST /api/settings/test-smtp | ✅ Done | SMTP mail-notification config test (supports none/ssl/starttls); frontend "send test email" button connected |
| Email notification | service/notify.go | ✅ Done | wraps SendNotificationMail / SendNotificationMailToMany / SendNotificationMailToUser; reads SMTP config from system settings and checks the enabled switch |
| Email notification | AccessKey-expiry disable notice | ✅ Done | after the DisableExpiredAccessKeys cron disables expired keys, emails users (aggregated per user); the body lists only key names + expiry times (HTML-escaped) and asks them to log in and delete; reuses the SMTP enabled switch |
| Quick access (frontend) | — | ✅ Done | user sidebar adds a "Quick Access" page; generates AccessKey-carrying mcpServers config for the CherryStudio client and one-click import via the `cherrystudio://mcp/install?servers=` deep link; servers without key validation (auth_enabled=false) need no AccessKey; client enable/disable and the access protocol are controlled by admin ops settings; only published servers are shown, filtered by key permission |
| i18n (backend) | all endpoints | ✅ Done | added the `i18n` package, the unified `resp` response helper and the global `Locale` language middleware. `message` returns CN/EN by the `Accept-Language` request header (Chinese is the source; English is a dictionary-style approximate translation); the LDAP test `data.message` and overview-chart fallback labels are localized; notification email / gateway-maintenance copy stays Chinese (out of scope this round) |
| i18n (frontend) | entire site | ✅ Done | introduced vue-i18n (zh-CN/en-US; follows the browser by default and can be toggled, remembered in localStorage `mcp-console-locale`); all copy across 13 components + stores/api moved to keys; `<html lang>` and `document.title` follow the language; dates/numbers formatted by locale; the switcher is on the login page and the top-right of the header |

- ✅ Done
- 🚧 In progress
- ⏳ To be developed
- ⏳ Planned
