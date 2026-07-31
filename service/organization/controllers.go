package organization

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
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

func (c *Handler) HandleGetOrganizationDetailsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	organization, err := c.store.GetOrganizationById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if organization == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("organization not found"))
		return
	}

	members, err := c.membershipStore.GetAllMembershipsByOrganizationId(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, OrganizationDetailsResponseDTO{
		Organization: organization,
		Members:      members,
	})
}

func (c *Handler) HandleChangeOwnerOfOrganizationById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var payload ChangeOwnerOfOrganizationRequestDTO
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authenticated"))
		return
	}

	// check if the user is the owner of the organization

	isOwner, err := c.store.CheckIfUserIsOwner(userId, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if !isOwner {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to change owner of this organization"))
		return
	}

	// change the owner of the organization
	err = c.store.ChangeOwnerOfOrganizationById(payload.NewOwnerID, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// also update the membership of the new owner (owner is tracked on orgs; memberships use admin)
	err = c.membershipStore.UpdateMembershipRole(r.Context(), payload.NewOwnerID, id, sqlc.MembershipsRoleAdmin)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Owner changed successfully")
}

func (c *Handler) HandleDeleteOrganizationById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// check 1: if the user is the owner of the organization, then we can delete the organization

	userId, _ := middleware.GetUserID(r.Context())

	owner, err := c.store.CheckIfUserIsOwner(userId, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if !owner {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to delete this organization"))
		return
	}

	// checks 2: If there are already members in the organization, then we cannot delete the organization
	var areMembers = false
	memberships, err := c.membershipStore.GetAllMembershipsByOrganizationId(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			areMembers = false
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}
	if len(memberships) > 0 {
		areMembers = true
	}

	if areMembers {
		utils.WriteError(w, http.StatusBadRequest, errors.New("organization has members cannot be deleted"))
		return
	}

	// delete the organization
	err = c.store.DeleteOrganizationById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Organization deleted successfully")
}

func (c *Handler) HandleGetOrganizationsListByUserId(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authenticated"))
		return
	}

	organizations, err := c.store.GetOrganizationsListByUserId(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, organizations)
}
