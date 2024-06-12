package game

import "github.com/beevik/guid"

// состояние в базе
//easyjson:json
type TTTState struct {
	Field   [15][15]byte `json:"field"`
	Players []guid.Guid  `json:"players"`
	State   string       `json:"state"`
	Winner  int          `json:"winner"`
	Turn    guid.Guid    `json:"turn"`
}

// состояние для посылки игрокам
//easyjson:json
type TTTSendState struct {
	Field   [15][15]byte `json:"field"`
	Players []guid.Guid  `json:"players"`
	Turn    guid.Guid    `json:"turn"`
	You     guid.Guid    `json:"you"`
	State   string       `json:"state"`
	Winner  int          `json:"winner"`
}
