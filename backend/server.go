package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"goba2/games/goba"
	"goba2/games/goba/util"
	"goba2/realtime"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type createJSON struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type joinJSON struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func marshall(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}

type Server struct {
	instances map[string]*realtime.Server[realtime.ID]
}

func NewServer() *Server {
	return &Server{
		instances: make(map[string]*realtime.Server[realtime.ID]),
	}
}

func (s *Server) Routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/create", s.GobaCreateEndpoint())
	r.HandleFunc("/join", s.GobaJoinEndpoint())

	// fileserver
	fs := http.FileServer(http.Dir("./client"))
	r.PathPrefix("/client/").Handler(
		http.StripPrefix("/client/", fs),
	)

	// homepage
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/index.html")
	})

	return r
}

func (s *Server) GobaCreateEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), time.Minute*10)
		instance := realtime.NewServer[realtime.ID](
			goba.New(64),
			&realtime.Config{
				Room:                realtime.NewBasicRoom(10),
				Metrics:             &realtime.EmptyMetrics{},
				SynchronousMessages: false,
			},
		)
		_ = instance.Open(ctx)
		code := "ABCD"
		s.instances[code] = instance

		go func() {
			<-ctx.Done()
			delete(s.instances, code)
		}()

		w.WriteHeader(http.StatusOK)
		w.Write(marshall(createJSON{
			Code:    code,
			Success: true,
			Error:   "",
		}))

	}
}

func (s *Server) GobaJoinEndpoint() http.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("1")
		query := r.URL.Query()
		identity := realtime.ID(query.Get("name"))
		code := strings.ToUpper(query.Get("code"))

		// upgrade the websocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic("websocket connection failed")
		}

		instance, exists := s.instances[code]
		if !exists {
			panic("code not found")
		}

		// add the user to the game
		ws := &realtime.Websocket{Conn: conn}

		util.WriteMessage("connection", marshall(&joinJSON{
			Success: true,
			Error:   "",
		}), ws)

		if err = instance.Connect(identity, ws); err != nil {
			ws.Close()
		}

		fmt.Println("4")
	}
}
