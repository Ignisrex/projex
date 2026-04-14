# Frontend Tasks

## Current State

The Next.js app scaffolding is in place. The projects list page (`app/page.js`) renders a list of projects and includes a "Create Project" form. The API client (`lib/api.js`) covers `getHealth`, `getProjects`, `createProject`, `getProject`, and `getProjectTasks` / `createTask`. Several API functions are missing, and the project detail page does not exist yet.

## Tasks

### 1. Complete the API Client (`lib/api.js`)

Add the missing functions so every backend endpoint has a corresponding client method:

- `updateProject(id, payload)` — `PUT /projects/{id}`
- `deleteProject(id)` — `DELETE /projects/{id}`
- `updateTask(id, payload)` — `PUT /tasks/{id}`
- `deleteTask(id)` — `DELETE /tasks/{id}`

All functions must go through the existing `request()` helper. No component should call `fetch` directly.

### 2. Create the Project Detail Page

Create `ui/app/projects/[id]/page.js`. This is the only missing page. It should:

- Fetch the project by ID on mount using `getProject(id)`.
- Display the project name and description.
- Fetch and display the task list using `getProjectTasks(id)`.
- Show each task's title, description, status, and priority.
- Provide a "Back to projects" link that returns to `/`.

Because the page requires client-side state (loading, form inputs), mark it with `"use client"` and use `useEffect` and `useState` hooks, consistent with the pattern used in `app/page.js`.

### 3. Add a Create Task Form to the Detail Page

Embed a form on the project detail page that lets the user create a new task:

- Fields: `title` (required), `description`, `status` (select: `todo` / `in_progress` / `done`, default `todo`), `priority` (number input, default `1`).
- On submit, call `createTask(projectId, payload)` and reload the task list.
- Display inline validation / error messages on failure.

### 4. Add Edit and Delete for Tasks

On each task in the list:

- **Delete**: Show a "Delete" button. On click, call `deleteTask(task.id)` and reload the task list.
- **Edit**: Show an "Edit" button that toggles an inline edit form pre-populated with the task's current values. On submit, call `updateTask(task.id, payload)` and reload.

### 5. Add Edit and Delete for Projects

On the project detail page:

- **Edit**: Show a form (or toggle) to update the project's name and description. On submit, call `updateProject(id, payload)` and refresh the displayed project.
- **Delete**: Show a "Delete Project" button. On click, call `deleteProject(id)` and redirect to `/` after success.

### 6. Add Edit and Delete for Projects on the List Page

On `app/page.js`, beside each project link:

- Add a "Delete" button that calls `deleteProject(project.id)` and refreshes the list.

### 7. Handle Loading and Error States

On both pages:

- Show a loading indicator while async calls are in flight.
- Display a user-readable error message (already partially done on the list page — extend the same pattern to the detail page).
- Clear errors before each new fetch attempt.
