package model

import (
	"github.com/beevik/guid"
)

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
	GameID      string
	Data        string
}

type SendMessage struct {
	PlayerID guid.Guid
	Message  string
}

type GameProcessor interface {
	GetInfo() *GameInfo
	Init(players []MatcherPlayer) (string, error) // TODO guids?
	Process(ctx ProcessorCtx, state string, msg *GameMsg) error
}

type ProcessorCtx interface {
	SetState(s string)
	AddSendMessage(msg SendMessage)
	AddSendMessages(msgs []SendMessage)
}
