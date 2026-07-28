package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/justKody/taskboard-go-api/db/sqlc"
	"github.com/justKody/taskboard-go-api/types"
)

type Store struct {
	queries *sqlc.Queries
}

type UserStore interface {
	GetUserByEmail(email string) (*types.User, error)
	CreateUser(params sqlc.CreateUserParams) (*types.User, error)
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		queries: sqlc.New(db),
	}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	user, err := s.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	return &types.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *Store) CreateUser(params sqlc.CreateUserParams) (*types.User, error) {
	user, err := s.queries.CreateUser(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return &types.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}, nil
}
