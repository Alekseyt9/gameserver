package game

import "github.com/beevik/guid"

type TTTState struct {
	field   [15][15]byte
	players []guid.Guid
	state   string // 'game', 'finished'
	winner  *guid.Guid
	turn    guid.Guid
}
