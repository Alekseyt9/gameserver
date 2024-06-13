package services

import (
	"gameserver/internal/services/model"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olahol/melody"
)

type WebSocketManager struct {
	playerManager *PlayerManager
	roomManager   *RoomManager
}

func NewWSManager(router *gin.Engine, pm *PlayerManager, rm *RoomManager) *WebSocketManager {
	m := &WebSocketManager{
		playerManager: pm,
		roomManager:   rm,
	}

	ws := melody.New()

	router.GET("/ws", func(c *gin.Context) {
		err := ws.HandleRequest(c.Writer, c.Request)
		if err != nil {
			log.Printf("ws.HandleRequest error: %v", err)
		}
	})

	ws.HandleConnect(func(s *melody.Session) {
		playerID, err := getPlayerID(s)
		if err != nil {
			log.Printf("Ошибка получения playerID из куки")
		}

		sendCh := m.playerManager.GetOrCreateChan(*playerID)

		go func() {
			for msg := range sendCh {
				err = s.Write([]byte(msg.Message))
				if err != nil {
					log.Printf("failed to write message: %v", err)
				}
			}
		}()
	})

	// TODO проверить, что есть playerID при отключении
	ws.HandleDisconnect(func(s *melody.Session) {
		playerID, err := getPlayerID(s)
		if err != nil {
			log.Printf("Ошибка получения playerID из куки")
		}

		m.playerManager.DeleteChan(*playerID)
	})

	ws.HandleMessage(func(s *melody.Session, data []byte) {
		playerID, err := getPlayerID(s)
		if err != nil {
			log.Printf("Ошибка получения playerID из куки %v", err)
		}

		msg, err := createGameMsg(data, *playerID)
		if err != nil {
			log.Printf("Ошибка создания GameMsg %v", err)
		}

		// комната уже есть, тк в игре
		room, err := m.roomManager.GetExistingRoom(s.Request.Context(), msg.GameID, msg.PlayerID)
		if err != nil {
			log.Printf("Ошибка получения комнаты %v", err)
		}

		ch := m.roomManager.GetOrCreateChan(room.ID)
		ch <- *msg
	})

	return m
}

func createGameMsg(data []byte, playerID uuid.UUID) (*model.GameMsg, error) {
	var m model.InMsg
	err := m.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return &model.GameMsg{
		Type:     m.Type,
		GameID:   m.GameID,
		PlayerID: playerID,
		Data:     m.Data,
	}, nil
}

func getPlayerID(s *melody.Session) (*uuid.UUID, error) {
	cookies := s.Request.Cookies()
	var playerID uuid.UUID
	var err error
	for _, cookie := range cookies {
		if cookie.Name == "playerID" {
			playerID, err = uuid.Parse(cookie.Value)
			break
		}
	}
	return &playerID, err
}
