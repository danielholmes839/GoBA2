package game

import (
	"context"
	"fmt"
	"goba2/netcode"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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

func (g *Game) OnMessage(user string, data []byte) {
	fmt.Printf("user: %s message: %s\n", user, string(data))
	g.users[user].conn.Write(data)

}

func (g *Game) OnConnect(user User, conn netcode.Connection) error {
	// connection succeeded
	fmt.Printf("game: %s new connection id: %s\n", g.name, user.id)
	g.users[user.id] = &UserInfo{
		user: user,
		conn: conn,
	}
	return nil
}

func (g *Game) OnDisconnect(user User) {
	// connection disconnected
	fmt.Printf("game: %s closed connection id: %s\n", g.name, user.id)
	delete(g.users, user.id)
}

func (g *Game) OnClose() {
	fmt.Println("game: shutdown!")
}

func (g *Game) OnOpen(ctx context.Context, engine netcode.Engine) {
	engine.After(time.Second*3, func() {
		fmt.Println("3 second after (after)")
	})

	engine.At(time.Now().Add(time.Second*3), func() {
		fmt.Println("3 second after (at)")
	})

	engine.Interval(time.Second*10, func() {
		fmt.Println("10 second interval")
	})

	fmt.Println("game: startup!")
}

func GameEndpoint() http.HandlerFunc {

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 30)
		cancel()
	}()

	mygame := NewGame("my-game")

	server := netcode.NewServer[User](
		mygame,
		&netcode.Config{
			Metrics:         &netcode.LocalServerMetrics{},
			ConnectionLimit: 100,
		},
	)

	if err := server.Open(ctx, 64); err != nil {
		panic(err)
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// upgrade the websocket connection
		token := r.URL.Query().Get("token")
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			fmt.Println(err)
			return
		}

		// add the user to the game
		ws := &netcode.Websocket{Conn: conn}

		if err = server.Connect(ctx, User{id: token}, ws); err != nil {
			fmt.Println("connection error:", err)
			ws.Close()
		}
	}
}
