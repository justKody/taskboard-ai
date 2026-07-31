package organization

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

type OrganizationStore interface {
	CreateOrganization(sqlc.CreateOrganizationParams) (*types.Organization, error)
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		queries: sqlc.New(db),
	}
}

func (s *Store) CreateOrganization(params sqlc.CreateOrganizationParams) (*types.Organization, error) {
	organization, err := s.queries.CreateOrganization(context.Background(), params)
	if err != nil {
		return nil, err
	}
	return &types.Organization{
		Id:        organization.ID,
		Name:      organization.Name,
		OwnerId:   organization.OwnerID,
		CreatedAt: organization.CreatedAt.Time,
	}, nil
}

func (s *Store) GetOrganizationById(id string) (*types.Organization, error) {
	organization, err := s.queries.GetOrganizationById(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &types.Organization{
		Id:        organization.ID,
		Name:      organization.Name,
		OwnerId:   organization.OwnerID,
		CreatedAt: organization.CreatedAt.Time,
	}, nil
}
