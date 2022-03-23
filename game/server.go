package game

import (
	"context"
	"fmt"
	"goba2/game/netcode"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
}

func (s *Server) GameEndpoint() http.HandlerFunc {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*600)
	mygame := NewGame("my-game")

	server := netcode.Server{
		ServerHooks:     mygame,
		ServerMetrics:   &netcode.LocalServerMetrics{},
		Name:            "my-server",
		MAX_CONNECTIONS: 5,
		Connections:     make(map[string]netcode.Connection),
		Tasks:           make(chan func()),
	}

	go server.Open(ctx, 64)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// upgrade the websocket connection
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			fmt.Println(err)
			return
		}

		// add the user to the game
		token := r.URL.Query().Get("token")
		ws := &netcode.Websocket{Id: token, Conn: conn}
		err = server.Connect(ctx, ws, ws)

		if err != nil {
			ws.Close()
			fmt.Println(err)
		}
	}
}
