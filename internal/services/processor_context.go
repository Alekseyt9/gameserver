package services

import (
	"gameserver/internal/services/model"

	"github.com/google/uuid"
)

type GameProcessorCtx struct {
	roomID       uuid.UUID
	gameID       string
	gameState    string
	sendMessages []model.SendMessage
}

func NewGameProcessorCtx(roomID uuid.UUID, gameID string) *GameProcessorCtx {
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
