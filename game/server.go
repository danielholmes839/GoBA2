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
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	mygame := NewGame("my-game")

	server := &netcode.Server[User]{
		ServerHooks:      mygame,
		ServerMetrics:    &netcode.LocalServerMetrics{},
		Name:             "my-server",
		CONNECTION_LIMIT: 5,
		Connections:      make(map[string]netcode.Connection),
	}

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
		err = server.Connect(ctx, User{id: token}, ws, ws)

		if err != nil {
			ws.Close()
			fmt.Println(err)
		}
	}
}
