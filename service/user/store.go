package user

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

type UserStore interface {
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	CreateUser(ctx context.Context, params sqlc.CreateUserParams) (*types.User, error)
	GetUserById(ctx context.Context, id string) (*types.User, error)
	UpdateUser(ctx context.Context, params sqlc.UpdateUserParams) (*types.User, error)
	ChangePassword(ctx context.Context, params sqlc.ChangePasswordParams) (*types.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetUsersList(ctx context.Context) ([]types.User, error)
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		queries: sqlc.New(db),
	}
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
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

func (s *Store) CreateUser(ctx context.Context, params sqlc.CreateUserParams) (*types.User, error) {
	user, err := s.queries.CreateUser(ctx, params)
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

func (s *Store) GetUserById(ctx context.Context, id string) (*types.User, error) {
	user, err := s.queries.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
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

func (s *Store) UpdateUser(ctx context.Context, params sqlc.UpdateUserParams) (*types.User, error) {
	user, err := s.queries.UpdateUser(ctx, params)
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

func (s *Store) ChangePassword(ctx context.Context, params sqlc.ChangePasswordParams) (*types.User, error) {
	user, err := s.queries.ChangePassword(ctx, params)
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

func (s *Store) DeleteUser(ctx context.Context, id string) error {
	err := s.queries.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user not found to delete")
		}
		return errors.New("failed to delete user")
	}
	return nil
}

func (s *Store) GetUsersList(ctx context.Context) ([]types.User, error) {
	users, err := s.queries.GetUsersList(ctx)
	if err != nil {
		return nil, err
	}

	usersList := make([]types.User, len(users))
	for i, user := range users {
		usersList[i] = types.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: user.CreatedAt,
		}
	}

	return usersList, nil
}
