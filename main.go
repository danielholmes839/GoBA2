package main

import (
	"goba2/game"
	"net/http"
)

func main() {
	server := game.Server{}
	http.HandleFunc("/game/connect", server.GameEndpoint())
	http.ListenAndServe("localhost:3000", nil)
}
