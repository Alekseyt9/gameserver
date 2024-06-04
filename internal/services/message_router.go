package services

import (
	"gameserver/internal/services/model"
)

type MessageRouter struct {
	gameManager   *GameManager
	roomManager   *RoomManager
	playerManager *PlayerManager
}

func (r *MessageRouter) Route(msg string) {
	cmsg := toClientMessage(msg)
	switch cmsg.MessageType {
	case GameMsgType:
		r.gameManager.Process(cmsg)
	case RoomMsgType:
		r.roomManager.Process(cmsg)
	case UserMsgType:
		r.playerManager.Process(cmsg)
	}
}

// частичный парсинг сообщения с клиента
func toClientMessage(msg string) *model.ClientMsg {
	return nil
}
