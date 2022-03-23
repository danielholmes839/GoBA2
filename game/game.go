package game

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Code      int             `json:"code"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type player struct {
	name      string
	connected bool
}

type Game struct {
	name    string
	counter int
	players map[string]*player
}

func NewGame(name string) *Game {
	return &Game{
		// game
		name:    name,
		counter: 0,
		players: make(map[string]*player),
	}
}

func (g *Game) Add(id, name string) {
	g.players[id] = &player{
		name:      name,
		connected: false,
	}
}

func (g *Game) Tick() {
	time.Sleep(time.Millisecond * 10)
	g.counter++
}

func (g *Game) OnConnect(id string) {
	// connection succeeded
	fmt.Printf("game: %s new connection id: %s\n", g.name, id)
	g.players[id].connected = true
}

func (g *Game) OnConnectError(id string, err error) {
	// connection failed
	fmt.Printf("game: %s new connection id: %s error: %s\n", g.name, id, err)
	delete(g.players, id)
}

func (g *Game) OnDisconnect(id string) {
	// connection disconnected
	fmt.Printf("game: %s closed connection id: %s\n", g.name, id)
	delete(g.players, id)
}
