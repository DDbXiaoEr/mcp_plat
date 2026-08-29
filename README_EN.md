# mcp_plat-console

MCP Service Platform · Admin Console

> **Personal Note**: All current modifications to this project have been done by AI. The author suffered a cerebral hemorrhage last year resulting in hemiplegia, and currently only has the use of one hand, so many features have not yet been implemented. I go to the hospital for rehabilitation every day, with limited time and energy. Without a job or income, tokens rely entirely on free new-user quotas and invitation rewards from various platforms. Development efficiency has reached its limit. Those with abundant tokens are welcome to contribute. Thanks to GLM and Alibaba Cloud Bailian.

## Features

`mcp_plat-console` is the admin backend of an MCP (Model Context Protocol) service platform, designed for **platform administrators**. It provides full-lifecycle management of MCP services, user permission control, and API gateway integration.

- **Full-lifecycle management** — From MCP Server registration and tool discovery to release, covering the entire service workflow
- **Flexible authentication** — Supports local accounts, LDAP, and CAS login to fit different organizational infrastructure
- **Fine-grained permissions** — Role-Based Access Control (RBAC) with AccessKey-level permission scoping down to individual services and tools
- **Tool-level control** — AccessKeys can be scoped to specific tools of each MCP Server. Tool lists are discovered automatically via the `tools/list` protocol, the frontend provides checkbox-based tool permission control, and runtime tool-level authorization is enforced through a custom APISIX plugin
- **Native gateway integration** — Deeply integrated with the APISIX API gateway, one-click route publishing, and a custom Go plugin for runtime authorization
- **Single-binary deployment** — The frontend is embedded at compile time, so a single executable is all you need to deploy

## Audience

| Role | Responsibility |
|------|------|
| **Super Admin** | Manage MCP Servers, users, role permissions, and system settings |
| **Regular User** | Manage personal AccessKeys and view usage history |

## Tech Stack

- **Backend**: Go 1.25 + Gin + GORM + PostgreSQL/MySQL + JWT + gRPC
- **Audit log storage**: Relational DB (PostgreSQL/MySQL) or ClickHouse (optional, separate from the main DB)
- **Frontend**: Vue 3 + Vite (`web/` directory)

## Documentation

- [Deployment architecture](deployment-architecture.md) (personal draft, in Chinese, for reference only — suggestions welcome)

## Test MCP Services

The project bundles a set of mock campus MCP Servers for functional testing and demos, located in the `testmcp/` directory.

### Service List

| Service | Port | Description |
|------|------|------|
| card-service | 8081 | Campus card service (balance query, transactions, loss reporting/unreporting, recharge) |
| academic-service | 8082 | Academic information service (timetable, grades, GPA, exam schedule) |
| network-service | 8083 | Campus network service (account info, traffic query, password reset, account activation/suspension) |
| library-service | 8084 | Library service (current loans, loan history, renewal, book search) |

### How to Run

```bash
cd testmcp

# Build all services
make build-local

# Run a single service
make run-card
make run-academic
make run-network
make run-library

# Or run them all at once
make run-all
```

### Built-in Test Accounts

All services share the following test accounts (student ID is passed as the `student_id` parameter):

| Student ID | Name | Card Status | Network Status | Notes |
|------|------|-----------|-----------|----------|
| 2021010101 | Zhang San | Normal | Active | 3 books borrowed |
| 2021010102 | Li Si | Normal | Active | 2 books borrowed |
| 2021010103 | Wang Wu | Normal | Active | 3 books borrowed, highest GPA |
| 2021010104 | Zhao Liu | **Reported lost** | **Suspended** | Lowest balance, running out of traffic |
| 2021010105 | Chen Qi | Normal | Active | 1 book borrowed |

> See `testmcp/test_accounts.json` for detailed test data.

## Quick Start

### 1. Configuration

```bash
# Copy the sample config file, then edit config.yaml with your DB, JWT, and other settings
# The main DB supports both PostgreSQL and MySQL
cp config.yaml.example config.yaml
```

Key `config.yaml` settings (see `config.yaml.example` for the full sample):

