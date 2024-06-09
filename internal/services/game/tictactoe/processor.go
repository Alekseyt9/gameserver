package game

import (
	"gameserver/internal/services/model"
	"math/rand"

	"github.com/beevik/guid"
)

type TTCProcessor struct {
}

var directions = [][2]int{
	{1, 0},
	{0, 1},
	{1, 1},
	{1, -1},
}

const (
	size     = 15
	empty    = 0
	winCount = 5
)

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

func (p *TTCProcessor) Process(ctx *model.GameProcessorCtx, st string, msg *model.GameMsg) error {
	m := getMessage(msg)
	s := getState(st)

	switch m.Kind {
	case "start":
		start(ctx, s, m)
	case "state":
		state(ctx, s, m, msg.PlayerID)
	case "move":
		move(ctx, s, m, msg.PlayerID)
	case "quit":
		quit(ctx, s, m, msg.PlayerID)
	}
	return nil
}

// сделать ход
func move(ctx *model.GameProcessorCtx, s *TTTState, m *TTTMessage, playerID guid.Guid) {
	if playerID != s.turn {
		ctx.SendError(playerID, "Сейчас ход другого игрока")
		return
	}

	d := m.Data.(TTTMoveData)
	if d.Move[0] > size-1 && d.Move[1] > size-1 {
		ctx.SendError(playerID, "Ход за границами поля")
		return
	}

	if s.field[d.Move[0]][d.Move[1]] != 0 {
		ctx.SendError(playerID, "Клетка уже занята")
		return
	}

	s.field[d.Move[0]][d.Move[1]] = figureOf(s, playerID)
	s.turn = playerID

	// проверить конец игры, ничью
	if checkWin(s.field, figureOf(s, s.players[0])) {
		s.winner = &s.players[0]
		s.state = "finished"
	} else if checkWin(s.field, figureOf(s, s.players[1])) {
		s.winner = &s.players[1]
		s.state = "finished"
	} else if checkDraw(s.field) {
		s.state = "finished"
	}

	ctx.SaveState(stateToStr(s))
}

// фигура игрока
func figureOf(s *TTTState, playerID guid.Guid) byte {
	return 0
}

// противоположный игрок
func oppositeOf(s *TTTState, playerID guid.Guid) guid.Guid {
	p := s.players[0]
	if p == playerID {
		p = s.players[1]
	}
	return p
}

// игрок покинул комнату
func quit(ctx *model.GameProcessorCtx, s *TTTState, m *TTTMessage, playerID guid.Guid) {
	if len(s.players) == 2 {
		winner := oppositeOf(s, playerID)
		s.winner = &winner
		s.state = "finished"

		// TODO разослать стейт

		ctx.SaveState(stateToStr(s))
	}
}

// создаем игроков в state, определяем первый ход
func start(ctx *model.GameProcessorCtx, s *TTTState, m *TTTMessage) {
	d := m.Data.(TTTStartData)
	s = &TTTState{
		state:   "game",
		players: d.Players[:],
		turn:    d.Players[rand.Intn(2)],
	}
	ctx.SaveState(stateToStr(s))
}

func state(ctx *model.GameProcessorCtx, s *TTTState, m *TTTMessage, playerID guid.Guid) {

}

func stateToStr(s *TTTState) string {
	return ""
}

func getState(state string) *TTTState {
	return nil
}

func getMessage(msg *model.GameMsg) *TTTMessage {
	return nil
}

func checkDraw(board [size][size]byte) bool {
	return false
}

func checkWin(board [size][size]byte, figure byte) bool {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if board[i][j] == figure {
				for _, dir := range directions {
					if checkDirection(board, figure, i, j, dir[0], dir[1]) {
						return true
					}
				}
			}
		}
	}

	return false
}

func checkDirection(board [size][size]byte, player byte, x, y, dx, dy int) bool {
	count := 0
	for k := 0; k < winCount; k++ {
		nx, ny := x+dx*k, y+dy*k
		if nx < 0 || nx >= size || ny < 0 || ny >= size || board[nx][ny] != player {
			return false
		}
		count++
	}

	return count == winCount
}
