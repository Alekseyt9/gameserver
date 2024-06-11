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

func (p *TTCProcessor) Process(ctx model.ProcessorCtx, st string, msg *model.GameMsg) error {
	m := createGameMessage(msg)
	s := strToState(st)

	switch m.Kind {
	case "start":
		err := start(ctx, s, m)
		if err != nil {
			return err
		}
	case "state":
		state(ctx, s, m, msg.PlayerID)
	case "move":
		move(ctx, s, m, msg.PlayerID)
	case "quit":
		err := quit(ctx, s, m, msg.PlayerID)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO
func createErrorMsg(playerID guid.Guid, s string) model.SendMessage {
	return model.SendMessage{}
}

// сделать ход
func move(ctx model.ProcessorCtx, s *TTTState, m *TTTMessage, playerID guid.Guid) {
	if playerID != s.turn {
		ctx.AddSendMessage(createErrorMsg(playerID, "Сейчас ход другого игрока"))
		return
	}

	d := m.Data.(TTTMoveData)
	if d.Move[0] > size-1 && d.Move[1] > size-1 {
		ctx.AddSendMessage(createErrorMsg(playerID, "Ход за границами поля"))
		return
	}

	if s.field[d.Move[0]][d.Move[1]] != 0 {
		ctx.AddSendMessage(createErrorMsg(playerID, "Клетка уже занята"))
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

	saveState(ctx, s)
}

// фигура игрока
func figureOf(s *TTTState, playerID guid.Guid) byte {
	// у первого - крестик
	if s.players[0] == playerID {
		return 1
	}
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
func quit(ctx model.ProcessorCtx, s *TTTState, m *TTTMessage, playerID guid.Guid) error {
	if len(s.players) == 2 {
		winner := oppositeOf(s, playerID)
		s.winner = &winner
		s.state = "finished"

		ctx.AddSendMessage(createStateSendMsg(ctx, winner))
		saveState(ctx, s)
	}

	return nil
}

// создаем игроков в state, определяем первый ход
func start(ctx model.ProcessorCtx, s *TTTState, m *TTTMessage) error {
	d := m.Data.(TTTStartData)
	s = &TTTState{
		state: "game",
		turn:  d.Players[rand.Intn(2)],
	}

	// крестик ходит первый
	if s.turn == d.Players[0] {
		s.players = []guid.Guid{d.Players[0], d.Players[1]}
	} else {
		s.players = []guid.Guid{d.Players[1], d.Players[0]}
	}

	saveState(ctx, s)

	return nil
}

func state(ctx model.ProcessorCtx, s *TTTState, m *TTTMessage, playerID guid.Guid) {
	msgs := make([]model.SendMessage, 0)

	for _, p := range s.players {
		msgs = append(msgs, createStateSendMsg(ctx, p))
	}

	ctx.AddSendMessages(msgs)
}

func saveState(ctx model.ProcessorCtx, s *TTTState) {
	ctx.SetState(stateToStr(s))
}

func stateToStr(s *TTTState) string {
	return ""
}

func strToState(state string) *TTTState {
	return nil
}

func createGameMessage(msg *model.GameMsg) *TTTMessage {
	return nil
}

func createStateSendMsg(ctx model.ProcessorCtx, playerID guid.Guid) model.SendMessage {
	return model.SendMessage{}
}

// проверить ничью (все поля заняты)
func checkDraw(board [size][size]byte) bool {
	c := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if board[i][j] > 0 {
				c++
			}
		}
	}
	return c == size*size
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
