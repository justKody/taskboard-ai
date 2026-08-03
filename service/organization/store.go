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
	CreateOrganization(ctx context.Context, params sqlc.CreateOrganizationParams) (*types.Organization, error)
	GetOrganizationById(ctx context.Context, id string) (*types.Organization, error)
	DeleteOrganizationById(ctx context.Context, id string) error
	CheckIfUserIsOwner(ctx context.Context, userId, id string) (bool, error)
	ChangeOwnerOfOrganizationById(ctx context.Context, userId, id string) error
	GetOrganizationsListByUserId(ctx context.Context, userId string) ([]types.Organization, error)
	CreateOrganizationInvite(ctx context.Context, params sqlc.CreateOrganizationInviteParams) (*types.OrganizationInvite, error)
	GetPendingOrganizationInvite(ctx context.Context, organizationId, userId string) (*types.OrganizationInvite, error)
	GetOrganizationInviteById(ctx context.Context, id string) (*types.OrganizationInvite, error)
	UpdateOrganizationInviteStatus(ctx context.Context, id string, status sqlc.OrganizationInviteStatus) (*types.OrganizationInvite, error)
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		queries: sqlc.New(db),
	}
}

func (s *Store) CreateOrganization(ctx context.Context, params sqlc.CreateOrganizationParams) (*types.Organization, error) {
	organization, err := s.queries.CreateOrganization(ctx, params)
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

func (s *Store) GetOrganizationById(ctx context.Context, id string) (*types.Organization, error) {
	organization, err := s.queries.GetOrganizationById(ctx, id)
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

func (s *Store) CheckIfUserIsOwner(ctx context.Context, userId string, id string) (bool, error) {
	ownerId, err := s.queries.CheckIfUserIsOwner(ctx, id)
	if err != nil {
		return false, err
	}
	return ownerId == userId, nil
}

func (s *Store) DeleteOrganizationById(ctx context.Context, id string) error {
	err := s.queries.DeleteOrganizationById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ChangeOwnerOfOrganizationById(ctx context.Context, userId, id string) error {
	err := s.queries.ChangeOwnerOfOrganizationById(ctx, sqlc.ChangeOwnerOfOrganizationByIdParams{
		OwnerID: userId,
		ID:      id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetOrganizationsListByUserId(ctx context.Context, userId string) ([]types.Organization, error) {
	organizations, err := s.queries.GetOrganizationsListByUserId(ctx, userId)
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

func (s *Store) CreateOrganizationInvite(ctx context.Context, params sqlc.CreateOrganizationInviteParams) (*types.OrganizationInvite, error) {
	invite, err := s.queries.CreateOrganizationInvite(ctx, params)
	if err != nil {
		return nil, err
	}
	return toOrganizationInvite(invite), nil
}

func (s *Store) GetPendingOrganizationInvite(ctx context.Context, organizationId, userId string) (*types.OrganizationInvite, error) {
	invite, err := s.queries.GetPendingOrganizationInvite(ctx, sqlc.GetPendingOrganizationInviteParams{
		OrganizationID: organizationId,
		UserID:         userId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toOrganizationInvite(invite), nil
}

func (s *Store) GetOrganizationInviteById(ctx context.Context, id string) (*types.OrganizationInvite, error) {
	invite, err := s.queries.GetOrganizationInviteById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toOrganizationInvite(invite), nil
}

func (s *Store) UpdateOrganizationInviteStatus(ctx context.Context, id string, status sqlc.OrganizationInviteStatus) (*types.OrganizationInvite, error) {
	invite, err := s.queries.UpdateOrganizationInviteStatus(ctx, sqlc.UpdateOrganizationInviteStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toOrganizationInvite(invite), nil
}

func toOrganizationInvite(invite sqlc.OrganizationInvite) *types.OrganizationInvite {
	return &types.OrganizationInvite{
		Id:             invite.ID,
		OrganizationId: invite.OrganizationID,
		UserId:         invite.UserID,
		InvitedBy:      invite.InvitedBy,
		Status:         string(invite.Status),
		CreatedAt:      invite.CreatedAt.Time,
	}
}
