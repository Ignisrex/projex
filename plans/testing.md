# Testing Tasks

## Overview

Testing is divided into three layers: backend unit/integration tests (Go), frontend component tests (React), and end-to-end tests covering the full stack. Complete the backend implementation first before writing tests that depend on real database access.

## Backend Tests (Go)

### 1. Repository Unit Tests with a Test Database

Once the real repository implementations exist, write tests in `server/internal/repositories/` that run against a real PostgreSQL instance (use the same Docker Compose DB or a dedicated test DB):

- `TestProjectRepository_CreateAndGet` — create a project, then fetch it by ID and assert all fields match.
- `TestProjectRepository_List` — create multiple projects, assert all are returned in descending `created_at` order.
- `TestProjectRepository_Update` — create a project, update name and description, assert the returned row reflects the changes.
- `TestProjectRepository_Delete` — create a project, delete it, assert `GetProject` returns `ErrNotFound`.
- `TestProjectRepository_DeleteCascades` — create a project with tasks, delete the project, assert the tasks are gone (ON DELETE CASCADE).
- Mirror the same set for `TaskRepository`.

Use `TestMain` to run migrations before the test suite starts and roll them back (or truncate tables) between tests.

### 2. Service Unit Tests with Mock Repositories

Write tests in `server/internal/services/` using mock or fake repository implementations (no database required):

- Verify that service methods delegate correctly to the repository layer.
- Test error propagation: when a repository returns `ErrNotFound`, confirm the service surfaces the same error to the handler.
- Since services currently have minimal logic, these tests primarily guard against accidental changes to the delegation contract.

Use a hand-written fake (struct that implements the interface) rather than a mocking library to stay consistent with the project's no-framework philosophy.

### 3. Handler Integration Tests with `httptest`

Write tests in `server/internal/http/handlers/` using `net/http/httptest`. Inject fake services so these tests don't need a database:

- `GET /health` — assert status 200 and a JSON body with a `status` key.
- `GET /projects` — happy path returns 200 with an array; service error returns 500.
- `GET /projects/{id}` — found returns 200; service returns `ErrNotFound` → handler returns 404.
- `POST /projects` — valid body returns 201 with the created project; malformed JSON returns 400.
- `PUT /projects/{id}` — valid update returns 200; not found returns 404; bad body returns 400.
- `DELETE /projects/{id}` — success returns 204; not found returns 404.
- Repeat the same matrix for task endpoints (`GET /projects/{id}/tasks`, `POST /projects/{id}/tasks`, `PUT /tasks/{id}`, `DELETE /tasks/{id}`).

### 4. Error Sentinel Coverage

Add a test that confirms `errors.Is(repositories.ErrNotFound, repositories.ErrNotFound)` works as expected and that handlers correctly translate it to HTTP 404. This guards against future refactoring accidentally breaking the error chain.

## Frontend Tests (React / Next.js)

### 5. API Client Unit Tests

Write tests for `ui/lib/api.js` using `jest` with `fetch` mocked via `jest.fn()` or `msw` (Mock Service Worker). Jest is not bundled with Next.js 14 and must be installed separately (`jest`, `jest-environment-jsdom`, `@testing-library/react`, `@testing-library/jest-dom`) and configured via `jest.config.js` before any tests can run.

- Each exported function calls the correct URL and HTTP method.
- A non-2xx response causes the function to throw with a meaningful message.
- A 204 response returns `null` without attempting to parse JSON.

### 6. Component Tests for the Projects List Page

Test `ui/app/page.js` using React Testing Library:

- On mount, `getProjects` is called and the returned projects are rendered as links.
- Submitting the "Create Project" form calls `createProject` with the correct payload and refreshes the list.
- An API error renders an error message.

### 7. Component Tests for the Project Detail Page

Test `ui/app/projects/[id]/page.js`:

- On mount, `getProject` and `getProjectTasks` are called with the correct ID.
- Tasks are rendered with their title, status, and priority.
- Submitting the "Create Task" form calls `createTask` and refreshes the task list.
- Clicking "Delete" on a task calls `deleteTask` and removes the task from the list.
- Clicking "Delete Project" calls `deleteProject` and redirects to `/`.

## End-to-End Tests

### 8. Happy-Path E2E Scenarios

Use Playwright (recommended) or Cypress to drive a real browser against the full stack running via Docker Compose. Cover the critical user flows:

- **Create a project**: Visit `/`, submit the create form, assert the new project appears in the list, and its link navigates to the detail page.
- **View project detail**: Navigate to a project detail page and assert the project name and an empty task list are rendered.
- **Create a task**: On the detail page, fill in the create-task form and assert the new task appears in the list.
- **Update a task**: Click "Edit" on a task, change its status to `in_progress`, save, and assert the displayed status updates.
- **Delete a task**: Click "Delete" on a task and assert it disappears from the list.
- **Delete a project**: Click "Delete Project" and assert the app redirects to `/` and the project no longer appears in the list.

### 9. Error and Edge Case Scenarios

- Navigate to a non-existent project ID (e.g., `/projects/99999`) and assert a meaningful "not found" message is shown rather than a crash.
- Submit the create-project form with an empty name field and assert the required-field validation prevents submission.
- Submit the create-task form with an empty title and assert the same.

## Running Tests

### Backend

```bash
# From the server directory
go test ./...

# With verbose output
go test -v ./...

# With race detector
go test -race ./...
```

### Frontend

```bash
# From the ui directory
npm test

# With coverage
npm test -- --coverage
```

### End-to-End

```bash
# Start the full stack first
docker compose up --build -d

# Run Playwright tests (once configured)
npx playwright test
```
