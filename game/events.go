package game

import (
	"encoding/json"
	"time"
)

type Event struct {
	Code      int             `json:"code"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}
