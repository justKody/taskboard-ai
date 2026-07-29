package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justKody/taskboard-go-api/middleware"
)

type Handler struct {
	store UserStore
}

func NewHandler(store UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.HandleLogin).Methods(http.MethodPost)
	router.HandleFunc("/logout", h.HandleLogout).Methods(http.MethodPost)
	router.HandleFunc("/signup", h.HandleSignup).Methods(http.MethodPost)

	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.Use(middleware.Auth)
	userRouter.HandleFunc("/me", h.HandleGetMe).Methods(http.MethodGet)
	userRouter.HandleFunc("/update", h.HandleUpdateUser).Methods(http.MethodPut)
	userRouter.HandleFunc("/change/password", h.HandleChangePassword).Methods(http.MethodPut)
	userRouter.HandleFunc("/delete", h.HandleDeleteUser).Methods(http.MethodDelete)
	userRouter.HandleFunc("/", h.HandleGetUsersList).Methods(http.MethodGet)
}
