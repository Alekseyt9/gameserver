package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

func RegisterPlayer(c *gin.Context) {
	playerID := guid.New()
	playerName := "player" + generateRandomString(5)

	maxAge := 2147483647
	c.SetCookie("playerID", playerID.String(), maxAge,
		"/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"playerName": playerName,
		"playerID":   playerID.String(),
	})
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
