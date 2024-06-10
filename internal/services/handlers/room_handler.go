package handlers

import (
	"net/http"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

type ConnectRoomRequest struct {
	GameID string `json:"gameID" binding:"required"`
}

func (h *Handler) ConnectRoom(c *gin.Context) {
	// TODO вынести в middleware
	pID, err := c.Cookie("playerID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "PlayerID cookie not found"})
		return
	}

	playerID, err := guid.ParseString(pID)
	if err != nil {
		// TODO
	}

	var req ConnectRoomRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	conRes, err := h.roomManager.PlayerConnect(c.Request.Context(), *playerID, req.GameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PlayerConnect"})
		return
	}

	// TODO
	c.JSON(http.StatusOK, gin.H{
		"message": conRes.State,
	})
}

func (h *Handler) QuitRoom(c *gin.Context) {
	// TODO вынести в middleware
	playerID, err := c.Cookie("playerID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "PlayerID cookie not found"})
		return
	}

	// TODO
	// здесь (если комната в игре) - посылаем начало игры с контентом
	// или в очередь на матчинг - возвращаем wait
	c.JSON(http.StatusOK, gin.H{
		"message":  "Connected to room successfully",
		"playerID": playerID,
	})
}
