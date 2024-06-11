package model

import "github.com/beevik/guid"

type MatcherRoom struct {
	ID            guid.Guid
	IsNew         bool   // новая комната
	Status        string // wait/game
	Players       []MatcherPlayer
	GameID        string
	StatusChanged bool
	State         string
}

type MatcherPlayer struct {
	IsNew    bool // новый игрок в комнате
	PlayerID guid.Guid
}

// запрос игрока на комнату
type RoomQuery struct {
	PlayerID guid.Guid
	GameID   string
}