```yaml
# Main DB (users, roles, AccessKeys, servers, system settings, etc.)
database:
  type: postgres           # or mysql
  postgres:
    host: localhost
    port: 5432
    user: root
    password: 123456
    dbname: mcp_platform

# Audit log DB (shared with audit-log-server), separate from the main DB
# type: postgres / mysql / clickhouse
audit_log_db:
  type: postgres
  postgres:
    host: localhost
    port: 5432
    user: root
    password: 123456
    dbname: mcp_history_db
  # To switch to ClickHouse (high volume):
  # type: clickhouse
  # clickhouse:
  #   addr: 127.0.0.1:9000
  #   user: default
  #   password: ""
  #   db: mcp_platform
  #   table: audit_logs
  #   engine: merge_tree            # standalone: merge_tree; cluster: replicated_merge_tree
  #   cluster: ""                   # required when engine is replicated_merge_tree
  #   ttl_days: 90                  # data retention in days

jwt_secret: "at-least-32-character-random-string"
access_key_secret: "another-at-least-32-character-random-string"
server_port: 8080

admin:
  username: admin
  password: admin123
```

### 2. Run

```bash
# Backend-only mode
make run

# Development mode (starts backend :8080 + frontend dev server :5174)
make dev
```

The frontend dev server automatically proxies `/api` requests to the backend on `:8080`.

### 3. Build

| Command | Description |
|------|------|
| `make build-server` | Backend-only binary (local debug) |
| `make build-web` | Build the frontend only (`web/dist/`) |
| `make build-embed` | Build frontend + embedded single binary |
| `make build-local` | Local full build: server + accesskey-auth-server + audit-log-server |
| `make build-local-release` | Local release build: embed + accesskey-auth-server + audit-log-server |
| `make build-linux-amd64-release` | Cross-compile linux amd64 (server / embed / gRPC services) |
| `make build-linux-arm64-release` | Cross-compile linux arm64 (server / embed / gRPC services) |
| `make build-accesskey-auth-server` | Build the AccessKey gRPC auth service |
| `make build-audit-log-server` | Build the audit log gRPC service |
| `make build-apisix-plugin` | Build the APISIX Go plugin runner (linux/$(ARCH)) |
| `make build-tools` | Build the tool set (datagen / accesskey-test / auditgen, etc.) |
| `make build-auditgen` | Build the audit data generator (`build/bin/tools/auditgen`) |
| `make auditgen` | Run the audit data generator (`DAYS=30 PER_DAY=200`) |

Git version info is injected automatically at build time:

```
mcp_plat-console version v1.0.0 (commit abc1234), built at 2026-07-11_06:42:21
```

## Docker Compose Deployment

The project ships with ready-to-use Docker Compose configs that spin up the full environment (databases, APISIX gateway, gRPC auth/audit services, OpenLDAP, etc.) with a single command.

### Prerequisites

- Docker Engine 20.10+
- Docker Compose v2 (`docker compose` command)

### 1. Build Images

```bash
# amd64 architecture
make docker-build

# arm64 architecture (Apple Silicon / Kunpeng, etc.)
make docker-build-arm64

# APISIX gateway image (built separately, installs the custom plugin)
make docker-build-mcp_plat_apisix

# Backend-only image (built separately, frontend not embedded; deploy/proxy the frontend separately)
make docker-build-mcp-plat-server          # linux/amd64
make docker-build-arm64-mcp-plat-server    # linux/arm64
```

The build produces the following images:

| Image | Description |
|------|------|
| `mcp_plat_embed:latest` | Main app (embedded frontend) |
| `mcp_plat_server:latest` | Main app (backend only, no embedded frontend) |
| `accesskey-auth-server:latest` | AccessKey gRPC auth service |
| `audit-log-server:latest` | Audit log gRPC service |

### 2. Start Services

Enter the `dockercompose/` directory and choose a mode based on your audit storage needs:

**PostgreSQL mode** (small/medium scale):

```bash
cd dockercompose
TAG=latest docker compose -f docker-compose-postgres.yml up -d
```

**ClickHouse mode** (high volume / high throughput):

```bash
cd dockercompose
TAG=latest docker compose -f docker-compose-clickhouse.yml up -d
```

> The `TAG` value corresponds to the Git commit hash used at `make docker-build` time; check `docker images` for the actual tag.

### 3. Access

| Service | Address | Description |
|------|------|------|
| Admin console | http://localhost:8080 | Web admin UI |
| APISIX gateway | http://localhost:80 | API gateway entry |
| Default account | admin / admin123 | Super admin |

### LDAP Test Accounts

Docker Compose automatically initializes OpenLDAP with pre-provisioned test accounts:

