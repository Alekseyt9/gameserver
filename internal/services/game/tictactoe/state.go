package game

import "github.com/beevik/guid"

type TTCState struct {
	field   [15][15]int
	players []guid.Guid
	state   string // 'game', 'draw', 'finished'
	winner  *guid.Guid
}
