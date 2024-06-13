package main

import (
	"gameserver/internal/run"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PageData struct {
	WebSocketURL string
}

func main() {
	cfg := &run.Config{}
	ParseFlags(cfg)
	SetEnv(cfg)

	err := run.Run(cfg)
	if err != nil {
		panic("Ошибка запуска сервера: " + err.Error())
	}
}
