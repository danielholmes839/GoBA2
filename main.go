package main

import (
	"goba2/backend"
	"net/http"
)

func main() {
	router := backend.New().Router()
	http.ListenAndServe(":3000", router)
}
