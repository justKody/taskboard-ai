package membership

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/justKody/taskboard-go-api/db/sqlc"
	"github.com/justKody/taskboard-go-api/types"
)

type Store struct {
	query *sqlc.Queries
}

type MemebershipStore interface {
	GetAllMembershipsByOrganizationId(ctx context.Context, orgId string) ([]types.Membership, error)
	UpdateMembershipRole(ctx context.Context, userId, orgId string, role sqlc.MembershipsRole) error
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		query: sqlc.New(db),
	}
}

func (s *Store) GetAllMembershipsByOrganizationId(ctx context.Context, orgId string) ([]types.Membership, error) {
	memberships, err := s.query.GetAllMembershipsByOrganizationId(ctx, orgId)
	if err != nil {
		return nil, err
	}

	membershipsList := make([]types.Membership, len(memberships))

	for i, membership := range memberships {

		role := string(membership.Role)
		joinedAt := membership.JoinedAt.Time.Format(time.RFC3339)
		membershipsList[i] = types.Membership{
			UserId:         membership.UserID,
			OrganizationId: membership.OrganizationID,
			Role:           role,
			JoinedAt:       joinedAt,
		}
	}

	return membershipsList, nil
}

func (s *Store) UpdateMembershipRole(ctx context.Context, userId, orgId string, role sqlc.MembershipsRole) error {
	return s.query.UpdateMembershipRole(ctx, sqlc.UpdateMembershipRoleParams{
		UserID:         userId,
		OrganizationID: orgId,
		Role:           role,
	})
}
