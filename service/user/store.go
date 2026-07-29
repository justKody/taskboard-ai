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
	GetUserByEmail(email string) (*types.User, error)
	CreateUser(params sqlc.CreateUserParams) (*types.User, error)
	GetUserById(id string) (*types.User, error)
	UpdateUser(params sqlc.UpdateUserParams) (*types.User, error)
	ChangePassword(params sqlc.ChangePasswordParams) (*types.User, error)
	DeleteUser(id string) error
	GetUsersList() ([]types.User, error)
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		queries: sqlc.New(db),
	}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	user, err := s.queries.GetUserByEmail(context.Background(), email)
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

func (s *Store) GetUserById(id string) (*types.User, error) {
	user, err := s.queries.GetUserById(context.Background(), id)
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

func (s *Store) UpdateUser(params sqlc.UpdateUserParams) (*types.User, error) {
	user, err := s.queries.UpdateUser(context.Background(), params)
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

func (s *Store) ChangePassword(params sqlc.ChangePasswordParams) (*types.User, error) {
	user, err := s.queries.ChangePassword(context.Background(), params)
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

func (s *Store) DeleteUser(id string) error {
	err := s.queries.DeleteUser(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user not found to delete")
		}
		return errors.New("failed to delete user")
	}
	return nil
}

func (s *Store) GetUsersList() ([]types.User, error) {
	users, err := s.queries.GetUsersList(context.Background())
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
