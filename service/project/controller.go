package project

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/justKody/taskboard-go-api/db/sqlc"
	"github.com/justKody/taskboard-go-api/middleware"
	"github.com/justKody/taskboard-go-api/utils"
)

func (h *Handler) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	var payload CreateProjectRequestDTO
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

	// only members of the organization can create projects in it
	membership, err := h.membershipStore.GetMembershipByUserAndOrganization(r.Context(), userId, payload.OrganizationID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if membership == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to create a project in this organization"))
		return
	}

	// only with role greater than admin can create a project
	if membership.Role != string(sqlc.MembershipsRoleAdmin) || membership.Role != string(sqlc.MembershipsRoleSuperAdmin) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Your role is not authorized to create project for this organization"))
	}

	params := sqlc.CreateProjectParams{
		OrganizationID: payload.OrganizationID,
		Name:           payload.Name,
		Description:    pgtype.Text{String: payload.Description, Valid: payload.Description != ""},
		CreatedBy:      userId,
	}

	project, err := h.store.CreateProject(r.Context(), params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, project)
}

func (h *Handler) handleListProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	organizationId := vars["id"]

	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authenticated"))
		return
	}

	// only members of the organization can see the projects
	membership, err := h.membershipStore.GetMembershipByUserAndOrganization(r.Context(), userId, organizationId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if membership == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to create a project in this organization"))
		return
	}

	projects, err := h.store.ListProject(r.Context(), organizationId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, projects)

}

func (h *Handler) handleUpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	organizationId := vars["id"]
	projectId := vars["projectId"]

	var payload UpdateProjectRequestDTO
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

	membership, err := h.membershipStore.GetMembershipByUserAndOrganization(r.Context(), userId, organizationId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if membership == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to update a project in this organization"))
		return
	}

	if membership.Role != string(sqlc.MembershipsRoleAdmin) && membership.Role != string(sqlc.MembershipsRoleSuperAdmin) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Your role is not authorized to update a project for this organization"))
		return
	}

	params := sqlc.UpdateProjectParams{
		ID:             projectId,
		OrganizationID: organizationId,
		Name:           payload.Name,
		Description:    pgtype.Text{String: payload.Description, Valid: payload.Description != ""},
		Status:         sqlc.ProjectStatus(payload.Status),
	}

	project, err := h.store.UpdateProject(r.Context(), params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if project == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("project not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, project)
}

func (h *Handler) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	organizationId := vars["id"]
	projectId := vars["projectId"]

	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authenticated"))
		return
	}

	membership, err := h.membershipStore.GetMembershipByUserAndOrganization(r.Context(), userId, organizationId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if membership == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to delete a project in this organization"))
		return
	}

	if membership.Role != string(sqlc.MembershipsRoleAdmin) && membership.Role != string(sqlc.MembershipsRoleSuperAdmin) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Your role is not authorized to delete a project for this organization"))
		return
	}

	err = h.store.DeleteProject(r.Context(), projectId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Project deleted successfully")
}
