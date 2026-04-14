---
name: projex-testing
description: Write and run tests for the projex Project Tracker app. Use this skill whenever adding tests, setting up Jest, writing Go tests, creating handler or repository tests, writing component tests, or setting up end-to-end tests. Invoke it for any task like "write tests", "add unit tests", "set up Jest", "test the API", or "write E2E tests".
---

# projex Testing

Full task list lives in `plans/testing.md`. Read it first to understand what's been done and what remains.

Complete backend implementation (repositories and handlers) before writing tests that require a real database.

## Backend Tests (Go)

### Repository Tests — `server/internal/repositories/`

Tests run against a real PostgreSQL instance. Use the `DATABASE_URL` environment variable (or a dedicated test DB). Use `TestMain` to apply goose migrations before the suite and truncate tables between tests.

Cover for both `ProjectRepository` and `TaskRepository`:

- `CreateAndGet` — create a record, fetch by ID, assert all fields match.
- `List` — create multiple records, assert correct order (`created_at DESC`).
- `Update` — create, update, assert returned row reflects changes.
- `Delete` — create, delete, assert `ErrNotFound` on subsequent fetch.
- `DeleteCascades` (projects only) — create project + tasks, delete project, assert tasks are gone (ON DELETE CASCADE).

ID types: repository interfaces use `string`; the DB uses `bigint`. The repository implementation converts with `strconv.ParseInt` / `strconv.FormatInt`.

### Service Tests — `server/internal/services/`

No database needed. Use hand-written fakes (structs that implement `ProjectRepository` / `TaskRepository`) — no mocking library. Verify:

- Methods delegate to the correct repository calls.
- `ErrNotFound` from a repository surfaces unchanged to the caller.

### Handler Tests — `server/internal/http/handlers/`

Use `net/http/httptest`. Inject fake services so no database is needed.

For every endpoint assert:
- Happy path status code and response body shape.
- `ErrNotFound` from the service → handler returns HTTP 404.
- Malformed JSON body → HTTP 400 with `{"error": "invalid request payload"}`.
- Creates return HTTP 201; deletes return HTTP 204.

Endpoints to cover:
`GET /health`, `GET /projects`, `GET /projects/{id}`, `POST /projects`, `PUT /projects/{id}`, `DELETE /projects/{id}`, `GET /projects/{id}/tasks`, `POST /projects/{id}/tasks`, `PUT /tasks/{id}`, `DELETE /tasks/{id}`

### Error Sentinel Test

Verify `errors.Is(repositories.ErrNotFound, repositories.ErrNotFound)` is true and that the handler translates it to HTTP 404, guarding against future refactoring.

### Run Backend Tests

```bash
# From server/
go test ./...
go test -v ./...
go test -race ./...
```

## Frontend Tests (React / Next.js)

Jest is **not** bundled with Next.js 14. Install and configure it first:

```bash
cd ui
npm install --save-dev jest jest-environment-jsdom @testing-library/react @testing-library/jest-dom @testing-library/user-event
```

Create `ui/jest.config.js` and `ui/jest.setup.js` per the [Next.js Jest docs](https://nextjs.org/docs/app/building-your-application/testing/jest).

### API Client Tests — `ui/lib/api.js`

Mock `fetch` globally via `jest.fn()` or use `msw` (Mock Service Worker). For each exported function assert:

- Correct URL and HTTP method are used.
- A non-2xx response throws with a message containing the status code.
- A 204 response returns `null` without parsing JSON.

### Projects List Page — `ui/app/page.js`

Use React Testing Library:

- On render, `getProjects` is called and returned projects appear as links.
- Submitting the create form calls `createProject` with the correct payload.
- A failed API call renders an error message.

### Project Detail Page — `ui/app/projects/[id]/page.js`

- On render, `getProject` and `getProjectTasks` are called with the correct ID.
- Tasks render with title, status, and priority.
- Submitting the create-task form calls `createTask`.
- Delete task button calls `deleteTask` and removes the task.
- Delete Project button calls `deleteProject` and triggers navigation to `/`.

### Run Frontend Tests

```bash
cd ui
npm test
npm test -- --coverage
```

## End-to-End Tests

Start the full stack first:

```bash
docker compose up --build -d
```

Use Playwright (`npm init playwright@latest` in a new `e2e/` directory at the repo root).

### Happy-Path Scenarios

1. Visit `/`, create a project, assert it appears in the list.
2. Click the project link, assert the detail page loads with the project name.
3. Create a task, assert it appears in the task list.
4. Edit the task status to `in_progress`, assert the displayed status updates.
5. Delete the task, assert it disappears.
6. Delete the project, assert redirect to `/` and project is gone from the list.

### Edge Cases

- Navigate to `/projects/99999` (non-existent ID) — assert a "not found" message, not a crash.
- Submit create-project form with empty name — assert submission is blocked by HTML required validation.
- Submit create-task form with empty title — assert same.

### Run E2E Tests

```bash
npx playwright test
```
