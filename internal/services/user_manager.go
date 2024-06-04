package services

import "gameserver/internal/services/model"

type UserManager struct {
}

/*
Создает нового пользователя (анонимного)
*/
func (*UserManager) CreateUser() *model.User {
	return nil
}
