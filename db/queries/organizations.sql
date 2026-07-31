-- name: CreateOrganization :one
INSERT INTO organizations (name, owner_id)
VALUES ($1, $2)
RETURNING id, name, owner_id, created_at;

