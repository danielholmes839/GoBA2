package main

import (
	"fmt"
	"goba2/backend"
	"net/http"
)

func main() {
	router := backend.NewServer().Routes()
	fmt.Println("server started...")
	http.ListenAndServe("localhost:3000", router)
}
