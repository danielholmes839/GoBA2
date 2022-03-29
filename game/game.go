package game

import (
	"encoding/json"
	"fmt"
	"goba2/realtime"
	"io"
	"time"
)

type User struct {
	Id string
}

func (u User) ID() string {
	return u.Id
}

type UserInfo struct {
	user User
	conn io.Writer
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

func (g *Game) OnMessage(user string, data []byte) {
	fmt.Printf("user: %s message: %s\n", user, string(data))
	g.users[user].conn.Write(data)

}

func (g *Game) OnConnect(user User, conn io.Writer) error {
	// connection succeeded
	fmt.Printf("game: %s new connection id: %s\n", g.name, user.Id)
	g.users[user.Id] = &UserInfo{
		user: user,
		conn: conn,
	}
	return nil
}

func (g *Game) OnDisconnect(user User) {
	// connection disconnected
	fmt.Printf("game: %s closed connection id: %s\n", g.name, user.Id)
	delete(g.users, user.Id)
}

func (g *Game) OnClose() {
	fmt.Println("game: shutdown!")
}

func (g *Game) OnOpen(engine realtime.Scheduler) {
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
		for _, connection := range g.users {
			connection.conn.Write(data)
		}
		counter++
		fmt.Println("1 second interval")
	})

	fmt.Println("game: startup!")
}
