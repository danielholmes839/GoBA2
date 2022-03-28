package backend

import (
	"encoding/json"
	"net/http"
)

func write(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	data, _ := json.Marshal(v)
	w.WriteHeader(status)
	w.Write(data)
}
