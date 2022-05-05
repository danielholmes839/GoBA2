package realtime

import (
	"github.com/gorilla/websocket"
)

const ServerShuttingDown = 4000

type Websocket struct {
	*websocket.Conn
}

func (ws *Websocket) Receive() ([]byte, error) {
	_, data, err := ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ws *Websocket) Write(data []byte) (int, error) {
	// write to the connection implements
	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		return 0, err
	}
	return len(data), nil
}
