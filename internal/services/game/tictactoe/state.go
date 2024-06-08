package game

import "github.com/beevik/guid"

type TTTState struct {
	field   [15][15]int
	players [2]guid.Guid
	state   string // 'game', 'draw', 'finished'
	winner  *guid.Guid
	turn    guid.Guid
}
