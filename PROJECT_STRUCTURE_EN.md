# Project Structure (English)

```
mcp_plat-console/
├── main.go                  # entry (!embed tag, backend-only build)
├── main_embed.go            # entry (embed tag, single-binary build with embedded frontend)
├── version.go               # Version / BuildTime variables
├── Makefile                 # build scripts (build / build-server / build-embed / build-web / dev / clean)
├── go.mod / go.sum
├── config.yaml.example      # main server config sample (copy to config.yaml to use)
├── config-generator.html    # standalone config generator page (HTML+pure JS, open directly in a browser): wizard-driven — first pick the deployment platform (Kubernetes incl. manifests / host or docker-compose, app configs only) × frontend mode (embedded / Nginx hosting), then configure storage/secrets/APISIX/frontend/K8s step by step; output adapts live (config.yaml, accesskey/audit/apisix YAML, optional nginx.conf) plus step_N K8s manifests (Namespace / external services / ConfigMap / Service / Deployment; Nginx mode adds mcp-plat-web and mcp-plat-server becomes ClusterIP), with CN/EN toggle
├── AGENTS.md                # project conventions (Chinese)
├── AGENTS_EN.md             # project conventions (English)
├── PROJECT_STRUCTURE.md     # this file (Chinese) — project structure quick reference
├── PROJECT_STRUCTURE_EN.md  # this file (English) — project structure quick reference
├── README.md
│
├── cmd/
│   ├── accesskey-auth-server/
│   │   └── main.go           # AccessKey auth gRPC service
│   ├── accesskey-test/
│   │   └── main.go           # AccessKey test tool
│   ├── apisix-runner/
│   │   ├── main.go           # APISIX runner entry
│   │   └── plugins/
│   │       ├── accesskey_verify.go   # AccessKey verification plugin
│   │       └── apisix_route_example.md
│   ├── audit-log-server/
│   │   └── main.go           # audit-log gRPC service
│   ├── auditgen/
│   │   └── main.go           # audit data generator (generates ~N days of audit logs from existing data)
│   └── datagen/
│       └── main.go           # test data generator
│
├── config/
│   └── config.go             # config loading (DB / JWT / port / admin account)
│
├── database/
│   └── database.go           # GORM init + AutoMigrate (entry point for the AuditStore abstraction)
│
├── auditstore/               # audit-log storage abstraction (configurable: relational DB / ClickHouse)
│   ├── store.go              # Store interface + factory New()
│   ├── gorm_store.go         # relational implementation (Postgres/MySQL) + periodic cleanup
│   └── clickhouse_store.go   # ClickHouse implementation (MergeTree / ReplicatedMergeTree + TTL)
│
├── logging/
│   └── logging.go            # logging init (default stdout + local file ./logs/mcp_plat.log with rotation; both disabled when Syslog is enabled)
│
├── model/
│   ├── user.go               # User model
│   ├── access_key.go         # AccessKey model
│   ├── usage_history.go      # UsageHistory model
│   ├── mcp_server.go         # MCPServer model
│   ├── role.go               # Role model
│   ├── role_server.go        # RoleServer association model
│   ├── audit_log.go          # AuditLog model (AccessKey stored by ID, can be stored in ClickHouse)
│   └── setting.go            # Setting key-value model (system settings)
│
├── handler/
│   ├── auth.go               # POST /api/auth/login, GET /api/auth/profile
│   ├── access_key.go         # CRUD /api/access-keys
│   ├── health.go             # GET /healthz (liveness), GET /readyz (readiness, checks DB)
│   ├── history.go            # GET /api/history
│   ├── mcp_server.go         # CRUD /api/servers
│   ├── overview.go           # GET /api/overview/stats (platform stats, admin only), GET /api/overview/call-trend (AI call trend, admin only)
│   ├── role.go               # CRUD /api/roles (admin only)
│   ├── rbac_user.go          # CRUD /api/users (admin only)
│   └── setting.go            # GET/PUT /api/settings (admin only)
│
├── i18n/                     # lightweight message translation package (dictionary-style, no third-party deps)
│   ├── i18n.go               # language detection Lang + Translate (sentence / regex template / longest-phrase-first replacement)
│   ├── zh_en.go              # Chinese→English translation table (all user-facing messages from handler/middleware/service)
│   └── i18n_test.go
│
├── resp/                     # unified JSON response helper (single exit point for c.JSON; translates message per request language)
│   ├── resp.go               # OK / Fail / JSON / Locale
│   └── resp_test.go
│
├── service/
│   ├── auth.go               # login business logic + JWT generation
│   ├── access_key.go         # AccessKey business logic
│   ├── cas.go                # CAS auth business logic
│   ├── ldap.go               # LDAP auth business logic
│   ├── history.go            # usage-history business logic
│   ├── kong.go               # Kong gateway publishing (upstreams/Service/routes only; no plugins pushed)
│   ├── mail.go               # email sending (SMTP, supports none/ssl/starttls)
│   ├── mcp_server.go         # MCPServer business logic
│   ├── mcp_server_test.go    # FetchTools unit test (mock + optional real server)
│   ├── notify.go             # unified email notification wrapper (SendNotificationMail etc.; SMTP config from system settings)
│   ├── overview.go           # platform-overview stats business logic (users/servers/tools/daily calls)
│   ├── role.go               # Role business logic
│   ├── rbac_user.go          # RBACUser business logic
│   └── setting.go            # system-settings business logic (key-value JSON storage)
│
├── middleware/
│   ├── auth.go               # JWT Bearer auth middleware + AdminRequired
│   └── locale.go             # language middleware (parses Accept-Language → c.Set("locale"))
│
├── plugin/
│   ├── access_key.go         # AccessKey gRPC plugin registration
│   ├── accesskey_client.go   # AccessKey gRPC client
│   ├── accesskey_grpc.pb.go  # gRPC generated code
│   ├── accesskey.pb.go       # protobuf generated code
│   ├── accesskey.proto       # AccessKey proto definition
│   ├── auditlog_client.go    # AuditLog gRPC client
│   ├── auditlog_grpc.pb.go   # gRPC generated code
│   ├── auditlog.pb.go        # protobuf generated code
│   ├── auditlog.proto        # AuditLog proto definition
│   └── grpc_lb.go            # gRPC multi-backend load balancing (least connections)
│
├── router/
│   └── router.go             # Gin route registration (used by main.go only; main_embed.go registers its own)
│
├── public/
│   ├── api_doc.md            # API endpoint documentation (Chinese; mandatory to maintain)
│   ├── api_doc_EN.md         # API endpoint documentation (English; mandatory to maintain)
│   ├── api-devstatus.md      # API development status (Chinese; mandatory to maintain)
│   └── api-devstatus_EN.md   # API development status (English; mandatory to maintain)
│
├── kubernetes/               # Kubernetes deployment manifests (namespace / configmap / services)
│   ├── 00-namespace.yaml
│   ├── 01-external-services.yaml
│   ├── 10-configmaps.yaml
│   ├── 40-accesskey-auth-server.yaml
│   └── 50-audit-log-server.yaml
│
├── testmcp/                  # test MCP services (card / network / library / academic)
│   ├── cmd/
│   │   ├── card/main.go      # card service entry
│   │   ├── network/main.go   # network service entry
│   │   ├── library/main.go   # library service entry
│   │   └── academic/main.go  # academic service entry
│   ├── internal/
│   │   ├── card/service.go   # card service logic
│   │   ├── network/service.go # network service logic
│   │   ├── library/service.go # library service logic
│   │   └── academic/service.go # academic service logic
│   ├── configs/              # per-service config files
│   ├── bin/                  # build artifacts
│   ├── go.mod / go.sum
│   ├── Makefile
│   └── test_accounts.json    # test accounts
│
└── web/                      # Vue 3 frontend (SPA)
    ├── AGENTS.md             # frontend conventions (Chinese)
    ├── AGENTS_EN.md          # frontend conventions (English)
    ├── PROJECT_STRUCTURE.md  # frontend project structure (Chinese)
    ├── PROJECT_STRUCTURE_EN.md # frontend project structure (English)
    ├── Makefile              # frontend build scripts (install / build / dev / preview / clean)
    ├── index.html            # HTML entry
    ├── package.json
    ├── vite.config.js        # Vite config (dev proxy /api → localhost:8080)
    ├── dist/                 # build artifacts (for embed packaging)
    │   ├── index.html
    │   ├── favicon.svg
    │   └── assets/
    └── src/
        ├── main.js           # Vue entry (mounts vue-i18n)
        ├── App.vue           # root component (conditionally rendered by login state)
        ├── api.js            # HTTP API wrapper
        ├── i18n.js           # vue-i18n instance (zh-CN/en-US) + roleLabel helper
        ├── locales/          # i18n entries: zh-CN/ (source Chinese) + en-US/ (English), split by view domain, merged by index.js
        ├── styles/
        │   └── global.css    # global CSS variables
        ├── stores/
        │   ├── auth.js       # auth state (login/logout/role code)
        │   ├── locale.js     # language state (persisted in localStorage mcp-console-locale)
        │   ├── nav.js        # sidebar nav state (no Vue Router; menu stores labelKey)
        │   ├── settings.js   # platform settings shared state (name/logo/jump links)
        │   └── theme.js      # light/dark theme state (persisted in localStorage)
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

## Build Modes

| Command | Build tag | Entry file | Description |
|---------|-----------|------------|-------------|
| `make build-server` / `go build .` | `!embed` (default) | `main.go` | backend only, no frontend |
| `make build-embed` / `go build -tags embed .` | `embed` | `main_embed.go` | single binary embedding `web/dist/` |

## API Route Overview

| Method | Path | Auth | Admin | Handler |
|--------|------|:---:|:---:|---------|
| GET | `/healthz` | no | - | handler/health.go → Healthz |
| GET | `/readyz` | no | - | handler/health.go → Readyz |
| POST | `/api/auth/login` | no | - | handler/auth.go → Login |
| GET | `/api/auth/method` | no | - | handler/auth.go → GetAuthMethod |
| GET | `/api/auth/platform` | no | - | handler/auth.go → GetPlatform |
| POST | `/api/auth/cas/validate` | no | - | handler/auth.go → CASValidate |
| GET | `/api/auth/profile` | yes | - | handler/auth.go → Profile |
| PUT | `/api/auth/profile` | yes | - | handler/auth.go → UpdateProfile |
| GET | `/api/access-keys` | yes | - | handler/access_key.go → List |
| POST | `/api/access-keys` | yes | - | handler/access_key.go → Create |
| PUT | `/api/access-keys/:id` | yes | - | handler/access_key.go → Update |
| DELETE | `/api/access-keys/:id` | yes | - | handler/access_key.go → Delete |
| GET | `/api/history` | yes | - | handler/history.go → List (⏳ planned) |
| GET | `/api/overview/stats` | yes | yes | handler/overview.go → Stats |
| GET | `/api/overview/call-trend` | yes | yes | handler/overview.go → CallTrend |
| GET | `/api/servers` | yes | - | handler/mcp_server.go → List |
| POST | `/api/servers` | yes | yes | handler/mcp_server.go → Create |
| PUT | `/api/servers/:id` | yes | yes | handler/mcp_server.go → Update |
| DELETE | `/api/servers/:id` | yes | yes | handler/mcp_server.go → Delete |
| POST | `/api/servers/fetch-tools` | yes | yes | handler/mcp_server.go → FetchTools |
| POST | `/api/servers/publish` | yes | yes | handler/mcp_server.go → Publish |
| POST | `/api/servers/maintenance` | yes | yes | handler/mcp_server.go → Maintenance |
| GET | `/api/roles` | yes | yes | handler/role.go → List |
| POST | `/api/roles` | yes | yes | handler/role.go → Create |
| PUT | `/api/roles/:id` | yes | yes | handler/role.go → Update |
| DELETE | `/api/roles/:id` | yes | yes | handler/role.go → Delete |
| GET | `/api/roles/:id/users` | yes | yes | handler/role.go → GetUsers |
| PUT | `/api/roles/:id/users` | yes | yes | handler/role.go → AssignUsers |
| GET | `/api/users` | yes | yes | handler/rbac_user.go → List |
| POST | `/api/users` | yes | yes | handler/rbac_user.go → Create |
| POST | `/api/users/batch` | yes | yes | handler/rbac_user.go → BatchCreate |
| PUT | `/api/users/:id` | yes | yes | handler/rbac_user.go → Update |
| DELETE | `/api/users/:id` | yes | yes | handler/rbac_user.go → Delete |
| GET | `/api/settings` | yes | yes | handler/setting.go → Get |
| PUT | `/api/settings/:key` | yes | yes | handler/setting.go → Save |
| GET | `/api/settings/:key` | yes | - | handler/setting.go → GetByKey |
| GET | `/api/settings/gateway-status` | yes | - | handler/setting.go → GatewayStatus |
| POST | `/api/settings/test-smtp` | yes | yes | handler/setting.go → TestSmtp |

> **i18n**: every endpoint's `message` is bilingual (CN/EN), selected by the `Accept-Language` request header (`en*` → English approximate translation, otherwise → original Chinese). Implemented by `middleware.Locale()` registered globally in `router.Setup()` plus the `resp` helper; the entry table lives in `i18n/`.

## Embedded Frontend Routing (important)

`main_embed.go`'s `NoRoute` handles SPA fallback:
- For a non-API path where the file does not exist, fall back to `index.html`
- Static asset paths (`/assets/...`, `/favicon.svg`) return the corresponding file directly

## Frontend Navigation Notes

The frontend does not use Vue Router; it switches components through the reactive `active` state in `stores/nav.js`:
- `App.vue` renders dynamically with `<component :is="...">`
- Nav items are defined per role (admin / user)
- The user menu includes a "Quick Access" page (QuickAccessView) that generates access configs carrying an AccessKey for common MCP clients (e.g. CherryStudio); the endpoint URL is assembled in the frontend from gateway settings (`defaultPublishDomain` + provider)
