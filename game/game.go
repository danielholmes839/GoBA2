package game

import (
	"encoding/json"
	"fmt"
	"goba2/game/netcode"
	"time"
)

type Event struct {
	Code      int             `json:"code"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type Game struct {
	name        string
	counter     int
	connections map[string]netcode.Connection
	tasks       chan func()
}

func NewGame(name string) *Game {
	return &Game{
		// game
		name:    name,
		counter: 0,
	}
}

func (g *Game) Tick() {
	time.Sleep(time.Millisecond * 10)
	g.counter++
}

func (g *Game) OnConnectOK(id string) {
	fmt.Printf("game: %s new connection id: %s\n", g.name, id)
}

func (g *Game) OnConnectError(id string, err error) {
	fmt.Printf("game: %s new connection id: %s error: %s\n", g.name, id, err)
}

func (g *Game) OnDisconnect(id string) {
	fmt.Printf("game: %s closed connection id: %s\n", g.name, id)
}
