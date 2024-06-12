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

type GameMsg struct {
	Type     string
	GameID   string
	PlayerID guid.Guid
	Data     map[string]interface{}
}

// частично десериализованное сообщение с клиента
//
//easyjson:json
type InMsg struct {
	Type   string                 `json:"type"`
	GameID string                 `json:"gameid"`
	Data   map[string]interface{} `json:"data"`
}

type SendMessage struct {
	PlayerID guid.Guid
	Message  string
}

type GameProcessor interface {
	GetInfo() *GameInfo
	Init(players []guid.Guid) (string, error)
	Process(ctx ProcessorCtx, state string, msg *GameMsg) error
}

type ProcessorCtx interface {
	SetState(s string)
	AddSendMessage(msg SendMessage)
	AddSendMessages(msgs []SendMessage)
}
