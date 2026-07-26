package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/justKody/taskboard-go-api/middleware"
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

	router.Use(middleware.Logger)
	_ := router.PathPrefix("/api/v1").Subrouter()

	// all handling

	log.Printf("\n🚀 Server starting on http://localhost:%s", s.addr)

	err := http.ListenAndServe(s.addr, router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
