package organization

import (
	"errors"
	"net/http"

	"github.com/justKody/taskboard-go-api/db/sqlc"
	"github.com/justKody/taskboard-go-api/middleware"
	"github.com/justKody/taskboard-go-api/utils"
)

func (c *Handler) HandleCreateOrganization(w http.ResponseWriter, r *http.Request) {
	var payload CreateOrganizationRequestDTO
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userId, ok := middleware.GetUserID(r.Context())
	if ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authenticated"))
		return
	}

	params := sqlc.CreateOrganizationParams{
		Name:    payload.Name,
		OwnerID: userId,
	}

	organization, err := c.store.CreateOrganization(params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// also make a membership for the owner

	utils.WriteJSON(w, http.StatusCreated, organization)
}

func (c *Handler) HandleGetOrganizationById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	organization, err := c.store.GetOrganizationById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}