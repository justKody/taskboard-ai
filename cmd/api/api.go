package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// all handling

	http.ListenAndServe(s.addr, router)
}
