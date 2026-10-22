package realtime

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justKody/taskboard-go-api/middleware"
)

type Handler struct {
	hub *Hub
}

func NewHandler() *Handler {
	hub := newHub()
	go hub.run()
	return &Handler{hub: hub}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	wsRouter := router.PathPrefix("/ws").Subrouter()

	wsRouter.Use(middleware.Auth)

	wsRouter.HandleFunc("/", h.HandleWs).Methods(http.MethodGet)
}
