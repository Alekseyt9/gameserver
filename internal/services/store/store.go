package store

import (
	"gameserver/internal/services/model"

	"github.com/beevik/guid"
)

type Store interface {
	GetUser(id guid.Guid) *model.Player
	CreateUser(*model.Player)

	//GetRoomState(id guid.Guid) string
	GetRoomState(playerID guid.Guid, gameType string) string
	SetRoomState(id guid.Guid, state string)

	CreateRoom()
	DropRoom(id guid.Guid)
}
