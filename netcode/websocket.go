package netcode

import (
	"io"

	"github.com/gorilla/websocket"
)

type Connection interface {
	io.Writer
	io.Closer
	Receive() ([]byte, error)
}

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
	err := ws.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}
