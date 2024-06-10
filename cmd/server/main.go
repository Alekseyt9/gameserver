package main

import "gameserver/internal/run"

type PageData struct {
	WebSocketURL string
}

func main() {
	cfg := &run.Config{}
	run.Run(cfg)
}
