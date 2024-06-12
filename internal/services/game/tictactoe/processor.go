package game

import (
	"fmt"
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

func (p *TTCProcessor) Init(players []model.MatcherPlayer) (string, error) {
	return start(players)
}

func (p *TTCProcessor) Process(ctx model.ProcessorCtx, st string, msg *model.GameMsg) error {
	m, err := createGameMessage(msg)
	if err != nil {
		return err
	}

	s, err := strToState(st)
	if err != nil {
		return err
	}

	switch m.Action {
	case "state":
		state(ctx, s, m, msg.PlayerID)
	case "move":
		move(ctx, s, m, msg.PlayerID)
	case "quit":
		err := quit(ctx, s, msg.PlayerID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *TTCProcessor) PlayerQuit(ctx model.ProcessorCtx, gameID string, playerID guid.Guid, st string) error {
	s, err := strToState(st)
	if err != nil {
		return err
	}

}

// создать сообщение об ошибке
func createErrorMsg(playerID guid.Guid, s string) model.SendMessage {
	return model.SendMessage{
		PlayerID: playerID,
		Message: fmt.Sprintf(`
		{
			"type": "game",
			"game": "tictactoe",
			"data": { 
				"action": "error",
				"data": {
					"message": "%s"
				}
			}
		}
		`, s),
	}
}

// сделать ход
func move(ctx model.ProcessorCtx, s *TTTState, m *TTTMessage, playerID guid.Guid) error {
	if playerID != s.Turn {
		ctx.AddSendMessage(createErrorMsg(playerID, "Сейчас ход другого игрока"))
		return nil
	}

	d, err := createMoveData(m.Data)
	if err != nil {
		return err
	}

	if d.Move[0] > size-1 && d.Move[1] > size-1 {
		ctx.AddSendMessage(createErrorMsg(playerID, "Ход за границами поля"))
		return nil
	}

	if s.Field[d.Move[0]][d.Move[1]] != 0 {
		ctx.AddSendMessage(createErrorMsg(playerID, "Клетка уже занята"))
		return nil
	}

	s.Field[d.Move[0]][d.Move[1]] = figureOf(s, playerID)
	s.Turn = playerID

	// проверить конец игры, ничью
	if checkWin(s.Field, figureOf(s, s.Players[0])) {
		s.Winner = 0
		s.State = "finished"
	} else if checkWin(s.Field, figureOf(s, s.Players[1])) {
		s.Winner = 1
		s.State = "finished"
	} else if checkDraw(s.Field) {
		s.State = "finished"
	}

	saveState(ctx, s)
	return nil
}

// фигура игрока
func figureOf(s *TTTState, playerID guid.Guid) byte {
	// у первого - крестик
	if s.Players[0] == playerID {
		return 1
	}
	return 0
}

// противоположный игрок
func oppositeOf(s *TTTState, playerID guid.Guid) guid.Guid {
	p := s.Players[0]
	if p == playerID {
		p = s.Players[1]
	}
	return p
}

func indexOf(s *TTTState, playerID guid.Guid) int {
	for i, v := range s.Players {
		if v == playerID {
			return i
		}
	}
	return -1
}

// игрок покинул комнату
func quit(ctx model.ProcessorCtx, s *TTTState, playerID guid.Guid) error {
	if len(s.Players) == 2 {
		winner := oppositeOf(s, playerID)
		s.Winner = indexOf(s, winner)
		s.State = "finished"

		m, err := createStateSendMsg(s, winner)
		if err != nil {
			return err
		}
		ctx.AddSendMessage(*m)
		saveState(ctx, s)
	}

	return nil
}

// создаем игроков в state, определяем первый ход
func start(players []model.MatcherPlayer) (string, error) {
	s := &TTTState{
		State:  "game",
		Turn:   players[rand.Intn(2)].PlayerID,
		Winner: -1,
	}

	// крестик ходит первый
	if s.Turn == players[0].PlayerID {
		s.Players = []guid.Guid{players[0].PlayerID, players[1].PlayerID}
	} else {
		s.Players = []guid.Guid{players[1].PlayerID, players[0].PlayerID}
	}

	json, err := stateToStr(s)
	if err != nil {
		return "", err
	}

	return json, nil
}

func state(ctx model.ProcessorCtx, s *TTTState, m *TTTMessage, playerID guid.Guid) error {
	msgs := make([]model.SendMessage, 0)

	for _, p := range s.Players {
		m, err := createStateSendMsg(s, p)
		if err != nil {
			return err
		}
		msgs = append(msgs, *m)
	}

	ctx.AddSendMessages(msgs)
	return nil
}

func saveState(ctx model.ProcessorCtx, s *TTTState) error {
	json, err := stateToStr(s)
	if err != nil {
		return err
	}
	ctx.SetState(json)
	return nil
}

func stateToStr(s *TTTState) (string, error) {
	json, err := s.MarshalJSON()
	return string(json), err
}

func strToState(state string) (*TTTState, error) {
	var newState TTTState
	err := newState.UnmarshalJSON([]byte(state))
	return &newState, err
}

func createGameMessage(msg *model.GameMsg) (*TTTMessage, error) {
	var res TTTMessage
	err := res.UnmarshalJSON([]byte(msg.Data))
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func createMoveData(s string) (*TTTMoveData, error) {
	var res TTTMoveData
	err := res.UnmarshalJSON([]byte(s))
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// создать сообщение с состоянием игры для посылки игроку
func createStateSendMsg(s *TTTState, playerID guid.Guid) (*model.SendMessage, error) {
	ss := TTTSendState{
		Field:   s.Field,
		Players: s.Players,
		Turn:    s.Turn,
		You:     playerID,
		State:   s.State,
	}

	json, err := ss.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return &model.SendMessage{
		PlayerID: playerID,
		Message:  string(json),
	}, nil
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
