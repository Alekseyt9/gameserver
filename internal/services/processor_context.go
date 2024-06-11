package services

import (
	"gameserver/internal/services/model"

	"github.com/beevik/guid"
)

type GameProcessorCtx struct {
	roomID       guid.Guid
	gameID       string
	gameState    string
	sendMessages []model.SendMessage
}

func NewGameProcessorCtx(roomID guid.Guid, gameID string) *GameProcessorCtx {
	return &GameProcessorCtx{
		roomID:       roomID,
		gameID:       gameID,
		sendMessages: make([]model.SendMessage, 0),
	}
}

func (c *GameProcessorCtx) SetState(s string) {
	c.gameState = s
}

func (c *GameProcessorCtx) AddSendMessage(msg model.SendMessage) {
	c.sendMessages = append(c.sendMessages, msg)
}

func (c *GameProcessorCtx) AddSendMessages(msgs []model.SendMessage) {
	c.sendMessages = append(c.sendMessages, msgs...)
}
