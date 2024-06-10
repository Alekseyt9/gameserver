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

const GameMsgType = "game"

func NewGameManager(store store.Store, pm *PlayerManager) *GameManager {
	man := &GameManager{
		procs:         make(map[string]model.GameProcessor),
		store:         store,
		playerManager: pm,
	}
	man.procs["tictactoe"] = &game.TTCProcessor{}

	return man
}

// поднимает state для пользователя+тип игры, передает в обработчик
func (m *GameManager) Process(ctx context.Context, msg *model.GameMsg) error {
	if msg.MessageType != GameMsgType {
		return errors.New("wrong MessageType")
	}

	room, err := m.store.GetRoom(ctx, msg.PlayerID, msg.GameID)
	if err != nil {
		return err
	}

	proc, err := m.getProcessor(msg.GameID)
	if err != nil {
		return err
	}

	procCtx := NewGameProcessorCtx(m.store, m.playerManager, room.ID, msg.GameID)
	proc.Process(procCtx, room.State, msg)

	return nil
}

func (m *GameManager) getProcessor(gameType string) (model.GameProcessor, error) {
	v, ok := m.procs[gameType]
	if !ok {
		return nil, errors.New("wrong processr")
	}
	return v, nil
}
