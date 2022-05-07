package backend

import (
	"context"
	"goba2/games/goba2"
	"goba2/realtime"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CORSMiddleware)

	r.Get("/connect", s.GameEndpoint())

	return r
}

func (s *Server) GameEndpoint() http.HandlerFunc {
	game := goba2.NewGame("app")

	server := realtime.NewServer[goba2.User](game, &realtime.Config{
		Room:                realtime.NewLimitRoom(10),
		Metrics:             &realtime.EmptyMetrics{},
		SynchronousMessages: true,
	})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()

	if err := server.Open(ctx); err != nil {
		panic(err)
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		identity := goba2.User{Id: name}

		// upgrade the websocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// add the user to the game
		user := goba2.User{Id: identity.ID()}
		ws := &realtime.Websocket{Conn: conn}

		if err = server.Connect(user, ws); err != nil {
			ws.Close()
		}
	}
}
