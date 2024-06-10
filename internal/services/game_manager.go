package services

import (
	"context"
	"errors"
	game "gameserver/internal/services/game/tictactoe"
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
)

type GameManager struct {
	procs map[string]model.GameProcessor
	store store.Store
}

const GameMsgType = "game"

func CreareGameManager(store store.Store) *GameManager {
	man := &GameManager{
		procs: make(map[string]model.GameProcessor),
		store: store,
	}
	man.procs["tictactoe"] = &game.TTCProcessor{}

	return man
}

func (g *GameManager) Process(ctx context.Context, msg *model.GameMsg) error {
	// поднимает state для пользователя+тип игры
	// передает в обработчик

	if msg.MessageType != GameMsgType {
		return errors.New("wrong MessageType")
	}

	room, err := g.store.GetRoom(ctx, msg.PlayerID, msg.GameID)
	if err != nil {
		return err
	}

	proc, err := g.getProcessor(msg.GameID)
	if err != nil {
		return err
	}

	procCtx := &model.GameProcessorCtx{}
	proc.Process(procCtx, room.State, msg)

	return nil
}

func (g *GameManager) getProcessor(gameType string) (model.GameProcessor, error) {
	v, ok := g.procs[gameType]
	if !ok {
		return nil, errors.New("wrong processr")
	}
	return v, nil
}
