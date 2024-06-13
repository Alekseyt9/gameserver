package run

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"gameserver/internal/services"
	"gameserver/internal/services/handlers"
	"gameserver/internal/services/store"
)

type PageData struct {
	WebSocketURL string
}

type Config struct {
	Address     string `env:"ADDRESS"`
	DataBaseDSN string `env:"DATABASE_DSN"`
}

func Run(cfg *Config) error {
	s, err := store.NewDBStore(cfg.DataBaseDSN)
	if err != nil {
		return err
	}
	pm := services.NewPlayerManager(s)
	gm := services.NewGameManager(s, pm)
	m, err := services.NewMatcher(s, pm, gm)
	if err != nil {
		return err
	}
	rm := services.NewRoomManager(s, gm, pm, m)

	r := Router(s, pm, rm, cfg)

	// TODO из командной строки или переменной окружения
	err = r.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}

func Router(s store.Store, pm *services.PlayerManager, rm *services.RoomManager, cfg *Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	services.NewWSManager(r, pm, rm)
	setupFileServer(r)
	setupHandlers(r, s, rm)
	return r
}

func setupHandlers(r *gin.Engine, s store.Store, rm *services.RoomManager) {
	h := handlers.New(s, rm)
	r.POST("/api/player/register", h.RegisterPlayer)
	r.POST("/api/room/connect", h.ConnectRoom)
	r.POST("/api/room/quit", h.QuitRoom)
}

func setupFileServer(r *gin.Engine) {
	r.Use(gzip.Gzip(gzip.BestCompression))

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
			WebSocketURL: "ws://localhost:8080/ws",
		}
		c.Writer.Header().Set("Content-Type", "text/html")
		tmpl.Execute(c.Writer, data)
	})
}
