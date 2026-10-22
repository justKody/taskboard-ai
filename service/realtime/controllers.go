package realtime

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/justKody/taskboard-go-api/middleware"
	"github.com/justKody/taskboard-go-api/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		// allowing no origin header and any localhost dev origin
		return origin == "" || strings.HasPrefix(origin, "http://localhost")
	},
}

func (h *Handler) HandleWs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("not authenticated"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return // Upgrade already wrote the HTTP error
	}

	client := newClient(h.hub, conn, userID)
	h.hub.register <- client

	go client.writePump()
	client.readPump()
}
