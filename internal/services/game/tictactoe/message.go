package game

import "github.com/beevik/guid"

type TTTMessage struct {
	Kind string
	Data any // в зависимости от kind
}

type TTTStartData struct {
	Players [2]guid.Guid
}

type TTTMoveData struct {
	Move [2]byte
}
