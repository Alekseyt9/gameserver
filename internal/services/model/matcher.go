package model

import "github.com/beevik/guid"

type MatcherRoom struct {
	IsNew   bool   // новая комната
	Status  string // wait/game
	Players []MatcherPlayer
}

type MatcherPlayer struct {
	IsNew    bool // новый игрок в комнате
	PlayerId guid.Guid
}
