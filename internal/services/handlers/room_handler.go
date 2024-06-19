package handlers

import (
	"gameserver/internal/services"
	"gameserver/internal/services/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoomRequest struct {
	GameID string `json:"gameID" binding:"required"`
}

func (h *Handler) ConnectRoom(c *gin.Context) {
	pID, err := c.Cookie("playerID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "PlayerID cookie not found"})
		return
	}
	playerID, err := uuid.Parse(pID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid.Parse(pID)"})
		return
	}

	gameID := c.Query("gameid")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing gameID parameter"})
		return
	}

	var conRes *services.PlayerConnectResult
	conRes, err = h.roomManager.PlayerConnect(c.Request.Context(), playerID, gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PlayerConnect"})
		return
	}
	h.log.Info("ConnectRoom", "playerID", playerID.String(), "result", conRes.State)

	if conRes.State == "game" {
		err = h.playerManager.SendToPlayer(playerID,
			model.NewSendMessage(playerID, conRes.RoomID, model.CreateStartGameMsg(conRes.ContentLink)))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SendToPlayer"})
			return
		}
	}

	h.wsManager.UpgradeToWebSocket(c, playerID)
}

func (h *Handler) QuitRoom(c *gin.Context) {
	plID, err := c.Cookie("playerID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "PlayerID cookie not found"})
		return
	}

	var req RoomRequest
	if err = c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var playerID uuid.UUID
	playerID, err = uuid.Parse(plID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("QuitRoom", "playerID", playerID.String())

	err = h.roomManager.PlayerQuit(c.Request.Context(), req.GameID, playerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quit room successfully",
	})
}
