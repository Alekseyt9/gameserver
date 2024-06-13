package store

import (
	"context"
	"errors"
	"gameserver/internal/services/model"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	CreatePlayer(ctx context.Context, player *model.Player) error
	GetPlayer(ctx context.Context, playerID uuid.UUID) (*model.Player, error)

	GetRoom(ctx context.Context, gameID string, playerID uuid.UUID) (*model.Room, error)
	SetRoomState(ctx context.Context, id uuid.UUID, state string) error

	// помечаем игрока комнаты на удаление
	MarkDropRoomPlayer(ctx context.Context, roomID uuid.UUID, playerID uuid.UUID) error

	// создаем/обновляем комнаты после матчинга
	CreateOrUpdateRooms(ctx context.Context, rooms []*model.MatcherRoom) error

	// загружаем комнаты в ожидании игроков
	LoadWaitingRooms(ctx context.Context) ([]*model.MatcherRoom, error)
}
