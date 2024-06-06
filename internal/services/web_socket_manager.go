package services

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type WebSocketManager struct {
}

func CreateWSManager(url string) {
	router := gin.Default()
	m := melody.New()

	router.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {

		/*
			go func() {
			  for msg := range ch {
				s.Write([]byte(msg))
			  }
			}()
		*/
	})

	router.Run(url)
}
