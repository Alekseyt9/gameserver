package model

import "github.com/beevik/guid"

type GameProcessorCtx struct {
	roomID guid.Guid
	gameID string
}

type SendMessage struct {
	playerID guid.Guid
	message  string
}

func (c *GameProcessorCtx) SaveState(s string) {

}

func (c *GameProcessorCtx) SendMessages(msgs []SendMessage) {

}

func (c *GameProcessorCtx) SendError(playerID guid.Guid, text string) {

}
