package model

import "github.com/beevik/guid"

// частично десериализованное сообщение с клиента
type ClientMsg struct {
	MessageType string
	PlayerID    guid.Guid
	GameType    string
	Data        string
}
