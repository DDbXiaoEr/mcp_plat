# Project Conventions

> **See `PROJECT_STRUCTURE.md` (Chinese) / `PROJECT_STRUCTURE_EN.md` (English) for the project layout quick reference** — full directory tree, API route table and build modes. Read them first at the start of each session to understand the whole codebase without re-exploring it.

## Tech Stack

- **Backend**: Go + Gin + GORM + PostgreSQL + JWT
- **Frontend**: Vue 3 + Vite (`web/` directory)

## Shared File Rules (mandatory)

- **`public/api_doc.md`** — request/response field documentation for the API; keep in sync whenever an endpoint is added or modified
- **`public/api-devstatus.md`** — API development status; update promptly after implementing or changing an endpoint
- **`PROJECT_STRUCTURE.md` / `PROJECT_STRUCTURE_EN.md`** — keep the directory tree and route table in sync (both the Chinese and English versions) whenever files are added or removed

## Backend Conventions

- Go code lives at the project root
- Layered architecture: `handler` → `service` → `model`
- Every HTTP endpoint returns JSON in the unified shape `{"code": 200, "message": "success", "data": {}}`; **must be emitted through the `resp` helpers (`resp.OK` / `resp.Fail`)**, never raw `c.JSON(gin.H{...})`, so that `message` is translated by `Accept-Language`
- The Chinese string in `message` is the translation-table key; when adding a new user-facing Chinese message, also add the corresponding English entry (a whole sentence or short phrase) to `i18n/zh_en.go` — English is an approximate translation
- Database access uses GORM; all models are defined under `model/`
- Auth uses JWT; middleware lives under `middleware/`
- Environment-variable template is `.env.example`

## Frontend Conventions

- See `web/AGENTS.md` (Chinese) / `web/AGENTS_EN.md` (English)
