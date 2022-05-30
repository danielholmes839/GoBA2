package goba

import (
	"context"
	"fmt"
	"goba2/games/goba/util"
	"goba2/realtime"
	"net/http"
	"strings"
	"time"

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

type infoJSON struct {
	LiveGames int `json:"games"`
}

type Endpoints struct {
	TPS         int
	PlayerLimit int // per game
	Timeout     time.Duration
	Instances   map[string]*realtime.Server[realtime.ID]
}

func (e *Endpoints) CreateEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create a new GoBA instance
		ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
		instance := realtime.NewServer[realtime.ID](
			NewGoba(e.TPS, cancel),
			&realtime.Config{
				Room:                realtime.NewBasicRoom(e.PlayerLimit),
				Metrics:             &realtime.EmptyMetrics{},
				SynchronousMessages: false,
			},
		)
		_ = instance.Open(ctx)

		code := util.Code()
		e.Instances[code] = instance


		go func() {
			// delete the instance when the context ends
			<-ctx.Done()
			delete(e.Instances, code)
		}()
		
		// 
		w.Write(util.Marshall(&createJSON{
			Code:    code,
			Success: true,
			Error:   "",
		}))

	}
}

func (e *Endpoints) JoinEndpoint() http.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// extract query params
		id := r.URL.Query().Get("name")
		code := strings.ToUpper(r.URL.Query().Get("code"))

		// create the websocket connection
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		identity := realtime.ID(id)
		conn := &realtime.Websocket{Conn: ws}

		instance := e.Instances[code]

		// validate the code
		if instance == nil {
			util.WriteMessage("connection", util.Marshall(&joinJSON{
				Success: false,
				Error:   fmt.Sprintf("'%s' does not exist", code),
			}), conn)
			conn.Close()
			return
		}

		// validate the username
		if err = validateName(id); err != nil {
			util.WriteMessage("connection", util.Marshall(&joinJSON{
				Success: false,
				Error:   err.Error(),
			}), conn)
			conn.Close()
			return
		}

		// validate successfully connected to the game
		if err = instance.Connect(identity, conn); err != nil {
			message := err.Error()

			switch err {
			case realtime.ErrRoomFull:
				message = fmt.Sprintf("Connection limit reached (%d)", e.PlayerLimit)
				break
			case realtime.ErrAlreadyConnected:
				message = fmt.Sprintf("The username '%s' is already taken", id)
				break
			}

			util.WriteMessage("connection", util.Marshall(&joinJSON{
				Success: false,
				Error:   message,
			}), conn)
			conn.Close()
			return
		}
	}
}

func (e *Endpoints) joinError() {

}

func (e *Endpoints) InfoEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(util.Marshall(&infoJSON{
			LiveGames: len(e.Instances),
		}))
	}
}
