package main

import (
	"net/http"
	"path/filepath"
	"text/template"

	"gameserver/internal/services/handlers"
	"gameserver/internal/services/store"

	"github.com/gin-gonic/gin"
)

type PageData struct {
	WebSocketURL string
}

type Config struct {
	ConnectionString string
}

func main() {
	r := gin.Default()

	cfg := &Config{}

	fileServer(r)
	regHandlers(r, cfg)
}

func regHandlers(r *gin.Engine, cfg *Config) {

	store, err := store.NewDBStore(cfg.ConnectionString)
	if err != nil {
		// TODO log
	}

	h := handlers.New(store)
	r.POST("/api/player/register", h.RegisterPlayer)
	r.POST("/api/room/connect", h.ConnectRoom)
	r.POST("/api/room/quit", h.QuitRoom)
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
