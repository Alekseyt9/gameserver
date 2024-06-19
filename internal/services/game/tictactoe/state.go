package game

import "github.com/google/uuid"

// состояние в базе
//
//easyjson:json
type TTTState struct {
	Field    [15][15]int `json:"field"`
	Players  []uuid.UUID `json:"players"`
	State    string      `json:"state"`
	Winner   int         `json:"winner"`
	Turn     uuid.UUID   `json:"turn"`
	WinLine  [][]int     `json:"winline"`
	LastMove [2]int      `json:"lastmove"`
}

// состояние для посылки игрокам
//
//easyjson:json
type TTTSendState struct {
	Field    [15][15]int `json:"field"`
	Players  []uuid.UUID `json:"players"`
	Turn     uuid.UUID   `json:"turn"`
	You      uuid.UUID   `json:"you"`
	State    string      `json:"state"`
	Winner   int         `json:"winner"`
	WinLine  [][]int     `json:"winline"`
	LastMove [2]int      `json:"lastmove"`
}
