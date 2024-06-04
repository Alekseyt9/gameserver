package services

import "gameserver/internal/services/model"

const UserMsgType = "register"

type PlayerManager struct {
}

func (g *PlayerManager) Process(msg *model.ClientMsg) error {
	return nil
}

/*
Создает нового пользователя (анонимного)
*/
func (*PlayerManager) createUser() *model.Player {
	return nil
}
