package store

import (
	"context"
	"errors"
	"gameserver/internal/services/model"

	"github.com/beevik/guid"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	GetPlayer(ctx context.Context, id guid.Guid) (*model.Player, error)
	CreatePlayer(ctx context.Context, player *model.Player) error

	GetRoom(ctx context.Context, playerID guid.Guid, gameID string) (*model.Room, error)
	SetRoomState(ctx context.Context, id guid.Guid, state string) error

	//CreateRoom(ctx context.Context) error
	//DropRoom(ctx context.Context, id guid.Guid) error

	CreateOrUpdateRooms(ctx context.Context, rooms []*model.MatcherRoom) error
	LoadWaitingRooms(ctx context.Context) ([]*model.MatcherRoom, error)
}
