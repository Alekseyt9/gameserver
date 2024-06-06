package services

import (
	"errors"
	game "gameserver/internal/services/game/tictactoe"
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"
)

type GameManager struct {
	procs map[string]GameProcessor
	store store.Store
}

const GameMsgType = "game"

func CreareGameManager(store store.Store) *GameManager {
	man := &GameManager{
		procs: make(map[string]GameProcessor),
		store: store,
	}
	man.procs["tictactoe"] = &game.TTCProcessor{}

	return man
}

func (g *GameManager) Process(msg *model.ClientMsg) error {
	// поднимает state для пользователя+тип игры
	// передает в обработчик

	if msg.GameType != GameMsgType {
		return errors.New("wrong MessageType")
	}

	state, err := g.store.GetRoomState(msg.PlayerID, msg.GameType)
	if err != nil {
		return err
	}

	proc, err := g.getProcessor(msg.GameType)
	if err != nil {
		return err
	}

	procCtx := &model.GameProcessorCtx{}
	proc.Process(procCtx, state, msg.Data)

	return nil
}

func (g *GameManager) getProcessor(gameType string) (GameProcessor, error) {
	v, ok := g.procs[gameType]
	if !ok {
		return nil, errors.New("wrong processr")
	}
	return v, nil
}
