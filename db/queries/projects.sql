-- name: CreateProject :one
INSERT INTO projects (organization_id, name, description, created_by)
VALUES ($1, $2, $3, $4)
RETURNING id, organization_id, name, description, status, created_by, created_at;

-- name: ListProject :many
SELECT id, name, description, status, created_by, created_at from projects
where organization_id = $1;

-- name: DeleteProject :exec
DELETE from projects
where id = $1;