package organization

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justKody/taskboard-go-api/middleware"
)

type Handler struct {
}

func (h *Handler) NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	organizationRouter := router.PathPrefix("/organization").Subrouter()

	// we will have middleware here to check if the user is authenticated
	organizationRouter.Use(middleware.Auth)

	// create organization
	organizationRouter.HandleFunc("/create", h.CreateOrganization).Methods(http.MethodPost)
}
