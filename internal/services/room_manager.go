package services

import (
	"context"
	"errors"
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
	"log/slog"
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
	log           *slog.Logger
}

type PlayerConnectResult struct {
	State       string
	ContentLink string
	RoomID      *uuid.UUID
}

const (
	chanBuffer = 100
)

func NewRoomManager(store store.Store, gm *GameManager, pm *PlayerManager, m *Matcher, log *slog.Logger) *RoomManager {
	return &RoomManager{
		store:         store,
		gameManager:   gm,
		chanMap:       make(map[uuid.UUID]chan model.GameMsg, chanBuffer),
		matcher:       m,
		playerManager: pm,
		log:           log,
	}
}

func (m *RoomManager) GetOrCreateChan(roomID uuid.UUID) chan model.GameMsg {
	m.chanLock.Lock()
	defer m.chanLock.Unlock()

	ch, ok := m.chanMap[roomID]
	if !ok {
		ch = make(chan model.GameMsg, chanBuffer)
		m.chanMap[roomID] = ch

		// создаем воркер для комнаты.
		go func() {
			ctx := context.Background()
			for msg := range ch {
				gctx, err := m.gameManager.Process(ctx, &msg)
				if err != nil {
					m.log.Error("m.gameManager.Process error", err)
				}

				err = m.processResult(gctx)
				if err != nil {
					m.log.Error("m.processResult error", err)
				}
			}
		}()
	}

	return ch
}

// сохранение стейта игры, рассылка сообщений.
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

// подключение к существующей комнате или создание комнаты.
// подключаться нужно каждый раз при коннекте игрока.
func (m *RoomManager) PlayerConnect(
	ctx context.Context,
	playerID uuid.UUID,
	gameID string) (*PlayerConnectResult, error) {
	room, err := m.GetExistingRoom(ctx, gameID, playerID)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}

	if room != nil {
		switch room.Status {
		case "game": // комната в игре есть - стартуем игру.
			return &PlayerConnectResult{
				State:       "game",
				ContentLink: m.gameManager.GetGameInfo(gameID).ContentURL,
				RoomID:      &room.ID,
			}, nil
		case "wait": // комната есть, но игрок ожидает игры
			return &PlayerConnectResult{
				State:  "wait",
				RoomID: &room.ID,
			}, nil
		}
	}

	// ставим в очередь в матчмейкинг.
	m.matcher.CheckAndAdd(model.RoomQuery{
		PlayerID: playerID,
		GameID:   gameID,
	})
	return &PlayerConnectResult{State: "wait"}, nil
}

// выход из комнаты (выйти можно только один раз).
func (m *RoomManager) PlayerQuit(ctx context.Context, gameID string, playerID uuid.UUID) error {
	room, err := m.GetExistingRoom(ctx, gameID, playerID)
	if err != nil {
		return err
	}
	ch := m.GetOrCreateChan(room.ID)

	// чтобы стейт игры не перезаписывался - в игру событие передаем через канал.
	ch <- m.createQuitGameMsg(gameID, playerID)

	// помечаем игроков, которые вышли.
	err = m.store.MarkDropRoomPlayer(ctx, room.ID, playerID)
	if err != nil {
		return err
	}

	return nil
}

func (m *RoomManager) createQuitGameMsg(gameID string, playerID uuid.UUID) model.GameMsg {
	return model.GameMsg{
		Type:     "game",
		GameID:   gameID,
		PlayerID: playerID,
		Data:     map[string]interface{}{"action": "quit"},
	}
}

func (m *RoomManager) GetExistingRoom(ctx context.Context, gameID string, playerID uuid.UUID) (*model.Room, error) {
	r, err := m.store.GetRoom(ctx, gameID, playerID, false)
	if err != nil {
		return nil, err
	}
	return r, nil
}
