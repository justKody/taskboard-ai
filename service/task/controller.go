package task

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/justKody/taskboard-go-api/db/sqlc"
	"github.com/justKody/taskboard-go-api/middleware"
	"github.com/justKody/taskboard-go-api/utils"
)

// only managers and above assigned to the project's organization may act on tasks
func canManageTasks(role string) bool {
	return role == string(sqlc.MembershipsRoleManager) ||
		role == string(sqlc.MembershipsRoleAdmin) ||
		role == string(sqlc.MembershipsRoleSuperAdmin)
}

// manager can create task for assigned project
func (h *Handler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["id"]

	var payload CreateTaskRequestDTO
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

	project, err := h.projectStore.GetProject(r.Context(), projectId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if project == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("Project not found"))
		return
	}

	membership, err := h.membershipStore.GetMembershipByUserAndOrganization(r.Context(), userId, project.OrganizationID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if membership == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to create a task in this project"))
		return
	}
	if !canManageTasks(membership.Role) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Your role is not authorized to create a task in this project"))
		return
	}

	var assignedTo pgtype.UUID
	if err := assignedTo.Scan(payload.AssignedTo); err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("Invalid assigned_to user id"))
		return
	}

	var dueDate pgtype.Date
	if err := dueDate.Scan(payload.DueDate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("Invalid due_date"))
		return
	}

	params := sqlc.CreateTaskParams{
		ProjectID:   projectId,
		Title:       payload.Title,
		Description: pgtype.Text{String: payload.Description, Valid: payload.Description != ""},
		Status:      sqlc.TaskStatus(payload.Status),
		Priority:    sqlc.Priority(payload.Priority),
		DueDate:     dueDate,
		AssignedTo:  assignedTo,
		CreatedBy:   userId,
	}

	task, err := h.store.CreateTask(r.Context(), params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, task)
}

// all given task to him?
func (h *Handler) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["id"]

	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authenticated"))
		return
	}

	project, err := h.projectStore.GetProject(r.Context(), projectId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if project == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("Project not found"))
		return
	}

	membership, err := h.membershipStore.GetMembershipByUserAndOrganization(r.Context(), userId, project.OrganizationID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if membership == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to view tasks in this project"))
		return
	}

	var assignedTo pgtype.UUID
	if err := assignedTo.Scan(userId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tasks, err := h.store.GetTask(r.Context(), projectId, assignedTo)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tasks)
}

// all tasks in project
func (h *Handler) HandleListProjectTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["id"]

	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authenticated"))
		return
	}

	project, err := h.projectStore.GetProject(r.Context(), projectId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if project == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("Project not found"))
		return
	}

	membership, err := h.membershipStore.GetMembershipByUserAndOrganization(r.Context(), userId, project.OrganizationID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if membership == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to view tasks in this project"))
		return
	}

	tasks, err := h.store.ListProjectTask(r.Context(), projectId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tasks)
}

// manager can update a task in an assigned project
func (h *Handler) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["id"]
	taskId := vars["taskId"]

	var payload UpdateTaskRequestDTO
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

	project, err := h.projectStore.GetProject(r.Context(), projectId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if project == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("Project not found"))
		return
	}

	membership, err := h.membershipStore.GetMembershipByUserAndOrganization(r.Context(), userId, project.OrganizationID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if membership == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to update a task in this project"))
		return
	}
	if !canManageTasks(membership.Role) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Your role is not authorized to update a task in this project"))
		return
	}

	var assignedTo pgtype.UUID
	if err := assignedTo.Scan(payload.AssignedTo); err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("Invalid assigned_to user id"))
		return
	}

	var dueDate pgtype.Date
	if err := dueDate.Scan(payload.DueDate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("Invalid due_date"))
		return
	}

	params := sqlc.UpdateTaskParams{
		ID:          taskId,
		ProjectID:   projectId,
		Title:       payload.Title,
		Description: pgtype.Text{String: payload.Description, Valid: payload.Description != ""},
		Status:      sqlc.TaskStatus(payload.Status),
		Priority:    sqlc.Priority(payload.Priority),
		DueDate:     dueDate,
		AssignedTo:  assignedTo,
	}
	task, err := h.store.UpdateTask(r.Context(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.WriteError(w, http.StatusNotFound, errors.New("Task not found"))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, task)
}

// delete the task
func (h *Handler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["id"]
	taskId := vars["taskId"]

	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authenticated"))
		return
	}

	project, err := h.projectStore.GetProject(r.Context(), projectId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if project == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("Project not found"))
		return
	}

	membership, err := h.membershipStore.GetMembershipByUserAndOrganization(r.Context(), userId, project.OrganizationID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if membership == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Not authorized to delete a task in this project"))
		return
	}
	if !canManageTasks(membership.Role) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Your role is not authorized to delete a task in this project"))
		return
	}

	if err := h.store.DeleteTask(r.Context(), taskId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Task deleted successfully")
}
