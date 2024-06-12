package handlers

import (
	"net/http"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

type RoomRequest struct {
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

	var req RoomRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	conRes, err := h.roomManager.PlayerConnect(c.Request.Context(), *playerID, req.GameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PlayerConnect"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": conRes.State,
	})
}

func (h *Handler) QuitRoom(c *gin.Context) {
	// TODO вынести в middleware
	plID, err := c.Cookie("playerID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "PlayerID cookie not found"})
		return
	}

	var req RoomRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	playerID, err := guid.ParseString(plID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.roomManager.PlayerQuit(c.Request.Context(), req.GameID, *playerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quit room successfully",
	})
}
