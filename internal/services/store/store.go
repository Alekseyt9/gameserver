package store

import (
	"gameserver/internal/services/model"

	"github.com/beevik/guid"
)

type Store interface {
	GetUser(id guid.Guid) (*model.Player, error)
	CreateUser(*model.Player) error

	//GetRoomState(id guid.Guid) string
	GetRoomState(playerID guid.Guid, gameType string) (string, error)
	SetRoomState(id guid.Guid, state string) error

	CreateRoom() error
	DropRoom(id guid.Guid) error
}
