package backend

import (
	"context"
	"goba2/backend/auth"
	"goba2/game"
	"goba2/realtime"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type Config struct {
	TokenVerifier auth.TokenVerifier
	Discord       *auth.OAuth2Config
}

type Server struct {
	Conf *Config
}

func NewServer(conf *Config) *Server {
	return &Server{
		Conf: conf,
	}
}

func (s *Server) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CORSMiddleware)

	discordRedirect, discordCallback := auth.OAuth2Endpoints(s.Conf.Discord)

	r.Get("/auth/discord", discordRedirect)
	r.Get("/auth/discord/callback", discordCallback)
	r.Get("/connect", s.GameEndpoint())

	r.Route("/", func(r chi.Router) {
		r.Use(AuthenticationMiddleware(s.Conf.TokenVerifier))
		r.Get("/@me", s.MeEndpoint())
	})

	return r
}

func (s *Server) MeEndpoint() http.HandlerFunc {
	return AuthHandler(func(w http.ResponseWriter, r *http.Request, identity *auth.Identity) {
		writeJSON(w, http.StatusOK, identity)
	})
}

func (s *Server) GameEndpoint() http.HandlerFunc {
	app := game.NewGame("app")

	server := realtime.NewServer[game.User](app, &realtime.Config{
		Room:                realtime.NewRoom(10),
		Metrics:             &realtime.EmptyMetrics{},
		SynchronousMessages: true,
	})

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
		// identity, err := auth.GetIdentity(s.Conf.TokenVerifier, r)
		// if err != nil {
		// 	return
		// }

		name := r.URL.Query().Get("name")
		identity := game.User{Id: name}

		// upgrade the websocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// add the user to the game
		user := game.User{Id: identity.ID()}
		ws := &realtime.Websocket{Conn: conn}

		if err = server.Connect(user, ws); err != nil {
			ws.Close()
		}
	}
}
