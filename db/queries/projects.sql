-- name: ListProjects :many
SELECT id, name, description, created_at, updated_at
FROM projects
ORDER BY created_at DESC;

-- name: GetProject :one
SELECT id, name, description, created_at, updated_at
FROM projects
WHERE id = $1;

-- name: CreateProject :one
INSERT INTO projects (name, description)
VALUES ($1, $2)
RETURNING id, name, description, created_at, updated_at;

-- name: UpdateProject :one
UPDATE projects
SET name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, created_at, updated_at;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;
