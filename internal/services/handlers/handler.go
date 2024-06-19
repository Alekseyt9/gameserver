package handlers

import (
	"gameserver/internal/services"
	"gameserver/internal/services/store"
	"log/slog"
)

type Handler struct {
	store         store.Store
	roomManager   *services.RoomManager
	log           *slog.Logger
	wsManager     *services.WebSocketManager
	playerManager *services.PlayerManager
}

func New(store store.Store, rm *services.RoomManager, pm *services.PlayerManager, ws *services.WebSocketManager, log *slog.Logger) *Handler {
	return &Handler{
		store:         store,
		roomManager:   rm,
		log:           log,
		wsManager:     ws,
		playerManager: pm,
	}
}
