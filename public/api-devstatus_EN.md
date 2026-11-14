# API Development Status (English)

> After implementing or modifying any API, update the status below promptly.

## Progress Status

| Module | Endpoint | Status | Notes |
|--------|----------|--------|-------|
| Health check | GET /healthz | ✅ Done | liveness probe, no auth, for Docker / K8s |
| Health check | GET /readyz | ✅ Done | readiness probe, checks DB connectivity, no auth |
| Auth | POST /api/auth/login | ✅ Done | Cookie session (Redis); sets HttpOnly cookie on success. With `admin.enabled: true` the admin account can log in; with `false` admin login is disabled and only the platform-configured auth method is used |
| Auth | POST /api/auth/logout | ✅ Done | deletes Redis session and clears cookie |
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
| MCP server | POST /api/servers/publish | ✅ Done | new accesskey_header param; new publish_domain param selected from the operation-settings publish-domain list, written into the route host/hosts and recorded on the server's publish_domain; on success marks status=published and records auth_enabled (always false for Kong); supports provider=kong (upstream/Service/routes only, no plugins); accesskey_verify and audit_log plugins accept multiple backend gRPC addresses (grpc_addrs) with least-connection load balancing |
| MCP server | POST /api/servers/maintenance | ✅ Done | set/cancel maintenance; APISIX returns 503 via the mocking plugin (response_status=503, legacy fallback mock/response_code), Kong uses request-termination; leaving maintenance restores the original route saved at publish (GatewayRoute); entering sets status=maintenance, restoring sets published; uses each server's recorded gateway_id to locate the gateway (falling back to the default gateway), so one batch may span multiple gateways |
| MCP server | GET /api/servers | ✅ Done | added description, status, auth_enabled fields |
| MCP server | POST /api/servers | ✅ Done | added description field |
| MCP server | PUT /api/servers/:id | ✅ Done | added description field |
| MCP server | DELETE /api/servers/:id | ✅ Done | |
| MCP server | POST /api/servers/fetch-tools | ✅ Done | frontend connected; returns `[{name, description}]` |
| Gateway server | GET /api/gateway-servers | ✅ Done | lists gateway instances; probes each gateway Admin API concurrently (3s timeout, online/offline + latency); the default gateway comes first |
| Gateway server | POST /api/gateway-servers | ✅ Done | creates a gateway instance; the first one becomes default automatically; APISIX needs an Admin Key, Kong may leave it empty |
| Gateway server | PUT /api/gateway-servers/:id | ✅ Done | updates a gateway instance; an empty admin_key keeps the existing key |
| Gateway server | DELETE /api/gateway-servers/:id | ✅ Done | deletes a gateway instance; deleting the default promotes the oldest remaining gateway |
| Gateway server | POST /api/gateway-servers/:id/default | ✅ Done | sets the default gateway (publish/maintenance/quick access all use it) |
| Gateway server | POST /api/gateway-servers/:id/test | ✅ Done | tests the connectivity of a saved gateway (5s timeout) |
| Gateway server | POST /api/gateway-servers/check | ✅ Done | tests unsaved form values; pass id when editing to reuse the saved Admin Key |
| Gateway server | legacy config migration | ✅ Done | at startup the old settings.api_gateway is migrated into a default gateway record and the old setting is deleted |
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
| System settings | PUT /api/settings/:key | ✅ Done | key: log / smtp / auth / user_ops / platform / publish_domain / network_security / audit_log / quick_access; API gateway configuration moved out of settings to /api/gateway-servers |
| System settings | PUT /api/settings/user_ops | ✅ Done | new-user auto role assignment switched to a single global filter attribute filterAttribute + per-role regex rules (roleRules: roleId + pattern), matched in order by attribute at login; first match assigns, none match → unassigned until admin handles it; removed the old defaultRoleId |
| System settings | PUT /api/settings/log | ✅ Done | hot reload after save; with Syslog disabled logs go to stdout + local file (default ./logs/mcp_plat.log) with size/count/days rotation and gzip compression; both disabled when Syslog is enabled |
| System settings | GET /api/settings/:key | ✅ Done | login only |
| System settings | GET /api/settings/gateway-status | ✅ Done | default gateway status (from API Gateway Server Management); never leaks admin config; returns provider / adminUrl / accesskeyHeader; publish domains split out into a standalone publish_domain setting |
| System settings | PUT/GET /api/settings/publish_domain | ✅ Done | publish-domain setting (domains array) supporting multiple domains; selected when publishing MCP servers; readable by regular users |
| System settings | PUT/GET /api/settings/quick_access | ✅ Done | quick-access settings (client enable/disable + http/https access protocol); readable by regular users |
| System settings | PUT/GET /api/settings/platform | ✅ Done | platform settings add loginBackground (login-page background image); frosted-glass login card + adaptive background fill |
| System settings | POST /api/settings/test-ldap | ✅ Done | LDAP connection and attribute-mapping test |
| System settings | POST /api/settings/test-smtp | ✅ Done | SMTP mail-notification config test (supports none/ssl/starttls); frontend "send test email" button connected |
| Email notification | service/notify.go | ✅ Done | wraps SendNotificationMail / SendNotificationMailToMany / SendNotificationMailToUser; reads SMTP config from system settings and checks the enabled switch |
| Email notification | AccessKey-expiry disable notice | ✅ Done | after the DisableExpiredAccessKeys cron disables expired keys, emails users (aggregated per user); the body lists only key names + expiry times (HTML-escaped) and asks them to log in and delete; reuses the SMTP enabled switch |
| MCP publish (frontend) | — | ✅ Done | the publish dialog adds an "API gateway server" dropdown (default gateway preselected) and the confirm step shows the target gateway; the selection is sent as gateway_id and stored on the server |
| Gateway management (frontend) | — | ✅ Done | adds an "API Gateway Servers" page to the admin sidebar with card-based instances (online status/latency/default badge/auth addresses) and create, edit, delete, set-default and connectivity-test actions; the old gateway configuration was removed from system settings |
| Quick access (frontend) | — | ✅ Done | user sidebar adds a "Quick Access" page; generates AccessKey-carrying mcpServers config for the CherryStudio client and one-click import via the `cherrystudio://mcp/install?servers=` deep link; servers without key validation (auth_enabled=false) need no AccessKey; client enable/disable and the access protocol are controlled by admin ops settings; only published servers are shown, filtered by key permission; the endpoint uses the server's recorded publish_domain first and falls back to the first publish domain in operation settings |
| i18n (backend) | all endpoints | ✅ Done | added the `i18n` package, the unified `resp` response helper and the global `Locale` language middleware. `message` returns CN/EN by the `Accept-Language` request header (Chinese is the source; English is a dictionary-style approximate translation); the LDAP test `data.message` and overview-chart fallback labels are localized; notification email / gateway-maintenance copy stays Chinese (out of scope this round) |
| i18n (frontend) | entire site | ✅ Done | introduced vue-i18n (zh-CN/en-US; follows the browser by default and can be toggled, remembered in localStorage `mcp-console-locale`); all copy across 13 components + stores/api moved to keys; `<html lang>` and `document.title` follow the language; dates/numbers formatted by locale; the switcher is on the login page and the top-right of the header |

- ✅ Done
- 🚧 In progress
- ⏳ To be developed
- ⏳ Planned
