package handlers

import (
	"gameserver/internal/services/model"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	playerRandCharCount = 5
)

func (h *Handler) RegisterPlayer(c *gin.Context) {
	pID, err := c.Cookie("playerID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "c.Cookie(playerID)"})
	}

	if pID != "" {
		var playerID uuid.UUID
		playerID, err = uuid.Parse(pID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "uuid.Parse"})
		}

		var player *model.Player
		player, err = h.store.GetPlayer(c.Request.Context(), playerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "store.GetPlayer"})
		}

		c.JSON(http.StatusOK, gin.H{
			"playerName": player.Name,
			"playerID":   player.ID.String(),
		})
	}

	player := model.Player{
		ID:   uuid.New(),
		Name: "player" + h.generateRandomString(playerRandCharCount),
	}

	maxAge := 2147483647
	c.SetCookie("playerID", player.ID.String(), maxAge,
		"/", "", false, true)

	err = h.store.CreatePlayer(c.Request.Context(), &player)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "create user in DB"})
	}

	c.JSON(http.StatusOK, gin.H{
		"playerName": player.Name,
		"playerID":   player.ID.String(),
	})
}

func (h *Handler) generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))] //nolint:gosec //rand
	}
	return string(b)
}
