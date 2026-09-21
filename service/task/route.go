package task

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justKody/taskboard-go-api/middleware"
	"github.com/justKody/taskboard-go-api/service/membership"
	"github.com/justKody/taskboard-go-api/service/project"
)

type Handler struct {
	store           TaskStore
	projectStore    project.ProjectStore
	membershipStore membership.MemebershipStore
}

func NewHandler(store TaskStore, projectStore project.ProjectStore, membershipStore membership.MemebershipStore) *Handler {
	return &Handler{
		store:           store,
		projectStore:    projectStore,
		membershipStore: membershipStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	taskRouter := router.PathPrefix("/task").Subrouter()

	taskRouter.Use(middleware.Auth)

	taskRouter.HandleFunc("/create/{id}", h.HandleCreateTask).Methods(http.MethodPost)
	taskRouter.HandleFunc("/get/{id}", h.HandleGetTask).Methods(http.MethodGet)
	taskRouter.HandleFunc("/update/{id}/{taskId}", h.HandleUpdateTask).Methods(http.MethodPut)
	taskRouter.HandleFunc("/delete/{id}/{taskId}", h.HandleDeleteTask).Methods(http.MethodDelete)
	taskRouter.HandleFunc("/list/{id}", h.HandleListProjectTask).Methods(http.MethodGet)
}
