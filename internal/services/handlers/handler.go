package handlers

import "gameserver/internal/services/store"

type Handler struct {
	store store.Store
}

func New(store store.Store) *Handler {
	return &Handler{
		store: store,
	}
}
