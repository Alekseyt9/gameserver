package handlers

import (
	"gameserver/internal/services"
	"gameserver/internal/services/store"
)

type Handler struct {
	store       store.Store
	roomManager *services.RoomManager
}

func New(store store.Store, rm *services.RoomManager) *Handler {
	return &Handler{
		store:       store,
		roomManager: rm,
	}
}
