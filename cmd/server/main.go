package main

import (
	"fmt"
	"net/http"
)

func main() {
	fileServer()
}

func fileServer() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
