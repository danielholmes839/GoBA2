package backend

import (
	"context"
	"fmt"
	"goba2/backend/auth"
	"goba2/game"
	"goba2/netcode"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type Config struct {
	TokenVerifier auth.TokenVerifier
	Discord       *auth.OAuth2Config
}

type API struct {
	Conf *Config
}

func NewAPI(conf *Config) *API {
	return &API{
		Conf: conf,
	}
}

func (api *API) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CORSMiddleware)

	discordRedirect, discordCallback := auth.OAuth2Endpoints(api.Conf.Discord)

	r.Get("/auth/discord", discordRedirect)
	r.Get("/auth/discord/callback", discordCallback)
	r.Get("/connect", api.GameEndpoint())

	r.Route("/", func(r chi.Router) {
		r.Use(AuthenticationMiddleware(api.Conf.TokenVerifier))
		r.Get("/@me", api.MeEndpoint())
	})

	return r
}

func (api *API) MeEndpoint() http.HandlerFunc {
	return Authenticated(func(w http.ResponseWriter, r *http.Request, identity *auth.Identity) {
		writeJSON(w, http.StatusOK, identity)
	})
}

func (api *API) GameEndpoint() http.HandlerFunc {
	mygame := game.NewGame("my-game")

	server := netcode.NewServer[game.User](
		mygame,
		&netcode.Config{
			Metrics:             &netcode.LocalServerMetrics{},
			ConnectionLimit:     100,
			SynchronousMessages: true,
		},
	)

	ctx := context.Background()

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
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		identity, err := auth.GetIdentity(api.Conf.TokenVerifier, r)
		if err != nil {
			return
		}

		// add the user to the game
		ws := &netcode.Websocket{Conn: conn}

		if err = server.Connect(ctx, game.User{Id: identity.ID()}, ws); err != nil {
			fmt.Println("connection error:", err)
			ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
			ws.Close()
		}

	}
}
