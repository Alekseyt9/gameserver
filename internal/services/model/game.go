package model

import "github.com/google/uuid"

// общая информация об игре.
type GameInfo struct {
	PlayerCount int    // количество игроков в игре.
	ContentURL  string // ссылка на index.html игры.
	TurnTimeout int    // ограничение времени на ход.
}

type GameMsg struct {
	Type     string
	GameID   string
	PlayerID uuid.UUID
	Data     map[string]interface{}
}

// частично десериализованное сообщение с клиента.
//
//easyjson:json
type InMsg struct {
	Type   string                 `json:"type"`
	GameID string                 `json:"gameid"`
	Data   map[string]interface{} `json:"data"`
}

type SendMessage struct {
	PlayerID uuid.UUID
	RoomID   *uuid.UUID
	Message  string
}

func NewSendMessage(playerID uuid.UUID, roomID *uuid.UUID, msg string) SendMessage {
	return SendMessage{
		PlayerID: playerID,
		RoomID:   roomID,
		Message:  msg,
	}
}

type GameProcessor interface {
	GetInfo() *GameInfo
	Init(players []uuid.UUID) (string, error)
	Process(ctx ProcessorCtx, state string, msg *GameMsg) error
}

type ProcessorCtx interface {
	SetState(s string)
	AddSendMessage(msg SendMessage)
	AddSendMessages(msgs []SendMessage)
}
