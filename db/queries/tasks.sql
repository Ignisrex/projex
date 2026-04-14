-- name: ListTasksByProject :many
SELECT id, project_id, title, description, status, priority, created_at, updated_at
FROM tasks
WHERE project_id = $1
ORDER BY created_at DESC;

-- name: CreateTask :one
INSERT INTO tasks (project_id, title, description, status, priority)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, project_id, title, description, status, priority, created_at, updated_at;

-- name: UpdateTask :one
UPDATE tasks
SET title = $2,
    description = $3,
    status = $4,
    priority = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING id, project_id, title, description, status, priority, created_at, updated_at;

-- name: DeleteTask :execrows
DELETE FROM tasks
WHERE id = $1;
