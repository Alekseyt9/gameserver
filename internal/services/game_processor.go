package services

import "gameserver/internal/services/model"

type GameProcessor interface {
	Process(ctx *model.GameProcessorCtx, state string, msg string) error
}
