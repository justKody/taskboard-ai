package organization

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justKody/taskboard-go-api/middleware"
	"github.com/justKody/taskboard-go-api/service/membership"
)

type Handler struct {
	store           OrganizationStore
	membershipStore membership.MemebershipStore
}

func NewHandler(store OrganizationStore, membershipStore membership.MemebershipStore) *Handler {
	return &Handler{
		store:           store,
		membershipStore: membershipStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	organizationRouter := router.PathPrefix("/organization").Subrouter()

	// we will have middleware here to check if the user is authenticated
	organizationRouter.Use(middleware.Auth)

	organizationRouter.HandleFunc("/create", h.HandleCreateOrganization).Methods(http.MethodPost)
	organizationRouter.HandleFunc("/get/{id}", h.HandleGetOrganizationDetailsById).Methods(http.MethodGet)
	organizationRouter.HandleFunc("/list", h.HandleGetOrganizationsListByUserId).Methods(http.MethodGet)
	organizationRouter.HandleFunc("/change-owner/{id}", h.HandleChangeOwnerOfOrganizationById).Methods(http.MethodPut)
	organizationRouter.HandleFunc("/delete/{id}", h.HandleDeleteOrganizationById).Methods(http.MethodDelete)
}
