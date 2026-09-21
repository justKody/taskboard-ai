package project

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justKody/taskboard-go-api/middleware"
	"github.com/justKody/taskboard-go-api/service/membership"
)

type Handler struct {
	store           ProjectStore
	membershipStore membership.MemebershipStore
}

func NewHandler(store ProjectStore, membershipStore membership.MemebershipStore) *Handler {
	return &Handler{
		store:           store,
		membershipStore: membershipStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	projectRouter := router.PathPrefix("/project").Subrouter()

	projectRouter.Use(middleware.Auth)

	projectRouter.HandleFunc("/create", h.HandleCreateProject).Methods(http.MethodPost)
}
