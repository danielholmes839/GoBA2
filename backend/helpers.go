package backend

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	data, _ := json.Marshal(v)
	w.WriteHeader(status)
	w.Write(data)
}
