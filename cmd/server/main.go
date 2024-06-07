package main

import (
	"net/http"
	"path/filepath"
	"text/template"

	"gameserver/internal/services/handlers"

	"github.com/gin-gonic/gin"
)

type PageData struct {
	WebSocketURL string
}

func main() {
	r := gin.Default()
	fileServer(r)
	regHandlers(r)
}

func regHandlers(r *gin.Engine) {
	r.POST("/api/player/register", handlers.RegisterPlayer)
}

func fileServer(r *gin.Engine) {
	contentDir := filepath.Join("..", "..", "internal", "content")
	r.StaticFS("/content", http.Dir(contentDir))

	r.GET("/", func(c *gin.Context) {
		if c.Request.URL.Path != "/" {
			c.String(http.StatusNotFound, "Page not found")
			return
		}
		tmplPath := filepath.Join(contentDir, "index.html")
		tmpl := template.Must(template.ParseFiles(tmplPath))

		// TODO задавать базовый адрес из командной строки
		data := PageData{
			WebSocketURL: "ws://localhost:3001/ws",
		}
		c.Writer.Header().Set("Content-Type", "text/html")
		tmpl.Execute(c.Writer, data)
	})

	err := r.Run(":8080")
	if err != nil {
		panic("Ошибка запуска сервера: " + err.Error())
	}
}
