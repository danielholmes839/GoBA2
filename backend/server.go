package backend

import (
	"goba2/games/goba"
	"goba2/realtime"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (s *Server) Router() *mux.Router {
	gobaEndpoints := &goba.Endpoints{
		TPS:         64,
		PlayerLimit: 10,
		Timeout:     time.Minute * 10,
		Instances:   make(map[string]*realtime.Server[realtime.ID]),
	}

	r := mux.NewRouter()
	r.HandleFunc("/goba/v1/create", gobaEndpoints.CreateEndpoint())
	r.HandleFunc("/goba/v1/join", gobaEndpoints.JoinEndpoint())
	r.HandleFunc("/goba/v1/info", gobaEndpoints.InfoEndpoint())

	// files
	fs := http.FileServer(http.Dir("./games/goba/client"))
	r.PathPrefix("/client/").Handler(
		http.StripPrefix("/client/", fs),
	)

	// homepage
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./games/goba/client/index.html")
	})

	return r
}
