package handlers

import (
	"gameserver/internal/services"
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

	var req RoomRequest
	if err = c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var conRes *services.PlayerConnectResult
	conRes, err = h.roomManager.PlayerConnect(c.Request.Context(), playerID, req.GameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PlayerConnect"})
		return
	}

	a := make(map[string]any, 0)
	a["state"] = conRes.State
	if conRes.State == "game" {
		a["contentLink"] = conRes.ContentLink
	}

	h.log.Info("ConnectRoom", "playerID", playerID.String(), "result", conRes.State)

	c.JSON(http.StatusOK, a)
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
