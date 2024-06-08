package game

import (
	"gameserver/internal/services/model"
	"math/rand"
)

type TTCProcessor struct {
}

func New() model.GameProcessor {
	return &TTCProcessor{}
}

func (p *TTCProcessor) GetInfo() *model.GameInfo {
	return &model.GameInfo{
		PlayerCount: 2,
		ContentURL:  "/content/games/tictactoe/index.html",
		TurnTimeout: 60 * 5,
	}
}

func (p *TTCProcessor) Process(ctx *model.GameProcessorCtx, state string, msg *model.GameMsg) error {
	m := p.getMessage(msg)
	s := p.getState(state)

	switch m.Kind {
	case "start":
		p.start(ctx, s, m)
	case "move":
		p.move(ctx, s, m)
	case "pquit":
		p.pquit(ctx, s, m)
	}
	return nil
}

// игрок покинул комнату
func (p *TTCProcessor) pquit(ctx *model.GameProcessorCtx, s *TTTState, m *TTTMessage) {
}

// сделать ход
func (p *TTCProcessor) move(ctx *model.GameProcessorCtx, s *TTTState, m *TTTMessage) {
}

// проверить конец игры
func (p *TTCProcessor) checkFinish(s *TTTState) {

}

// создаем игроков в state, определяем первый ход
func (p *TTCProcessor) start(ctx *model.GameProcessorCtx, s *TTTState, m *TTTMessage) {
	d := m.Data.(TTTStartData)
	s = &TTTState{
		state:   "game",
		players: d.Players,
		turn:    d.Players[rand.Intn(2)],
	}
	ctx.SaveState(p.stateToStr(s))
}

func (p *TTCProcessor) stateToStr(s *TTTState) string {
	return ""
}

func (p *TTCProcessor) getState(state string) *TTTState {
	return nil
}

func (p *TTCProcessor) getMessage(msg *model.GameMsg) *TTTMessage {
	return nil
}
