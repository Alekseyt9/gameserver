package handlers

import (
	"gameserver/internal/services"
	"gameserver/internal/services/store"
	"log/slog"
)

type Handler struct {
	store       store.Store
	roomManager *services.RoomManager
	log         *slog.Logger
}

func New(store store.Store, rm *services.RoomManager, log *slog.Logger) *Handler {
	return &Handler{
		store:       store,
		roomManager: rm,
		log:         log,
	}
}
