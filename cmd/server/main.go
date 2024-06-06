package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

type PageData struct {
	WebSocketURL string
}

func main() {
	fileServer()
}

func fileServer() {
	contentDir := filepath.Join("..", "..", "internal", "content")
	fs := http.FileServer(http.Dir(contentDir))
	http.Handle("/content/", http.StripPrefix("/content/", fs))

	// шаблон для передачи адреса WebSocket соединени
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmplPath := filepath.Join(contentDir, "index.html")
		tmpl := template.Must(template.ParseFiles(tmplPath))

		data := PageData{
			WebSocketURL: "ws://localhost:3001/ws", // todo задавать базовый адрес из командной строки
		}
		tmpl.Execute(w, data)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
