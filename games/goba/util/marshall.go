package util

import "encoding/json"

func Marshall(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}