package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type APIServer struct {
	addr string
	db   *pgx.Conn
}

func NewApiServer(addr string, db *pgx.Conn) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// subRouter := router.PathPrefix("/api/v1").Subrouter()

	// all handling

	http.ListenAndServe(s.addr, router)
}
