package user

import (
	"github.com/jackc/pgx/v5"
	"github.com/justKody/taskboard-go-api/types"
)

type Store struct {
	db *pgx.Conn
}

type UserStore interface {
	GetUserByID(id string) (*types.User, error)
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUserByID(id string) (*types.User, error) {
	// query := "SELECT id, name, email, password FROM users WHERE id = $1"
	// row := s.db.QueryRow(context.Background(), query, id)
	// var user types.User
	// err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	// return &user, err
	return nil, nil
}
