package main

import (
	"goba2/backend"
	"net/http"
)

func main() {
	server := backend.NewServer()
	router := server.Routes()
	http.ListenAndServe("localhost:3000", router)
}
