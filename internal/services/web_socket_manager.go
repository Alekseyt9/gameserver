package services

import (
	"net/http"

	"github.com/olahol/melody"
)

type WebSocketManager struct {
}

func (*WebSocketManager) Create(url string) {
	m := melody.New()

	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "index.html")
		})
	*/

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		//m.Broadcast(msg)
	})

	http.ListenAndServe(":5000", nil)
}
