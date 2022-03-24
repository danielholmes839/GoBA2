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

type User struct {
	id string
}

func (u User) ID() string {
	return u.id
}

type UserInfo struct {
	user User
	conn netcode.Connection
}

type Game struct {
	name    string
	counter int
	users   map[string]*UserInfo
}

func NewGame(name string) *Game {
	return &Game{
		// game
		name:    name,
		counter: 0,
		users:   make(map[string]*UserInfo),
	}
}

func (g *Game) Tick() {
	time.Sleep(time.Millisecond * 10)
	g.counter++
}

func (g *Game) OnConnect(user User, conn netcode.Connection) error {
	// connection succeeded
	fmt.Printf("game: %s new connection id: %s\n", g.name, user.id)
	return nil
}

func (g *Game) OnDisconnect(user User) {
	// connection disconnected
	fmt.Printf("game: %s closed connection id: %s\n", g.name, user.id)
	delete(g.users, user.id)
}

func (g *Game) OnShutdown() {
	fmt.Println("game: shutdown!")
}

func (g *Game) OnStartup(engine netcode.Engine) {
	engine.After(time.Second*3, func() {
		fmt.Println("3 second after (after)")
	})

	engine.At(time.Now().Add(time.Second*3), func() {
		fmt.Println("3 second after (at)")
	})

	engine.Interval(time.Second, func() {
		fmt.Println("my interval")
	})

	fmt.Println("game: startup!")
}
