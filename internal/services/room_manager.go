package services

import (
	"context"
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
	"sync"

	"github.com/google/uuid"
)

type RoomManager struct {
	store         store.Store
	gameManager   *GameManager
	playerManager *PlayerManager
	matcher       *Matcher
	chanMap       map[uuid.UUID]chan model.GameMsg
	chanLock      sync.RWMutex
}

type PlayerConnectResult struct {
	State string
}

func NewRoomManager(store store.Store, gm *GameManager, pm *PlayerManager, m *Matcher) *RoomManager {
	return &RoomManager{
		store:         store,
		gameManager:   gm,
		chanMap:       make(map[uuid.UUID]chan model.GameMsg, 100),
		matcher:       m,
		playerManager: pm,
	}
}

func (m *RoomManager) GetOrCreateChan(roomID uuid.UUID) chan model.GameMsg {
	m.chanLock.Lock()
	defer m.chanLock.Unlock()

	ch, ok := m.chanMap[roomID]
	if !ok {
		ch = make(chan model.GameMsg, 100)
		m.chanMap[roomID] = ch

		// создаем воркер для комнаты
		go func() {
			ctx := context.Background()
			for msg := range ch {
				gctx, err := m.gameManager.Process(ctx, &msg)
				if err != nil {
					//TODO log
				}
				m.processResult(gctx)
			}
		}()
	}

	return ch
}

// сохранение стейта игры, рассылка сообщений
func (m *RoomManager) processResult(gctx *GameProcessorCtx) error {
	ctx := context.Background()

	if gctx.gameState != "" {
		err := m.store.SetRoomState(ctx, gctx.roomID, gctx.gameState)
		if err != nil {
			return err
		}
	}

	for _, msg := range gctx.sendMessages {
		chp := m.playerManager.GetChan(msg.PlayerID)
		if chp != nil {
			*chp <- msg
		}
	}

	return nil
}

func (m *RoomManager) DeleteChan(roomID uuid.UUID) {
	m.chanLock.Lock()
	defer m.chanLock.Unlock()

	ch, ok := m.chanMap[roomID]
	if ok {
		close(ch)
		delete(m.chanMap, roomID)
	}
}

// подключение к существующей комнате или создание комнаты
// подключаться нужно каждый раз при коннекте игрока
func (m *RoomManager) PlayerConnect(ctx context.Context, playerID uuid.UUID, gameID string) (*PlayerConnectResult, error) {
	room, err := m.GetExistingRoom(ctx, gameID, playerID)
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
func (m *RoomManager) PlayerQuit(ctx context.Context, gameID string, playerID uuid.UUID) error {
	room, err := m.GetExistingRoom(ctx, gameID, playerID)
	if err != nil {
		return err
	}
	ch := m.GetOrCreateChan(room.ID)

	// чтобы стейт игры не перезаписывался - в игру событие передаем через канал
	ch <- m.createQuitGameMsg(gameID, playerID)

	// помечаем игроков, которые вышли
	err = m.store.MarkDropRoomPlayer(ctx, room.ID, playerID)
	if err != nil {
		return err
	}

	return nil
}

func (m *RoomManager) createQuitGameMsg(gameID string, playerId uuid.UUID) model.GameMsg {
	return model.GameMsg{
		Type:     "game",
		GameID:   gameID,
		PlayerID: playerId,
		Data:     map[string]interface{}{"action": "quit"},
	}
}

func (m *RoomManager) GetExistingRoom(ctx context.Context, gameID string, playerID uuid.UUID) (*model.Room, error) {
	r, err := m.store.GetRoom(ctx, gameID, playerID)
	if err != nil {
		return nil, err
	}
	return r, nil
}
