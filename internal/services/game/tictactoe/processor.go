package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"gameserver/internal/services/model"
	"math/rand"

	"github.com/google/uuid"
)

type TTCProcessor struct {
}

var directions = [][2]int{ //nolint:gochecknoglobals // это константа.
	{1, 0},
	{0, 1},
	{1, 1},
	{1, -1},
}

const (
	size          = 15
	empty         = 0
	winCount      = 5
	playerCount   = 2
	turnTimeout   = 60 * 5
	contentURL    = "/content/games/tictactoe/index.html"
	playerFigure1 = 1
	playerFigure2 = 2
	actionState   = "state"
	actionMove    = "move"
	actionQuit    = "quit"
	stateFinished = "finished"
)

func New() model.GameProcessor {
	return &TTCProcessor{}
}

func (p *TTCProcessor) GetInfo() *model.GameInfo {
	return &model.GameInfo{
		PlayerCount: playerCount,
		ContentURL:  contentURL,
		TurnTimeout: turnTimeout,
	}
}

func (p *TTCProcessor) Init(players []uuid.UUID) (string, error) {
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
	case actionState:
		err = state(ctx, s, msg.PlayerID)
		if err != nil {
			return err
		}
	case actionMove:
		err = move(ctx, s, m, msg.PlayerID)
		if err != nil {
			return err
		}
	case actionQuit:
		err = quit(ctx, s, msg.PlayerID)
		if err != nil {
			return err
		}
	}
	return nil
}

// создать сообщение об ошибке.
func createErrorMsg(playerID uuid.UUID, s string) model.SendMessage {
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

// сделать ход.
func move(ctx model.ProcessorCtx, s *TTTState, m *TTTMessage, playerID uuid.UUID) error {
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
	s.Turn = oppositeOf(s, playerID)
	s.LastMove = d.Move

	// проверить конец игры, ничью.
	var line *[][]int
	chWin, line := checkWin(s.Field, figureOf(s, s.Players[0]))
	if chWin {
		s.Winner = 0
		s.State = stateFinished
		s.WinLine = *line
	} else {
		chWin, line = checkWin(s.Field, figureOf(s, s.Players[1]))
		if chWin {
			s.Winner = 1
			s.State = stateFinished
			s.WinLine = *line
		} else if checkDraw(s.Field) {
			s.State = stateFinished
		}
	}

	err = saveState(ctx, s)
	if err != nil {
		return err
	}

	for _, p := range s.Players {
		var msg *model.SendMessage
		msg, err = createStateSendMsg(s, p)
		if err != nil {
			return err
		}

		ctx.AddSendMessage(*msg)
	}

	return nil
}

// фигура игрока.
func figureOf(s *TTTState, playerID uuid.UUID) byte {
	// у первого - крестик.
	if s.Players[0] == playerID {
		return playerFigure1
	}
	return playerFigure2
}

// противоположный игрок.
func oppositeOf(s *TTTState, playerID uuid.UUID) uuid.UUID {
	p := s.Players[0]
	if p == playerID {
		p = s.Players[1]
	}
	return p
}

func indexOf(s *TTTState, playerID uuid.UUID) int {
	for i, v := range s.Players {
		if v == playerID {
			return i
		}
	}
	return -1
}

// игрок покинул комнату.
func quit(ctx model.ProcessorCtx, s *TTTState, playerID uuid.UUID) error {
	if len(s.Players) == playerCount {
		winner := oppositeOf(s, playerID)
		s.Winner = indexOf(s, winner)
		s.State = "finished"

		m, err := createStateSendMsg(s, winner)
		if err != nil {
			return err
		}
		ctx.AddSendMessage(*m)

		err = saveState(ctx, s)
		if err != nil {
			return err
		}
	}

	return nil
}

// создаем игроков в state, определяем первый ход.
func start(players []uuid.UUID) (string, error) {
	s := &TTTState{
		State:  "game",
		Turn:   players[rand.Intn(playerCount)], //nolint:gosec //rand
		Winner: -1,
	}

	// крестик ходит первый.
	if s.Turn == players[0] {
		s.Players = []uuid.UUID{players[0], players[1]}
	} else {
		s.Players = []uuid.UUID{players[1], players[0]}
	}

	json, err := stateToStr(s)
	if err != nil {
		return "", err
	}

	return json, nil
}

func state(ctx model.ProcessorCtx, s *TTTState, playerID uuid.UUID) error {
	msg, err := createStateSendMsg(s, playerID)
	if err != nil {
		return err
	}

	ctx.AddSendMessage(*msg)
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
	a, ok := msg.Data["action"]
	d, hasData := msg.Data["data"]

	if ok {
		m := &TTTMessage{
			Action: a.(string),
		}
		if hasData {
			v, err := json.Marshal(d)
			if err != nil {
				return nil, err
			}
			m.Data = string(v)
		}
		return m, nil
	}
	return nil, errors.New("wrong msg format")
}

func createMoveData(s string) (*TTTMoveData, error) {
	var res TTTMoveData
	err := res.UnmarshalJSON([]byte(s))
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// создать сообщение с состоянием игры для посылки игроку.
func createStateSendMsg(s *TTTState, playerID uuid.UUID) (*model.SendMessage, error) {
	ss := TTTSendState{
		Field:    s.Field,
		Players:  s.Players,
		Turn:     s.Turn,
		You:      playerID,
		State:    s.State,
		Winner:   s.Winner,
		WinLine:  s.WinLine,
		LastMove: s.LastMove,
	}

	json, err := ss.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return &model.SendMessage{
		PlayerID: playerID,
		Message: fmt.Sprintf(`
		{
			"type": "game",
			"gameid": "tictactoe",
			"data": {
				"action": "state",
				"data": %s
			}
		}
		`, string(json)),
	}, nil
}

// проверить ничью (все поля заняты).
func checkDraw(board [size][size]byte) bool {
	c := 0
	for i := range size {
		for j := range size {
			if board[i][j] > 0 {
				c++
			}
		}
	}
	return c == size*size
}

func checkWin(board [size][size]byte, figure byte) (bool, *[][]int) {
	for i := range size {
		for j := range size {
			if board[i][j] == figure {
				for _, dir := range directions {
					if checkDirection(board, figure, i, j, dir[0], dir[1]) {
						return true, &[][]int{{i, j}, {i + dir[0]*winCount, j + dir[1]*winCount}}
					}
				}
			}
		}
	}

	return false, nil
}

func checkDirection(board [size][size]byte, player byte, x, y, dx, dy int) bool {
	count := 0
	for k := range winCount {
		nx, ny := x+dx*k, y+dy*k
		if nx < 0 || nx >= size || ny < 0 || ny >= size || board[nx][ny] != player {
			return false
		}
		count++
	}

	return count == winCount
}
