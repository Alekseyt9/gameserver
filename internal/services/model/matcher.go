package model

import "github.com/beevik/guid"

type MatcherRoom struct {
	Players []MatcherPlayer
}

type MatcherPlayer struct {
	IsNew    bool // новый	игрок в комнате
	PlayerId guid.Guid
}
