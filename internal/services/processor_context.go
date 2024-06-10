package services

import (
	"context"
	"gameserver/internal/services/model"
	"gameserver/internal/services/store"

	"github.com/beevik/guid"
)

type GameProcessorCtx struct {
	store         store.Store
	playerManager *PlayerManager
	roomID        guid.Guid
	gameID        string
}

func NewGameProcessorCtx(store store.Store, pm *PlayerManager, roomID guid.Guid, gameID string) *GameProcessorCtx {
	return &GameProcessorCtx{
		store:         store,
		roomID:        roomID,
		gameID:        gameID,
		playerManager: pm,
	}
}

func (c *GameProcessorCtx) SaveState(s string) error {
	ctx := context.Background()
	err := c.store.SetRoomState(ctx, c.roomID, s)
	return err
}

func (c *GameProcessorCtx) SendMessages(msgs []model.SendMessage) {
	for _, msg := range msgs {
		c.SendMessage(msg)
	}
}

func (c *GameProcessorCtx) SendMessage(msg model.SendMessage) {
	chp := c.playerManager.GetChan(msg.PlayerID)
	if chp != nil {
		*chp <- msg
	}
}

func (c *GameProcessorCtx) SendError(playerID guid.Guid, text string) {

}
