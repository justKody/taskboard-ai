-- name: CreateProject :one
INSERT INTO projects (organization_id, name, description, created_by)
VALUES ($1, $2, $3, $4)
RETURNING id, organization_id, name, description, status, created_by, created_at;