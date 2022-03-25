package game

import (
	"context"
	"fmt"
	"goba2/netcode"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
}

func (s *Server) GameEndpoint() http.HandlerFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 30)
		cancel()
	}()

	mygame := NewGame("my-game")

	server := netcode.NewServer[User](mygame, 5).
		WithMetrics(&netcode.LocalServerMetrics{})

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
