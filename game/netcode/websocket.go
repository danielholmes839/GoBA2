package netcode

import (
	"context"
	"io"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	*websocket.Conn
	Id string
}

func (ws *Websocket) ID() string {
	return ws.Id
}

func (ws *Websocket) Open(ctx context.Context, handler io.Writer, close Callback) error {
	// create a context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		// websocket closes
		<-ctx.Done()
		ws.Close()
		close()
	}()

	for {
		// read messages
		_, data, err := ws.ReadMessage()
		if err != nil {
			return err
		}
		handler.Write(data)
	}
}

func (ws *Websocket) Write(data []byte) (int, error) {
	// write to the connection implements
	err := ws.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}
