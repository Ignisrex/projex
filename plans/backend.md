# Backend Tasks

## Current State

The Go REST API skeleton is fully scaffolded and compiles cleanly. All routes, HTTP handlers, service interfaces, and repository interfaces are wired. However, **all repository methods return `ErrNotImplemented`** — no real database access exists yet. The `server/internal/dbgen/` directory (sqlc-generated code) has not been generated and must be created before any database-backed work can begin.

## Tasks

### 1. Generate sqlc Code

Run sqlc to produce typed Go database access methods from the SQL queries in `db/queries/`.

```bash
sqlc generate -f ./db/sqlc.yaml
```

This will create `server/internal/dbgen/` containing Go models and query methods for both `projects` and `tasks`. Do not edit the generated files — regenerate them whenever `db/queries/` changes.

### 2. Implement the Project Repository

Create `server/internal/repositories/project_repository.go` as a package-private struct with a `New*` constructor, following the same pattern as services and stubs in this codebase.

- Store a `*dbgen.Queries` field initialised once in the constructor via `dbgen.New(db)` — do not call `dbgen.New` inside individual methods.
- `ListProjects` → `q.ListProjects(ctx)`
- `GetProject` → `q.GetProject(ctx, id)` — the repository interface takes `id string` but sqlc expects `int64` (bigserial PK); parse with `strconv.ParseInt` and return `ErrNotFound` on `sql.ErrNoRows`
- `CreateProject` → `q.CreateProject(ctx, params)` — map `domain.Project` fields to the sqlc params struct
- `UpdateProject` → `q.UpdateProject(ctx, params)` — parse string ID to int64; return `ErrNotFound` when no row is updated
- `DeleteProject` → `q.DeleteProject(ctx, id)` — parse string ID to int64; return `ErrNotFound` when nothing was deleted

Map between `dbgen` row types and `domain.Project` in the repository so upper layers never import `dbgen` directly.

### 3. Implement the Task Repository

Create `server/internal/repositories/task_repository.go` following the same constructor pattern.

- Store a `*dbgen.Queries` field initialised once in the constructor.
- `ListTasksByProject` → `q.ListTasksByProject(ctx, projectID)` — parse `projectID` string to int64
- `CreateTask` → `q.CreateTask(ctx, params)` — include all fields: `title`, `description`, `status`, `priority`; parse `projectID` string to int64
- `UpdateTask` → `q.UpdateTask(ctx, params)` — parse string ID to int64; return `ErrNotFound` when no row is updated
- `DeleteTask` → `q.DeleteTask(ctx, id)` — parse string ID to int64; return `ErrNotFound` when nothing was deleted

### 4. Wire Real Repositories in main.go

Update `server/cmd/api/main.go` to:

- Open a `*sql.DB` connection using `DATABASE_URL` from config.
- Verify connectivity with `db.Ping()` at startup; log and exit if it fails.
- Replace `NewProjectRepositoryStub()` and `NewTaskRepositoryStub()` with the real constructor calls.
- Close the DB connection on shutdown.

### 5. Add Proper Error Handling in Handlers

Update `server/internal/http/handlers/project_handler.go` so that handlers distinguish between error types instead of always returning 501:

- Not-found errors → HTTP 404 with `{"error": "not found"}`
- Bad request/decode errors (already returning 400) — keep as-is
- All other errors → HTTP 500 with `{"error": "internal server error"}`

Define a sentinel error (e.g., `repositories.ErrNotFound`) and check for it in handlers using `errors.Is`.

### 6. Add Structured Logging

`router.go` already registers `middleware.Logger` from chi, which covers HTTP request logging. The remaining work is in `main.go`:

- Replace the existing `log.Printf` / `log.Fatalf` calls with `log/slog` equivalents so startup events (server address, DB ping result, fatal errors) are emitted as structured JSON log lines.
- Use Go's standard `log/slog` package (available since Go 1.21, compatible with Go 1.23 used here); no third-party logging library is needed.

### 7. Run Goose Migrations

Before the first end-to-end run, apply the initial migration to the running PostgreSQL instance:

```bash
goose -dir ./db/migrations postgres "<DATABASE_URL>" up
```

This creates the `projects` and `tasks` tables with the correct schema and constraints.
