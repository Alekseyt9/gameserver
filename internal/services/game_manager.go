package services

import (
	"context"
	"errors"
	game "gameserver/internal/services/game/tictactoe"
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
)

type GameManager struct {
	procs         map[string]model.GameProcessor
	store         store.Store
	playerManager *PlayerManager
}

func NewGameManager(store store.Store, pm *PlayerManager) *GameManager {
	man := &GameManager{
		procs:         make(map[string]model.GameProcessor),
		store:         store,
		playerManager: pm,
	}
	man.procs["tictactoe"] = &game.TTCProcessor{}

	return man
}

func (m *GameManager) GetGameInfo(gameID string) *model.GameInfo {
	p, ok := m.procs[gameID]
	if ok {
		return p.GetInfo()
	}
	return nil
}

// поднимает state для пользователя+тип игры, передает в обработчик
func (m *GameManager) Process(ctx context.Context, msg *model.GameMsg) (*GameProcessorCtx, error) {
	if msg.MessageType != "game" {
		return nil, errors.New("wrong MessageType")
	}

	room, err := m.store.GetRoom(ctx, msg.GameID, msg.PlayerID)
	if err != nil {
		return nil, err
	}

	proc, err := m.getProcessor(msg.GameID)
	if err != nil {
		return nil, err
	}

	procCtx := NewGameProcessorCtx(room.ID, msg.GameID)
	proc.Process(procCtx, room.State, msg)

	return procCtx, nil
}

// инициализация игры
func (m *GameManager) Init(gameID string, players []model.MatcherPlayer) (string, error) {
	proc, err := m.getProcessor(gameID)
	if err != nil {
		return "", err
	}
	state, err := proc.Init(players)
	return state, err
}

func (m *GameManager) getProcessor(gameType string) (model.GameProcessor, error) {
	v, ok := m.procs[gameType]
	if !ok {
		return nil, errors.New("wrong processr")
	}
	return v, nil
}
