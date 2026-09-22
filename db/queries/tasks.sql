-- name: CreateTask :one
INSERT INTO tasks (project_id, title, description, status, priority, due_date, assigned_to, created_by)
VALUES (
  sqlc.arg(project_id),
  sqlc.arg(title),
  sqlc.arg(description),
  COALESCE(sqlc.narg(status)::task_status, 'todo'),
  COALESCE(sqlc.narg(priority)::priority, 'low'),
  sqlc.arg(due_date),
  sqlc.arg(assigned_to),
  sqlc.arg(created_by)
)
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
DELETE from tasks
where id = $1;

-- name: ListProjectTask :many
SELECT id, project_id, title, description, status, priority, due_date, assigned_to, created_by, created_at FROM tasks
where project_id = $1;