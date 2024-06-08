package model

import "github.com/beevik/guid"

// общая информация об игре
type GameInfo struct {
	PlayerCount int    // количество игроков в игре
	ContentURL  string // ссылка на index.html игры
	TurnTimeout int    // ограничение времени на ход
}

// частично десериализованное сообщение с клиента
type GameMsg struct {
	MessageType string
	PlayerID    guid.Guid
	GameType    string
	Data        string
}

type GameProcessor interface {
	Process(ctx *GameProcessorCtx, state string, msg GameMsg) error
	GetInfo() *GameInfo
}
