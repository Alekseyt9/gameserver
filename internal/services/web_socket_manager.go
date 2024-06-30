package services

import (
	"gameserver/internal/services/model"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olahol/melody"
)

type WebSocketManager struct {
	playerManager *PlayerManager
	roomManager   *RoomManager
	log           *slog.Logger
	ws            *melody.Melody
}

const (
	playerIDKey = "playerID"
	roomIDKey   = "roomID"
)

func NewWSManager(pm *PlayerManager, rm *RoomManager, log *slog.Logger) *WebSocketManager {
	m := &WebSocketManager{
		playerManager: pm,
		roomManager:   rm,
		log:           log,
		ws:            melody.New(),
	}
	ws := m.ws

	ws.HandleConnect(m.handleConnect)
	ws.HandleDisconnect(m.handleDisconnect)
	ws.HandleMessage(m.handleMessage)

	return m
}

func (m *WebSocketManager) handleConnect(s *melody.Session) {
	playerID, err := getPlayerID(s)
	if err != nil {
		m.log.Error("Ошибка получения playerID из куки")
	}
	sendCh := m.playerManager.GetOrCreateChan(*playerID)

	go func() {
		for msg := range sendCh {
			// сохраняем playerID и roomID в сессии
			_, ok := s.Keys[playerIDKey]
			if !ok {
				s.Keys[playerIDKey] = msg.PlayerID
			}
			_, ok = s.Keys[roomIDKey]
			if !ok && msg.RoomID != nil {
				s.Keys[roomIDKey] = *msg.RoomID
			}

			err = s.Write([]byte(msg.Message))
			if err != nil {
				m.log.Error("failed to write message", "error", err)
			}
		}
	}()
}

func (m *WebSocketManager) handleDisconnect(s *melody.Session) {
	pID, ok := s.Keys[playerIDKey]
	if !ok {
		m.log.Error("Ошибка получения playerID из сессии")
		return
	}
	playerID, ok := pID.(uuid.UUID)
	if !ok {
		m.log.Error("Ошибка приведения playerID к типу uuid.UUID")
		return
	}

	m.playerManager.DeleteChan(playerID)
}

func (m *WebSocketManager) handleMessage(s *melody.Session, data []byte) {
	pID, ok := s.Keys[playerIDKey]
	if !ok {
		m.log.Error("Ошибка получения playerID из сессии")
		return
	}
	playerID, ok := pID.(uuid.UUID)
	if !ok {
		m.log.Error("Ошибка приведения playerID к типу uuid.UUID")
		return
	}

	rID, ok := s.Keys[roomIDKey]
	if !ok {
		m.log.Error("Ошибка получения roomID из сессии")
		return
	}
	roomID, ok := rID.(uuid.UUID)
	if !ok {
		m.log.Error("Ошибка приведения roomID к типу uuid.UUID")
		return
	}

	m.log.Info("Message recieved", "msg", string(data))

	msg, err := createGameMsg(data, playerID)
	if err != nil {
		m.log.Error("Ошибка создания GameMsg", "error", err)
	}

	ch := m.roomManager.GetOrCreateChan(roomID)
	ch <- *msg
}

func (m *WebSocketManager) UpgradeToWebSocket(c *gin.Context, playerID uuid.UUID) {
	err := m.ws.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{
		playerIDKey: playerID,
	})
	if err != nil {
		m.log.Error("HandleRequestWithKeys")
		return
	}
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
