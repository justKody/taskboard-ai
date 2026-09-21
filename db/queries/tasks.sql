-- name: CreateTask :one
INSERT INTO tasks (project_id, title, description, status, priority, due_date, assigned_to, created_by)
VALUES ($1, $2, $3, COALESCE($4, 'todo'), COALESCE($5, 'low'), $6, $7, $8)
RETURNING id, project_id, title, description, status, priority, due_date, assigned_to, created_by, created_at;

-- name: GetTask :many
SELECT id, project_id, title, description, status, priority, due_date, assigned_to, created_by, created_at FROM tasks
where project_id = $1 AND assigned_to = $2;

-- name: UpdateTask :one
UPDATE tasks
SET title = $3, description = $4, status = $5, priority = $6, due_date = $7, assigned_to = $8
WHERE id = $1 AND project_id = $2
RETURNING id, project_id, title, description, status, priority, due_date, assigned_to, created_by, created_at;

-- name: DeleteTask :exec
DELETE from projects
where id = $1;

-- name: ListProjectTask :many
SELECT id, project_id, title, description, status, priority, due_date, assigned_to, created_by, created_at FROM tasks
where project_id = $1;