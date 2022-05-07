package backend

import (
	"context"
	"goba2/games/goba2"
	"goba2/realtime"
	"net/http"

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
	server := realtime.NewServer[realtime.ID](
		goba2.NewGame(), 
		&realtime.Config{
			Room:                realtime.NewBasicRoom(10),
			Metrics:             &realtime.EmptyMetrics{},
			SynchronousMessages: true,
		},
	)

	ctx := context.Background()

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
		identity := realtime.ID(name)

		// upgrade the websocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// add the user to the game
		ws := &realtime.Websocket{Conn: conn}

		if err = server.Connect(identity, ws); err != nil {
			ws.Close()
		}
	}
}
