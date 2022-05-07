package goba2

import (
	"context"
	"encoding/json"
	"fmt"
	"goba2/realtime"
	"io"
	"time"
)

type Connection struct {
	io.WriteCloser
	id string
}

type Game struct {
	counter     int
	connections map[string]*Connection
}

func NewGame() *Game {
	return &Game{
		// game
		counter:     0,
		connections: make(map[string]*Connection),
	}
}

func (g *Game) HandleMessage(id string, data []byte) {
	fmt.Printf("user: %s message: %s\n", id, string(data))

	// unmarshall the event
	event := Event{}
	if err := json.Unmarshal(data, &event); err != nil {
		return
	}

	switch event.Code {
	case 1:
		// unmarshall move
		move := EventMove{}
		if err := json.Unmarshal(event.Data, &move); err != nil {
			return
		}
		fmt.Println(move)
	}
}

func (g *Game) HandleConnect(identity realtime.ID, conn realtime.Connection) error {
	// connection succeeded
	id := identity.ID()
	fmt.Printf("new connection id: %s\n", id)

	g.connections[id] = &Connection{
		WriteCloser: conn,
		id:          id,
	}

	return nil
}

func (g *Game) HandleDisconnect(id string) {
	// connection disconnected
	fmt.Printf("closed connection id: %s\n", id)
	delete(g.connections, id)
}

func (g *Game) HandleClose() {
	fmt.Println("game: shutdown!")
}

func (g *Game) HandleOpen(ctx context.Context, engine realtime.Engine) {
	engine.After(time.Second*3, func() {
		fmt.Println("3 second after (after)")
	})

	engine.At(time.Now().Add(time.Second*3), func() {
		fmt.Println("3 second after (at)")
	})

	counter := 0
	type interval struct {
		Counter int `json:"counter"`
	}

	engine.Interval(time.Millisecond*1000, func() {
		data, _ := json.Marshal(interval{Counter: counter})
		for _, connection := range g.connections {
			connection.Write(data)
		}
		counter++
		fmt.Println("1 second interval")
	})

	fmt.Println("game: startup!")
}
