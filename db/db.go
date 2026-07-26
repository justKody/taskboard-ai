package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func NewPostgresStorage(connStr string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {

		log.Fatal("Error Connecting to Db", err)
	}

	return conn
}
