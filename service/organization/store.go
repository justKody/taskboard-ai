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
	GetOrganizationById(id string) (*types.Organization, error)
	DeleteOrganizationById(id string) error
	CheckIfUserIsOwner(userId, id string) (bool, error)
	ChangeOwnerOfOrganizationById(userId, id string) error
	GetOrganizationsListByUserId(userId string) ([]types.Organization, error)
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

func (s *Store) CheckIfUserIsOwner(userId string, id string) (bool, error) {
	ownerId, err := s.queries.CheckIfUserIsOwner(context.Background(), id)
	if err != nil {
		return false, err
	}
	return ownerId == userId, nil
}

func (s *Store) DeleteOrganizationById(id string) error {
	err := s.queries.DeleteOrganizationById(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ChangeOwnerOfOrganizationById(userId, id string) error {
	err := s.queries.ChangeOwnerOfOrganizationById(context.Background(), sqlc.ChangeOwnerOfOrganizationByIdParams{
		OwnerID: userId,
		ID:      id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetOrganizationsListByUserId(userId string) ([]types.Organization, error) {
	organizations, err := s.queries.GetOrganizationsListByUserId(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	result := make([]types.Organization, len(organizations))
	for i, org := range organizations {
		result[i] = types.Organization{
			Id:        org.ID,
			Name:      org.Name,
			OwnerId:   org.OwnerID,
			CreatedAt: org.CreatedAt.Time,
		}
	}
	return result, nil
}
