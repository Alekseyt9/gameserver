package services

import (
	"context"
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"

	"github.com/beevik/guid"
)

type RoomManager struct {
	store       store.Store
	gameManager *GameManager
	matcher     *Matcher
	chanMap     map[guid.Guid]chan model.GameMsg
}

type PlayerConnectResult struct {
	State string
}

func NewRoomManager(store store.Store, gm *GameManager, m *Matcher) *RoomManager {
	return &RoomManager{
		store:       store,
		gameManager: gm,
		chanMap:     make(map[guid.Guid]chan model.GameMsg, 100),
		matcher:     m,
	}
}

func (m *RoomManager) GetOrCreateChan(roomID guid.Guid) chan model.GameMsg {
	ch, ok := m.chanMap[roomID]
	if !ok {
		ch = make(chan model.GameMsg, 100)
		m.chanMap[roomID] = ch

		// создаем воркер для комнаты
		go func() {
			ctx := context.Background()
			for msg := range ch {
				m.gameManager.Process(ctx, &msg)
			}
		}()
	}

	return ch
}

func (m *RoomManager) DeleteChan(roomID guid.Guid) {
	ch, ok := m.chanMap[roomID]
	if ok {
		close(ch)
		delete(m.chanMap, roomID)
	}
}

// подключение к существующей комнате или создание комнаты
// подключаться нужно каждый раз при коннекте игрока
func (m *RoomManager) PlayerConnect(ctx context.Context, playerID guid.Guid, gameID string) (*PlayerConnectResult, error) {
	room, err := m.GetExistingRoom(ctx, playerID, gameID)
	if err != nil {
		return nil, err
	}

	if room != nil && room.Status == "game" {
		// комната в игре есть - стартуем игру
		return &PlayerConnectResult{State: "game"}, nil
	} else {
		// ставим в очередь в матчмейкинг
		m.matcher.CheckAndAdd(model.RoomQuery{
			PlayerID: playerID,
			GameID:   gameID,
		})
		return &PlayerConnectResult{State: "wait"}, nil
	}
}

// выход из комнаты (выйти можно только один раз)
func (m *RoomManager) PlayerQuit(msg *model.GameMsg) {

}

func (m *RoomManager) GetExistingRoom(ctx context.Context, playerID guid.Guid, gameID string) (*model.Room, error) {
	r, err := m.store.GetRoom(ctx, playerID, gameID)
	if err != nil {
		return nil, err
	}
	return r, nil
}
