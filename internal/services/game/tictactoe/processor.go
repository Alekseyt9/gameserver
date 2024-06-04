package game

import "gameserver/internal/services/model"

type TTCProcessor struct {
}

func (p *TTCProcessor) Process(ctx *model.GameProcessorCtx, state string, msg string) error {
	return nil
}
