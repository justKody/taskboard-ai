package project

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/justKody/taskboard-go-api/db/sqlc"
	"github.com/justKody/taskboard-go-api/types"
)

type Store struct {
	queries *sqlc.Queries
}

type ProjectStore interface {
	CreateProject(ctx context.Context, params sqlc.CreateProjectParams) (*types.Project, error)
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		queries: sqlc.New(db),
	}
}

func (s *Store) CreateProject(ctx context.Context, params sqlc.CreateProjectParams) (*types.Project, error) {
	project, err := s.queries.CreateProject(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &types.Project{
		Id:             project.ID,
		OrganizationID: project.OrganizationID,
		Name:           project.Name,
		Description:    project.Description.String,
		Status:         string(project.Status),
		CreatedBy:      project.CreatedBy,
		CreatedAt:      project.CreatedAt.Time,
	}, nil
}
