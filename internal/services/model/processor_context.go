package model

import "github.com/beevik/guid"

type GameProcessorCtx struct {
	roomID guid.Guid
	gameID string
}

type SendMessage struct {
	message  string
	playerID guid.Guid
}

func (c *GameProcessorCtx) SaveState(s string) {

}

func (c *GameProcessorCtx) SendMessages() {

}
