-- name: GetUserByEmail :one
SELECT id, name, email, password, created_at
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING id, name, email, password, created_at;
