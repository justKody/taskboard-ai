package organization

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justKody/taskboard-go-api/middleware"
	"github.com/justKody/taskboard-go-api/service/membership"
	"github.com/justKody/taskboard-go-api/service/user"
)

type Handler struct {
	store           OrganizationStore
	membershipStore membership.MemebershipStore
	userStore       user.UserStore
}

func NewHandler(store OrganizationStore, membershipStore membership.MemebershipStore, userStore user.UserStore) *Handler {
	return &Handler{
		store:           store,
		membershipStore: membershipStore,
		userStore:       userStore,
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
	// user invites
	organizationRouter.HandleFunc("/invite", h.HandleInviteUserToOrganization).Methods(http.MethodPost)
	organizationRouter.HandleFunc("/invite/{id}", h.HandleAcceptOrRejectOrganizationInvite).Methods(http.MethodPut)
}
