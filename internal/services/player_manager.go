package services

import (
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
)

type PlayerManager struct {
	store store.Store
}

/*
Создает нового пользователя (анонимного)
*/
func (*PlayerManager) CreateUser() *model.Player {
	return nil
}