| Type | Account range | uid format (filter regex, used for auto role assignment) | Email | Password |
|------|----------|----------|------|------|
| Student | student01 ~ student20 | `^student\d{2}$` | student01@xxx.edu.cn | 123456 |
| Teacher | teacher01 ~ teacher20 | `^teacher\d{2}$` | teacher01@xxx.edu.cn | 123456 |
| Staff | staff01 ~ staff10 | `^staff\d{2}$` | staff01@xxx.edu.cn | 123456 |

> The username used at login depends on the LDAP `userFilter` configured in system settings. The default filter is `(uid=%s)`, meaning you log in with the uid (e.g. `student01`). Admins can customize the matching rule under **System Settings → Authentication → LDAP User Filter**.
>
> See `dockercompose/ldap/init-data.ldif` for the test account data.

### 4. Stop Services

```bash
cd dockercompose
docker compose -f docker-compose-postgres.yml down    # PostgreSQL mode
# or
docker compose -f docker-compose-clickhouse.yml down  # ClickHouse mode
```

### Service Ports

| Service | Port | Description |
|------|------|------|
| mcp_plat_embed | 8080 | Admin console |
| APISIX | 80, 443, 9180 | API gateway |
| mainpostgres | 5432 | Main database |
| auditpostgres | 5433 | Audit database (PG mode only) |
| ClickHouse | 8123, 9000 | Audit database (CH mode only) |
| accesskey-auth-server | 9090 | AccessKey auth service |
| audit-log-server | 9091 | Audit log service |
| OpenLDAP | 389 | LDAP directory service |
| Etcd | 2379 | APISIX config center |

### Default Admin Account

| Username | Password |
|--------|------|
| admin | admin123 |

## Features

### Authentication

- Local database login (bcrypt password hashing)
- LDAP directory authentication with attribute mapping (auto-syncs user info on login)
- CAS single sign-on (frontend auto-detects the ticket parameter in the URL)
- Login method switchable via system settings, with online testing of LDAP connection and mapping
- JWT Bearer Token auth, distinguishing regular users from admins

### Platform Overview

- Dashboard: platform users, MCP server count, tool count, daily AI call count
- AI call trend line chart (last 7 / 30 days) with manual and auto refresh

### MCP Service Management

- Create, edit, and delete MCP Servers
- Service info: name, gateway path, backend address (multiple IP:Port), department, protocol (SSE / Streamable HTTP / stdio), protocol version (2025-03-26 / 2025-06-18 / 2026-07-28)
- **Tool discovery**: connect to an MCP Server and automatically call the `tools/list` protocol to pull available tools, with dropdown or custom address selection; picks randomly among multiple addresses
- **One-click publish**: batch-select MCP Servers to register as APISIX gateway routes; multiple addresses publish as round-robin load-balanced upstreams, with optional AccessKey auth plugin and custom headers
- Compatible with MCP 2025-06-18 spec, supports Bearer Token auth

### AccessKey Management

- Users can create multiple API keys for calling MCP services (the limit is configurable in system settings)
- Each key can specify the MCP Servers and tools it is allowed to access (fine-grained permission scope)
- Supports enable/disable and expiration dates; a scheduled task auto-disables expired keys
- Works with the APISIX gateway plugin for runtime AccessKey validation

### RBAC Permissions

- **Role management**: create/edit/delete roles; roles are associated with accessible MCP Servers
- **User management**: admins can create, edit, and delete users and assign users to roles
- Batch role assignment via transfer boxes (with account search)

### System Settings

System configuration is managed in groups and edited on a single page:

| Group | Items |
|------|--------|
| `platform` | Platform name, Logo URL, site address |
| `api_gateway` | APISIX Admin API address/key, publish domain, gRPC auth address, AccessKey header |
| `auth` | Login method switch (local / ldap / cas), LDAP/CAS/OAuth connection params |
| `log` | syslog forwarding (TCP/UDP) |
| `smtp` | Mail service configuration |
| `user_ops` | User default role, AccessKey limit, expiry-check cron expression |
| `network_security` | Network security policy (IP whitelist) |
| `audit_log` | Audit log toggle and audit-log-server gRPC address |

### APISIX Gateway Integration

- **Go plugin runner** (`cmd/apisix-runner`): a custom plugin running on the APISIX side that performs request-level AccessKey validation
- **gRPC auth service** (`cmd/accesskey-auth-server`): a standalone AccessKey verification service used by the gateway plugin
- **gRPC audit service** (`cmd/audit-log-server`): a standalone audit log collection service that writes to the audit DB
- Routes and upstreams are auto-registered to the APISIX Admin API when publishing MCP services

