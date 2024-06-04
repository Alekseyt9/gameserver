package store

import (
	"gameserver/internal/services/model"

	"github.com/beevik/guid"
)

type Store interface {
	GetUser(id guid.Guid) *model.User
	CreateUser(*model.User)

	GetRoomState(id guid.Guid) string
	SetRoomState(id guid.Guid, state string)

	CreateRoom()
	DropRoom(id guid.Guid)
}
