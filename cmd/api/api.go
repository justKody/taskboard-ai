package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/justKody/taskboard-go-api/middleware"
	"github.com/justKody/taskboard-go-api/service/membership"
	"github.com/justKody/taskboard-go-api/service/organization"
	user "github.com/justKody/taskboard-go-api/service/user"
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
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)

	organizationStore := organization.NewStore(s.db)
	membershipStore := membership.NewStore(s.db)
	organizationHandler := organization.NewHandler(organizationStore, membershipStore, userStore)
	organizationHandler.RegisterRoutes(subRouter)
	// all handling

	fmt.Printf("\n\n🚀 Server starting on http://localhost:%s\n\n\n", s.addr)

	err := http.ListenAndServe(s.addr, router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
