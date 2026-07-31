package organization

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justKody/taskboard-go-api/middleware"
)

type Handler struct {
	store *Store
}

func (h *Handler) NewHandler(store *Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	organizationRouter := router.PathPrefix("/organization").Subrouter()

	// we will have middleware here to check if the user is authenticated
	organizationRouter.Use(middleware.Auth)

	// create organization
	organizationRouter.HandleFunc("/create", h.HandleCreateOrganization).Methods(http.MethodPost)
	organizationRouter.HandleFunc("/get/{id}", h.HandleGetOrganizationById).Methods(http.MethodGet)
}
