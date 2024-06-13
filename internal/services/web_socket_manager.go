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
		ws.HandleRequest(c.Writer, c.Request)
	})

	ws.HandleConnect(func(s *melody.Session) {
		playerID, err := getPlayerID(s)
		if err != nil {
			log.Printf("Ошибка получения playerID из куки")
		}

		sendCh := m.playerManager.GetOrCreateChan(*playerID)

		go func() {
			for msg := range sendCh {
				s.Write([]byte(msg.Message))
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
			log.Printf("Ошибка получения playerID из куки %w", err)
		}

		msg, err := createGameMsg(data, *playerID)
		if err != nil {
			log.Printf("Ошибка создания GameMsg %w", err)
		}

		// комната уже есть, тк в игре
		room, err := m.roomManager.GetExistingRoom(s.Request.Context(), msg.GameID, msg.PlayerID)
		if err != nil {
			log.Printf("Ошибка получения комнаты %w", err)
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
