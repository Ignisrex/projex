---
name: projex-frontend
description: Implement the Next.js frontend for the projex Project Tracker app. Use this skill whenever working on the UI, pages, components, API client, or any task in the ui/ directory. Invoke it for tasks like "build the frontend", "create the project detail page", "add task editing", "complete the API client", or "implement the UI".
---

# projex Frontend Implementation

Full task list lives in `plans/frontend.md`. Read it first to understand what remains and pick up where work left off.

## Conventions (strictly enforced)

- All API calls go through `ui/lib/api.js` only. Components must never call `fetch` directly.
- API base URL comes from `process.env.NEXT_PUBLIC_API_BASE_URL` (default `http://localhost:8080`) — already wired in `lib/api.js`.
- Use the Next.js App Router (`ui/app/` directory).
- Pages that need client-side state (data fetching, forms) must have `"use client"` at the top.
- No advanced state management libraries (no Redux, Zustand, etc.) — `useState` / `useEffect` only, following the pattern in `ui/app/page.js`.

## Task 1 — Complete the API Client

Edit `ui/lib/api.js` to add the four missing functions. All must use the existing `request()` helper:

```js
export async function updateProject(id, payload) { /* PUT /projects/{id} */ }
export async function deleteProject(id) { /* DELETE /projects/{id} */ }
export async function updateTask(id, payload) { /* PUT /tasks/{id} */ }
export async function deleteTask(id) { /* DELETE /tasks/{id} */ }
```

`deleteProject` and `deleteTask` call `DELETE` endpoints that return 204 — `request()` already returns `null` for 204 responses.

## Task 2 — Create the Project Detail Page

Create `ui/app/projects/[id]/page.js`. This page does not exist yet.

- Mark `"use client"` at the top.
- Extract `id` from the route using Next.js App Router's `useParams()` hook.
- On mount (`useEffect`), call `getProject(id)` and `getProjectTasks(id)` in parallel.
- Render the project name and description.
- Render the task list: title, description, status, priority.
- Include a "← Back to projects" link to `/`.
- Follow the structure and style of the existing `ui/app/page.js`.

## Task 3 — Create Task Form

On the project detail page, embed a form to create a new task below the task list:

- Fields: `title` (text, required), `description` (text), `status` (select: `todo` / `in_progress` / `done`, default `todo`), `priority` (number, default `1`, min `1`).
- On submit, call `createTask(id, { title, description, status, priority })` and reload the task list.
- Clear the form fields on success.
- Show an inline error message on failure.

## Task 4 — Task Edit and Delete

On each task row in the detail page:

- **Delete button**: calls `deleteTask(task.id)` then reloads the task list.
- **Edit button**: toggles an inline form pre-populated with the task's current values. On submit, calls `updateTask(task.id, payload)` and reloads. On cancel, dismisses the form without saving.
- Only one task should be in edit mode at a time (track `editingTaskId` in state).

## Task 5 — Project Edit and Delete

On the project detail page:

- **Edit**: toggle an inline form with the project's current name and description. On submit, call `updateProject(id, payload)` and refresh the displayed project data.
- **Delete Project button**: call `deleteProject(id)`, then redirect to `/` using Next.js `useRouter().push('/')`.

## Task 6 — Delete on the Projects List Page

On `ui/app/page.js`, add a delete button beside each project link:

- On click, call `deleteProject(project.id)` and reload the project list.

## Task 7 — Loading and Error States

Apply consistently on both pages:

- Use a boolean `loading` state; show a simple loading message while fetches are in progress.
- Use an `error` string state; display it in red (the list page already does this — replicate the pattern).
- Clear `error` at the start of every fetch/submit.

## Running Locally

```bash
cd ui
npm install
npm run dev
```

UI is available at `http://localhost:3000`. Requires the backend to be running at `http://localhost:8080` (or set `NEXT_PUBLIC_API_BASE_URL` in `ui/.env.local`).
