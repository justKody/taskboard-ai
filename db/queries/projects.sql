-- name: CreateProject :one
INSERT INTO projects (organization_id, name, description, created_by)
VALUES ($1, $2, $3, $4)
RETURNING id, organization_id, name, description, status, created_by, created_at;

-- name: GetProject :one
SELECT id, organization_id, name, description, status, created_by, created_at FROM projects
WHERE id = $1;

-- name: ListProject :many
SELECT id, name, description, status, created_by, created_at from projects
where organization_id = $1;

-- name: UpdateProject :one
UPDATE projects
SET name = $3, description = $4, status = $5
WHERE id = $1 AND organization_id = $2
RETURNING id, organization_id, name, description, status, created_by, created_at;

-- name: DeleteProject :exec
DELETE from projects
where id = $1;