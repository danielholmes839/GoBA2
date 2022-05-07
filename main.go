package main

import (
	"goba2/backend"
	"net/http"
)

func main() {
	router := backend.NewServer().Routes()
	http.ListenAndServe("localhost:3000", router)
}
