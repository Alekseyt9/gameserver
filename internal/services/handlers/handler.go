package handlers

import "gameserver/internal/services/store"

type Handler struct {
	store store.Store
}

func Create(store store.Store) *Handler {
	return &Handler{
		store: store,
	}
}
