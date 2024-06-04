package services

import (
	game "gameserver/internal/services/game/tictactoe"
	"gameserver/internal/services/model"
)

type GameManager struct {
	procs map[string]*GameProcessor
}

func CreareGameManager() *GameManager {
	man := &GameManager{
		procs: make(map[string]*GameProcessor),
	}
	man.procs["tictactoe"] = &game.TTCProcessor{}

	return man
}

func Process(msg *model.ClientMsg) {

}
