# Project Rules: Project Tracker (`projex`)

## Project Overview
Full-stack project and task management application.
- **Backend**: Go 1.23, `chi` v5 router, PostgreSQL via sqlc-generated queries
- **Frontend**: Next.js 14 (React 18), plain `fetch`-based API client
- **Database**: PostgreSQL 16, schema managed by `goose`, queries by `sqlc`
- **Orchestration**: Docker Compose with two isolated bridge networks

## Repository Structure
```
projex/
├── server/              # Go REST API
│   ├── cmd/api/         # Entrypoint (main.go)
│   └── internal/
│       ├── config/      # Env var loading
│       ├── domain/      # Pure domain types (Project, Task)
│       ├── repositories/# Repository interfaces + stubs
│       ├── services/    # Service interfaces + implementations
│       ├── http/        # chi router setup
│       │   └── handlers/# HTTP handlers + writeJSON helper
│       └── dbgen/       # sqlc-generated code (not committed; must be generated)
├── ui/                  # Next.js frontend
│   ├── app/             # Next.js App Router pages and layout
│   └── lib/api.js       # Fetch-based API client
├── db/
│   ├── migrations/      # goose migration files
│   ├── queries/         # Raw SQL consumed by sqlc
│   └── sqlc.yaml        # sqlc config (outputs to server/internal/dbgen)
└── docker-compose.yml
```

## Current Implementation State
- All route wiring, handler, service, and repository interfaces are scaffolded and compile.
- **Repository implementations are stubs** — every method returns `ErrNotImplemented`. The real implementations backed by sqlc + PostgreSQL have not been written yet.
- `server/internal/dbgen/` does not exist yet. It must be generated via `sqlc generate -f ./db/sqlc.yaml` before any database-backed code can be written.
- The UI project detail page (`/projects/[id]`) is not yet created.

## Architecture Conventions

### Backend Layering (strictly enforced)
Requests must always flow through layers in this order — never skip layers:
```
HTTP Handler → Service → Repository → (sqlc-generated DB methods)
```
- **Handlers** (`internal/http/handlers/`): Parse HTTP input, call service, write JSON response. No business logic, no direct DB access.
- **Services** (`internal/services/`): Orchestrate business logic. Call repositories only through their interfaces.
- **Repositories** (`internal/repositories/`): All DB access. Must use sqlc-generated methods from `internal/dbgen` — no ad hoc SQL strings in services or handlers.
- **Domain** (`internal/domain/`): Plain Go structs shared across layers. No methods or logic.

### Interface-first design
- Repositories and services are always consumed via their Go interfaces (`ProjectRepository`, `TaskRepository`, `ProjectService`, `TaskService`).
- Concrete implementations are package-private structs (`projectService`, `taskService`, etc.) exposed only via `New*` constructor functions.
- This enables the existing stubs to be swapped for real implementations without touching callers.

### Dependency injection
All dependencies are wired manually in `cmd/api/main.go`. Do not use a DI framework.

## Database Rules
- Schema changes go in a new `goose` migration file under `db/migrations/`. Never modify existing migration files.
- New queries go in the appropriate file under `db/queries/` with sqlc annotations, then regenerate `dbgen`.
- After changing any `.sql` query file, always re-run: `sqlc generate -f ./db/sqlc.yaml`
- The sqlc config outputs to `server/internal/dbgen` with package name `dbgen`, using `database/sql`, with JSON tags and interface emission enabled.

## API Conventions
- All responses are JSON (`Content-Type: application/json`), written via the shared `writeJSON(w, status, payload)` helper in `handlers/response.go`.
- Successful creates return HTTP 201; successful deletes return HTTP 204 (no body).
- Decode errors return HTTP 400 with `{"error": "invalid request payload"}`.
- Route parameters are extracted with `chi.URLParam(r, "id")`.

## Configuration
All runtime config comes from environment variables with defaults in `internal/config/config.go`:
| Env Var | Default | Description |
|---|---|---|
| `PORT` | `8080` | API listen port |
| `DATABASE_URL` | `postgres://postgres:postgres@postgres:5432/project_tracker?sslmode=disable` | PostgreSQL DSN |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:3000` | Comma-separated allowed origins |

Never hardcode secrets or host-specific values in source files. Use `.env` files locally (see `.env.example`).

## Frontend Conventions
- API calls go through `ui/lib/api.js` only. Components must not call `fetch` directly.
- The API base URL is read from `process.env.NEXT_PUBLIC_API_BASE_URL` (default: `http://localhost:8080`).
- Use the Next.js App Router (`app/` directory). Pages that need client state use `"use client"`.

## Running the Project

### With Docker Compose (recommended)
```bash
docker compose up --build
```
- UI: http://localhost:3000
- API: http://localhost:8080
- DB: internal only

### Local development
```bash
# Server
cd server && go run ./cmd/api

# UI
cd ui && npm install && npm run dev
```

### Database migrations
```bash
goose -dir ./db/migrations postgres "postgres://postgres:postgres@localhost:5432/project_tracker?sslmode=disable" up
```

### Regenerate sqlc code
```bash
sqlc generate -f ./db/sqlc.yaml
```

## Data Model Summary
- **`projects`**: `id (bigserial PK)`, `name`, `description`, `created_at`, `updated_at`
- **`tasks`**: `id (bigserial PK)`, `project_id (FK → projects.id ON DELETE CASCADE)`, `title`, `description`, `status` (`todo`|`in_progress`|`done`), `priority (int)`, `created_at`, `updated_at`
- Index on `tasks(project_id)`.

## Out of Scope for v1
- Authentication / authorization
- Background jobs or queues
- WebSockets
- Advanced frontend state management (no Redux, Zustand, etc.)
