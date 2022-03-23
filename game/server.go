package game

import (
	"context"
	"fmt"
	"goba2/game/netcode"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
}

func (s *Server) GameEndpoint() http.HandlerFunc {
	ctx := context.Background()
	g := NewGame(1)
	go g.Run(ctx, 64)

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
		ws := &netcode.Websocket{Conn: conn}
		err = g.Connect(ctx, ws, &User{id: token})

		if err != nil {
			ws.Close()
			fmt.Println(err)
		}
	}
}
