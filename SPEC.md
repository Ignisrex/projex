# Project Tracker — Product & Technical Specification (v1)
## 1. Overview
Project Tracker is a full-stack web application for managing projects and their tasks.
The system enables users to:
- Create, edit, view, and delete projects
- Create, edit, view, and delete tasks associated with a project
- Track task status and priority within each project

This specification defines architecture, data model, API surface, project structure, containerization, and delivery expectations for v1.

## 2. Goal
Build a basic project and task management application with a clear separation of concerns across frontend, backend, and database layers, using type-safe SQL access in Go.

## 3. Tech Stack
- Backend: Go
- Frontend: Next.js (React)
- Database: PostgreSQL
- SQL tooling: sqlc (typed query generation), goose (migrations)
- Containerization: Docker (per service)
- Orchestration: docker-compose

## 4. Repository Structure
```text
/
├─ server/               # Go REST API
│  └─ Dockerfile
├─ ui/                   # Next.js frontend
│  └─ Dockerfile
├─ db/                   # Database layer
│  ├─ migrations/        # goose migration files
│  └─ queries/           # raw SQL files consumed by sqlc
└─ docker-compose.yml
```

## 5. Functional Requirements
### 5.1 Frontend (`/ui`)
- Build a Next.js app with minimal, clean UI
- Provide:
  - Page: list all projects
  - Page: view one project and its tasks
  - Forms: create/edit projects
  - Forms: create/edit tasks
- Integrate with backend via `fetch` (or a simple API client wrapper)
- Frontend must never connect directly to PostgreSQL

### 5.2 Backend (`/server`)
- Build a Go REST API with modular structure:
  - handlers (HTTP layer)
  - services (business logic)
  - repositories/data access (sqlc-generated layer integration)
- Required endpoints:
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
- Include basic structured logging
- Configuration must come from environment variables
- Backend must only access the database through sqlc-generated code

### 5.3 Database (`/db`)
- Use PostgreSQL
- Manage schema with goose migrations in `db/migrations`
- Define raw SQL queries in `db/queries`
- Configure sqlc to generate typed Go database access code from those queries

## 6. Data Model
### 6.1 `projects` table
- `id` (uuid or serial primary key)
- `name`
- `description`
- `created_at`
- `updated_at`

### 6.2 `tasks` table
- `id` (primary key)
- `project_id` (foreign key → `projects.id`)
- `title`
- `description`
- `status` (`todo` | `in_progress` | `done`)
- `priority` (integer or enum)
- `created_at`
- `updated_at`

### 6.3 Relational Constraints
- One project has many tasks
- Every task must belong to exactly one project
- Deleting a project should define explicit task behavior (recommended for v1: `ON DELETE CASCADE`)

## 7. API Contract (v1)
### 7.1 Health
- `GET /health`
  - Response: application health/status payload

### 7.2 Projects
- `GET /projects`
  - Returns list of projects
- `GET /projects/{id}`
  - Returns one project by ID
- `POST /projects`
  - Creates a project
- `PUT /projects/{id}`
  - Updates a project
- `DELETE /projects/{id}`
  - Deletes a project

### 7.3 Tasks
- `GET /projects/{id}/tasks`
  - Returns all tasks for a project
- `POST /projects/{id}/tasks`
  - Creates task under a project
- `PUT /tasks/{id}`
  - Updates a task
- `DELETE /tasks/{id}`
  - Deletes a task

## 8. sqlc Requirements
- Keep SQL query files under `db/queries`
- Configure sqlc for PostgreSQL and Go code generation
- Generate strongly typed query methods and models
- Backend data access must call sqlc-generated types/methods only (no ad hoc SQL in handlers/services)

## 9. goose Requirements
- All schema changes must be represented as migration files in `db/migrations`
- Include an initial migration that creates:
  - `projects`
  - `tasks`
  - relevant indexes/constraints (including FK constraint)

## 10. Containerization & Orchestration
### 10.1 Dockerfiles
- `server/Dockerfile` required
- `ui/Dockerfile` required
- Each image should be service-specific and self-contained

### 10.2 docker-compose (`docker-compose.yml`)
Define services:
- `postgres`
- `server`
- `ui`

Compose configuration must include:
- Environment variables for each service
- Persistent volume for PostgreSQL data
- Service dependencies (`server` depends on `postgres`; `ui` depends on `server`)

### 10.3 Networking Requirements
Create two docker networks:
- Public network: `ui` ↔ `server`
- Private network: `server` ↔ `postgres`

Rules:
- `ui` and `server` are exposed on the public network
- `server` and `postgres` communicate on the private network
- PostgreSQL must not be publicly exposed
- Frontend and backend ports are exposed (e.g., `3000`, `8080`)

## 11. Configuration
Use environment variables for all runtime config, including (at minimum):
- API bind address/port
- Database connection URL or equivalent DB config fields
- Frontend API base URL
- Environment mode (dev/prod as needed)

No hardcoded secrets or host-specific values in source code.

## 12. Code Quality & Architecture Constraints
- Keep architecture simple, readable, and modular
- Separate concerns cleanly (HTTP handlers, business logic, data access, models)
- Favor clarity over optimization
- Avoid overengineering and unnecessary abstractions in v1
- Build incrementally and keep the system runnable after each major step

## 13. Implementation Order (Required)
1. Define schema and create initial goose migration
2. Write SQL queries for projects and tasks
3. Configure sqlc and generate Go code
4. Build Go API using sqlc layer
5. Scaffold Next.js UI
6. Connect UI to backend APIs
7. Add Dockerfiles for `server` and `ui`
8. Add `docker-compose.yml` with required networks/services
9. Verify end-to-end functionality

## 14. Non-goals (v1)
- No authentication/authorization
- No background jobs
- No websockets
- No advanced frontend state management

## 15. Deliverables
- `server/` (Go API)
- `ui/` (Next.js app)
- `db/` (migrations + sqlc queries/config)
- `docker-compose.yml`
- `server/Dockerfile`
- `ui/Dockerfile`
- `README` documenting:
  - setup steps
  - docker-compose usage
  - migration commands
  - sqlc generation commands
