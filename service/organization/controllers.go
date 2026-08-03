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
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authenticated"))
		return
	}

	params := sqlc.CreateOrganizationParams{
		Name:    payload.Name,
		OwnerID: userId,
	}

	// create the organization
	organization, err := c.store.CreateOrganization(r.Context(), params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// create the membership for the owner
	var createMembershipParams = sqlc.CreateMembershipParams{
		UserID:         userId,
		OrganizationID: organization.Id,
		Role:           sqlc.MembershipsRoleSuperAdmin,
	}

	_, err = c.membershipStore.CreateMembership(r.Context(), createMembershipParams)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, organization)
}

func (c *Handler) HandleGetOrganizationDetailsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// get the organization details
	organization, err := c.store.GetOrganizationById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if organization == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("organization not found"))
		return
	}

	// get the memberships for the organization
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

	isOwner, err := c.store.CheckIfUserIsOwner(r.Context(), userId, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if !isOwner {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to change owner of this organization"))
		return
	}

	// change the owner of the organization
	err = c.store.ChangeOwnerOfOrganizationById(r.Context(), payload.NewOwnerID, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// also update the membership of the new owner (owner is tracked on orgs; memberships use super_admin)
	err = c.membershipStore.UpdateMembershipRole(r.Context(), payload.NewOwnerID, id, sqlc.MembershipsRoleSuperAdmin)
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

	owner, err := c.store.CheckIfUserIsOwner(r.Context(), userId, id)
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
	err = c.store.DeleteOrganizationById(r.Context(), id)
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

	organizations, err := c.store.GetOrganizationsListByUserId(r.Context(), userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, organizations)
}

func (c *Handler) HandleInviteUserToOrganization(w http.ResponseWriter, r *http.Request) {
	var payload InviteUserToOrganizationRequestDTO
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// get the inviter's user id
	invitedBy, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authenticated"))
		return
	}

	if invitedBy == payload.UserID {
		utils.WriteError(w, http.StatusBadRequest, errors.New("cannot invite yourself"))
		return
	}

	// check if the organization exists
	organization, err := c.store.GetOrganizationById(r.Context(), payload.OrganizationID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if organization == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("organization not found"))
		return
	}

	// check if the inviter is a super_admin or admin
	inviterMembership, err := c.membershipStore.GetMembershipByUserAndOrganization(r.Context(), invitedBy, payload.OrganizationID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if inviterMembership == nil ||
		(inviterMembership.Role != string(sqlc.MembershipsRoleSuperAdmin) &&
			inviterMembership.Role != string(sqlc.MembershipsRoleAdmin)) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to invite users to this organization"))
		return
	}

	// check if the invitee exists
	invitee, err := c.userStore.GetUserById(r.Context(), payload.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if invitee == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("user not found"))
		return
	}

	// check if the invitee is already a member of the organization
	existingMembership, err := c.membershipStore.GetMembershipByUserAndOrganization(r.Context(), payload.UserID, payload.OrganizationID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if existingMembership != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("user is already a member of this organization"))
		return
	}
	// check if the invitee has a pending invite
	pendingInvite, err := c.store.GetPendingOrganizationInvite(r.Context(), payload.OrganizationID, payload.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if pendingInvite != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("user already has a pending invite"))
		return
	}

	invite, err := c.store.CreateOrganizationInvite(r.Context(), sqlc.CreateOrganizationInviteParams{
		OrganizationID: payload.OrganizationID,
		UserID:         payload.UserID,
		InvitedBy:      invitedBy,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, invite)
}

func (c *Handler) HandleAcceptOrRejectOrganizationInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inviteId := vars["id"]

	var payload AcceptOrRejectOrganizationInviteRequestDTO
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

	invite, err := c.store.GetOrganizationInviteById(r.Context(), inviteId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if invite == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("invite not found"))
		return
	}

	if invite.UserId != userId {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to respond to this invite"))
		return
	}

	if invite.Status != string(sqlc.OrganizationInviteStatusPending) {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invite is no longer pending"))
		return
	}

	status := sqlc.OrganizationInviteStatus(payload.Status)

	switch status {
	case sqlc.OrganizationInviteStatusAccepted:
		existingMembership, err := c.membershipStore.GetMembershipByUserAndOrganization(r.Context(), userId, invite.OrganizationId)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if existingMembership != nil {
			utils.WriteError(w, http.StatusBadRequest, errors.New("user is already a member of this organization"))
			return
		}

		_, err = c.membershipStore.CreateMembership(r.Context(), sqlc.CreateMembershipParams{
			UserID:         userId,
			OrganizationID: invite.OrganizationId,
			Role:           sqlc.MembershipsRoleMember,
		})
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	case sqlc.OrganizationInviteStatusRejected:
	default:
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid invite status"))
		return
	}

	updatedInvite, err := c.store.UpdateOrganizationInviteStatus(r.Context(), inviteId, status)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if updatedInvite == nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invite is no longer pending"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedInvite)
}
