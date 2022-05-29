package util

import (
	"encoding/json"
	"time"
)

// ClientEvent type
type ClientEvent struct {
	Client    string
	Category  string          `json:"category"`
	Name      string          `json:"event"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// NewClientEvent func
func NewClientEvent(data []byte) (*ClientEvent, error) {
	event := &ClientEvent{}
	err := json.Unmarshal(data, event)
	return event, err
}

// GetCategory func
func (event *ClientEvent) GetCategory() string {
	return event.Category
}

// GetName func
func (event *ClientEvent) GetName() string {
	return event.Name
}

// GetTimestamp func
func (event *ClientEvent) GetTimestamp() int64 {
	return event.Timestamp
}

// GetData func
func (event *ClientEvent) GetData() []byte {
	return event.Data
}

// ServerEvent type
type ServerEvent struct {
	Subscription string          `json:"subscription"`
	Name         string          `json:"name"`
	Timestamp    int64           `json:"timestamp"`
	Data         json.RawMessage `json:"data"`
}

// NewServerEvent func
func NewServerEvent(subscription string, name string, data []byte) *ServerEvent {
	return &ServerEvent{Subscription: subscription, Name: name, Timestamp: time.Now().Unix(), Data: data}
}

// Serialize func
func (serverEvent *ServerEvent) Serialize() []byte {
	data, _ := json.Marshal(serverEvent)
	return data
}
