package game

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

type Event struct {
	Code      int             `json:"code"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type EventQueue interface {
	io.Writer
	Read() (string, Event)
}

type GameEventQueue struct {
	sync.Mutex
	Sender string
	Events []Event
}

func (gq *GameEventQueue) Write(p []byte) (n int, err error) {
	event := &Event{}
	if err := json.Unmarshal(p, event); err != nil {
		return 0, err
	}
	gq.Lock()
	gq.Events = append(gq.Events, *event)
	gq.Unlock()
	return len(p), nil
}
