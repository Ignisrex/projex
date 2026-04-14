# Project Tracker (Scaffold)
This repository contains a full-stack scaffold for Project Tracker with:
- `server/`: Go + chi REST API skeleton
- `ui/`: Next.js frontend skeleton
- `db/`: PostgreSQL migration and sqlc query scaffolding
- `docker-compose.yml`: local orchestration

Current state:
- Endpoints and pages are wired, but core business/data logic is intentionally not implemented yet.
- Repository/service abstractions are in place for incremental implementation.

## Project structure
```text
.
├─ server/
├─ ui/
├─ db/
│  ├─ migrations/
│  └─ queries/
├─ docker-compose.yml
├─ SPEC.md
└─ README.md
```

## Prerequisites
- Docker + Docker Compose
- Go 1.23+ (for local server runs)
- Node 20+ (for local UI runs)
- goose CLI (for migrations)
- sqlc CLI (for typed query generation)

## Run with Docker Compose
From repository root:
```bash
docker compose up --build
```

Services:
- UI: `http://localhost:3000`
- API: `http://localhost:8080`
- DB: internal only (not exposed publicly)

## Local development (without Compose)
### Server
```bash
cd server
go mod tidy
go run ./cmd/api
```

### UI
```bash
cd ui
npm install
npm run dev
```

## Migrations (goose)
Example command from repo root (adjust for your local goose setup):
```bash
goose -dir ./db/migrations postgres "postgres://postgres:postgres@localhost:5432/project_tracker?sslmode=disable" up
```

## Regenerate sqlc code
From repo root:
```bash
sqlc generate -f ./db/sqlc.yaml
```

Generated code output path:
- `server/internal/dbgen`

## API scaffolded routes
- `GET /health`
- `GET /projects`
- `GET /projects/{id}`
- `POST /projects`
- `PUT /projects/{id}`
- `DELETE /projects/{id}`
- `GET /projects/{id}/tasks`
- `POST /projects/{id}/tasks`
- `PUT /tasks/{id}`
- `DELETE /tasks/{id}`
