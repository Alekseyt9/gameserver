package store

import (
	"context"
	"errors"
	"gameserver/internal/services/model"

	"github.com/beevik/guid"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	CreatePlayer(ctx context.Context, player *model.Player) error

	GetRoom(ctx context.Context, gameID string, playerID guid.Guid) (*model.Room, error)
	SetRoomState(ctx context.Context, id guid.Guid, state string) error
	DropRoomPlayer(ctx context.Context, roomID guid.Guid, playerID guid.Guid) error
	CreateOrUpdateRooms(ctx context.Context, rooms []*model.MatcherRoom) error
	LoadWaitingRooms(ctx context.Context) ([]*model.MatcherRoom, error)
}
