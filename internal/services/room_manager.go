package services

import (
	"gameserver/internal/services/model"

	"github.com/beevik/guid"
)

type RoomManager struct {
}

type RoomMsg struct {
	Action string
}

type RoomInfo struct {
}

const (
	RoomMsgType   = "matching"
	ActionConnect = "connect"
	ActionQuit    = "quit"
)

func (r *RoomManager) Process(msg *model.ClientMsg) {
	rmsg := r.parseRoomMsg(msg.Data)

	switch rmsg.Action {
	case ActionConnect:
		r.playerConnect(msg)
	case ActionQuit:
		r.playerQuit(msg)
	}
}

// подключение к существующей комнате/восстановление, если уже в комнате/создание комнаты
func (r *RoomManager) playerConnect(msg *model.ClientMsg) {
	room := r.getExistingRoom(msg.PlayerID, msg.GameType)
	if room != nil {
		// комната уже есть, стартуем игру
	} else {
		// нужно найти комнату в ожидании или создать новую
		// потом старт игры
	}
}

func (r *RoomManager) playerQuit(msg *model.ClientMsg) {

}

func (r *RoomManager) parseRoomMsg(msg string) *RoomMsg {
	return nil
}

func (r *RoomManager) getExistingRoom(userId guid.Guid, gameType string) *RoomInfo {
	return nil
}

/*
Получает существующую комнату для игрока или создает новую
Существующую: если игра с таким типом для игрока уже есть или игры нет, но игрок подключается к свободной комнате
Новую: если нет свободной комнаты для игрока и нет существующей игры игры
*/
func (r *RoomManager) getOrCreateRoom(userId guid.Guid, gameType string) {
	// запрос в хранилище по игроку и типу игры
	// ? нужна ли комната в базе, котоорая ожидает игроков?
	// наверное да, тк чтобы игроки могли подключаться, когда другие офлайн
}