### Kong Gateway Integration (experimental)

- Switch the gateway type to **Kong** in system settings to publish MCP services via the Kong Admin API
- Only creates/updates **Upstream, Service, and Route** — **no plugins are deployed**
- Therefore, AccessKey auth, audit logging, and path rewriting do not work in the Kong scenario; the gateway path uses the URI path configured on the server (address)
- Local Docker testing: Admin API defaults to `http://127.0.0.1:8001`; the Admin API Key can be left empty when RBAC is not enabled

### Audit Log

- Audit log storage abstraction layer (`auditstore/`) with relational (PostgreSQL/MySQL) and ClickHouse implementations
- **Separate databases**: the main business DB and the audit DB are configured independently (`audit_log_db`)
- ClickHouse supports MergeTree / ReplicatedMergeTree engines and TTL retention; relational DBs support periodic cleanup
- Audit logs are stored by AccessKey ID (plaintext keys are not stored), reducing data volume

### Usage History

- Records usage logs for every MCP tool call (user, AccessKey, MCP service, endpoint, status)
- Filter and paginate by AccessKey, MCP service, tool name, and time range, with page jump support

### Dual Build Modes

- **Backend-only** (`!embed` tag): compiles only the API service; the frontend is deployed separately
- **Single binary** (`embed` tag): the frontend SPA is embedded into the Go binary for one-file deployment

## API

| Method | Path | Auth | Description |
|--------|------|:---:|------|
| POST | `/api/auth/login` | - | Login |
| GET | `/api/auth/method` | - | Get auth method (local / ldap / cas) |
| GET | `/api/auth/platform` | - | Get public platform info (login page name/logo/background) |
| POST | `/api/auth/cas/validate` | - | CAS login validation |
| GET | `/api/auth/profile` | Bearer | Profile |
| GET | `/api/access-keys` | Bearer | List AccessKeys |
| POST | `/api/access-keys` | Bearer | Create AccessKey |
| PUT | `/api/access-keys/:id` | Bearer | Update AccessKey |
| DELETE | `/api/access-keys/:id` | Bearer | Delete AccessKey |
| GET | `/api/history` | Bearer | Usage history (audit logs) |
| GET | `/api/audit-logs` | Bearer | Audit log list |
| GET | `/api/servers` | Bearer | MCP Server list |
| POST | `/api/servers` | Admin | Create MCP Server |
| PUT | `/api/servers/:id` | Admin | Update MCP Server |
| DELETE | `/api/servers/:id` | Admin | Delete MCP Server |
| POST | `/api/servers/fetch-tools` | Admin | Pull tools/list |
| POST | `/api/servers/publish` | Admin | Batch publish to API gateway |
| GET | `/api/roles` | Admin | Role list |
| POST | `/api/roles` | Admin | Create role |
| PUT | `/api/roles/:id` | Admin | Update role |
| DELETE | `/api/roles/:id` | Admin | Delete role |
| GET | `/api/roles/:id/users` | Admin | Get users assigned to role |
| PUT | `/api/roles/:id/users` | Admin | Batch assign users to role |
| GET | `/api/users` | Admin | User list |
| POST | `/api/users` | Admin | Create user |
| PUT | `/api/users/:id` | Admin | Update user |
| DELETE | `/api/users/:id` | Admin | Delete user |
| GET | `/api/overview/stats` | Admin | Platform overview stats |
| GET | `/api/overview/call-trend` | Admin | AI call trend |
| GET | `/api/settings` | Admin | Get all system settings |
| PUT | `/api/settings/:key` | Admin | Save system setting |
| POST | `/api/settings/test-ldap` | Admin | Test LDAP connection and mapping |
| POST | `/api/settings/test-smtp` | Admin | Test SMTP mail sending |
| GET | `/api/settings/:key` | Bearer | Get a single setting |
| GET | `/api/settings/gateway-status` | Bearer | Check API gateway config status |

> See `public/api_doc.md` for detailed request/response fields and `public/api-devstatus.md` for development progress (both in Chinese)

## License

This project is licensed under the [GNU Affero General Public License v3.0](LICENSE) (AGPL-3.0).

Copyright (C) 2026 Zhaoquan Wang
