package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ConnectRoom(c *gin.Context) {
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
