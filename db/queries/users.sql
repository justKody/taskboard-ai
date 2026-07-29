-- name: GetUserByEmail :one
SELECT id, name, email, password, created_at
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING id, name, email, password, created_at;

-- name: GetUserById :one
SELECT id, name, email, password, created_at
FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET name = $2
WHERE id = $1
RETURNING id, name, email, password, created_at;

-- name: ChangePassword :one
UPDATE users
SET password = $2
WHERE id = $1
RETURNING id, name, email, password, created_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUsersList :many
SELECT id, name, email, password, created_at
FROM users;