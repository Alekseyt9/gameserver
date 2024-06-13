package model

import uuid "github.com/google/uuid"

type MatcherRoom struct {
	ID            uuid.UUID
	IsNew         bool   // новая комната
	Status        string // wait/game
	Players       []*MatcherPlayer
	GameID        string
	StatusChanged bool
	State         string
}

type MatcherPlayer struct {
	IsNew    bool // новый игрок в комнате
	PlayerID uuid.UUID
}

// запрос игрока на комнату
type RoomQuery struct {
	PlayerID uuid.UUID
	GameID   string
}
