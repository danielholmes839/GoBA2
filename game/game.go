package game

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goba2/game/netcode"
	"os"
	"strconv"
	"time"
)

type Event struct {
	Code      int             `json:"code"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type User struct {
	id string
}

type Game struct {
	*netcode.Server
	name        string
	counter     int
	connections map[string]netcode.Connection
	tasks       chan func()
}

func NewGame(name string) *Game {
	game := &Game{
		// server
		Server: &netcode.Server{
			ServerMetrics: &netcode.LocalServerMetrics{},
			Name:          name,
			Tasks:         make(chan func()),
		},
		// game
		name:        name,
		counter:     0,
		connections: make(map[string]netcode.Connection),
		tasks:       make(chan func()),
	}
	game.TickFunc = game.tick
	return game
}

func (g *Game) Connect(ctx context.Context, conn netcode.Connection, user *User) error {
	// add the connection
	if err := g.connect(conn, user); err != nil {
		return err
	}

	fName := fmt.Sprintf("./logs/user-%s.txt", user.id)
	os.Remove(fName)
	f, _ := os.Create(fName)

	// open connection
	go conn.Open(context.Background(), f, func() {
		g.disconnect(user)
	})

	return nil
}

func (g *Game) connect(conn netcode.Connection, user *User) error {
	// connect a player
	var err error
	g.Do(func() {
		if _, taken := g.connections[user.id]; taken {
			err = errors.New("user already connected")
			return
		}
		g.connections[user.id] = conn
	})

	if err == nil {
		fmt.Printf("game: %s, user connected: %s\n", g.name, user.id)
	}

	return err
}

func (g *Game) disconnect(user *User) {
	// disconnect a player
	g.Do(func() {
		delete(g.connections, user.id)
	})

	fmt.Printf("game: %s, user disconnected: %s\n", g.name, user.id)
}

func (g *Game) tick() {
	time.Sleep(time.Millisecond * 10)
	g.counter++
	// fmt.Println(g.counter, len(g.connections))

	// send some data
	data := []byte(strconv.Itoa(len(g.connections)))

	for _, connection := range g.connections {
		connection.Write(data)
	}
}
