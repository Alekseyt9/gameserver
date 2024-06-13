package handlers

import (
	"gameserver/internal/services/model"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) RegisterPlayer(c *gin.Context) {
	player := model.Player{
		ID:   uuid.New(),
		Name: "player" + h.generateRandomString(5),
	}

	maxAge := 2147483647
	c.SetCookie("playerID", player.ID.String(), maxAge,
		"/", "", false, true)

	err := h.store.CreatePlayer(c.Request.Context(), &player)
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
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
