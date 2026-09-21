package task

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/justKody/taskboard-go-api/db/sqlc"
)

type Store struct {
	queries *sqlc.Queries
}

type TaskStore interface {
	CreateTask(ctx context.Context, params sqlc.CreateTaskParams) (sqlc.Task, error)
	GetTask(ctx context.Context, projectId string, assignedTo pgtype.UUID) ([]sqlc.Task, error)
	UpdateTask(ctx context.Context, params sqlc.UpdateTaskParams) (sqlc.Task, error)
	DeleteTask(ctx context.Context, taskId string) error
	ListProjectTask(ctx context.Context, projectId string) ([]sqlc.Task, error)
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		queries: sqlc.New(db),
	}
}

func (s *Store) CreateTask(ctx context.Context, params sqlc.CreateTaskParams) (sqlc.Task, error) {
	return s.queries.CreateTask(ctx, params)
}

func (s *Store) GetTask(ctx context.Context, projectId string, assignedTo pgtype.UUID) ([]sqlc.Task, error) {
	return s.queries.GetTask(ctx, sqlc.GetTaskParams{
		ProjectID:  projectId,
		AssignedTo: assignedTo,
	})
}

func (s *Store) UpdateTask(ctx context.Context, params sqlc.UpdateTaskParams) (sqlc.Task, error) {
	return s.queries.UpdateTask(ctx, params)
}

func (s *Store) DeleteTask(ctx context.Context, taskId string) error {
	return s.queries.DeleteTask(ctx, taskId)
}

func (s *Store) ListProjectTask(ctx context.Context, projectId string) ([]sqlc.Task, error) {
	return s.queries.ListProjectTask(ctx, projectId)
}
