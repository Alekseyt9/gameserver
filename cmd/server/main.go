package main

import services "gameserver/internal"

type PageData struct {
	WebSocketURL string
}

func main() {
	cfg := &services.Config{}
	services.Run(cfg)
}
