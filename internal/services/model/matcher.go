package model

import "github.com/beevik/guid"

type MatcherRoom struct {
	IsNew   bool   // новая комната
	Status  string // wait/game
	Players []MatcherPlayer
	GameID  string
}

type MatcherPlayer struct {
	IsNew    bool // новый игрок в комнате
	PlayerId guid.Guid
}

// запрос игрока на комнату
type RoomQuery struct {
	PlayerID guid.Guid
	GameID   string
}
