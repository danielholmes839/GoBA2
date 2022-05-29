package gameplay

import (
	"context"
	"encoding/json"
	"fmt"
	"goba2/games/goba/gameplay/geometry"
	"goba2/games/goba/util"
	"goba2/realtime"
	"time"
)

// ClientInfo struct
type ClientInfo struct {
	team     *Team
	champion *Champion
	score    *Score
}

// Game struct
type Game struct {
	shutdown context.CancelFunc

	// Game settings
	tps int // "Ticks per second" the number of

	events *realtime.Queue[*util.ClientEvent] // Events received are put in this queue
	global *realtime.Subscription

	// Game variables
	clients     map[string]*ClientInfo
	teams       map[*Team]struct{}
	projectiles map[*Projectile]string

	// Structures
	walls  []*geometry.Rectangle
	bushes []*geometry.Rectangle
}

// NewGame func
func NewGame(tps int, shutdown context.CancelFunc) *Game {
	// events queue
	events := realtime.NewQueue[*util.ClientEvent]()

	g := &Game{
		shutdown:    shutdown,
		tps:         tps,
		events:      events,
		global:      realtime.NewSubscription("global"),
		clients:     make(map[string]*ClientInfo),
		teams:       make(map[*Team]struct{}),
		projectiles: make(map[*Projectile]string),

		// Structures
		walls: []*geometry.Rectangle{
			geometry.NewRectangle(-2000, -2000, 7000, 2000),
			geometry.NewRectangle(-2000, -2000, 2000, 7000),
			geometry.NewRectangle(-2000, 3000, 7000, 2000),
			geometry.NewRectangle(3000, -2000, 2000, 7000),

			geometry.NewRectangle(400, 400, 950, 100),
			geometry.NewRectangle(1650, 400, 950, 100),

			geometry.NewRectangle(400, 2500, 950, 100),
			geometry.NewRectangle(1650, 2500, 950, 100),
		},

		bushes: []*geometry.Rectangle{
			geometry.NewRectangle(500, 500, 500, 300), geometry.NewRectangle(1250, 500, 500, 300), geometry.NewRectangle(2000, 500, 500, 300),
			geometry.NewRectangle(500, 2200, 500, 300), geometry.NewRectangle(1250, 2200, 500, 300), geometry.NewRectangle(2000, 2200, 500, 300),
			geometry.NewRectangle(0, 1000, 300, 1000), geometry.NewRectangle(2700, 1000, 300, 1000),
		},
	}

	g.teams[NewTeam("Red Team", "#ff0000", geometry.NewPoint(1500, 200))] = struct{}{}
	g.teams[NewTeam("Blue Team", "#0000ff", geometry.NewPoint(1500, 2800))] = struct{}{}
	return g
}

func (g *Game) HandleOpen(ctx context.Context, engine realtime.Engine) {
	fmt.Println("game started")
	t := time.Duration(int(time.Second) / g.tps)

	engine.Interval(t, g.tick)

	engine.After(time.Second*3, func() {
		// cancel the context after 3 seconds if no one joins the game
		if len(g.clients) == 0 {
			g.shutdown()
		}
	})
}

func (g *Game) HandleClose() {
	fmt.Println("game stopped")
}

func (g *Game) HandleConnect(identity realtime.ID, conn realtime.Connection) error {
	id := identity.ID()
	champion := NewChampion(id)
	team := g.getNextTeam()
	team.addClient(id, conn, champion)
	g.addClientInfo(id, champion, team)

	// successfully connected
	util.WriteMessage("connection", util.Marshall(struct {
		// joinJSON from GoBA package
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}{
		Success: true,
		Error:   "",
	}), conn)

	// Send the new player setup data
	util.WriteMessage("setup", NewSetupUpdate(g, id), conn)

	// Send clients the updated teams
	g.global.Subscribe(id, conn)
	util.BroadcastMessage("update-teams", NewTeamsUpdate(g), g.global)

	fmt.Printf("'%s' connected\n", id)
	return nil
}

func (g *Game) HandleDisconnect(id string) {
	// Disconnect the client
	g.global.Unsubscribe(id)
	g.getClientTeam(id).removeClient(id) // Remove client from team
	delete(g.clients, id)                // Remove client from game

	// Send clients the updated teams
	util.BroadcastMessage("update-teams", NewTeamsUpdate(g), g.global)
	fmt.Printf("'%s' disconnected\n", id)

	if len(g.clients) == 0 {
		g.shutdown()
	}
}

// Handle func
func (g *Game) HandleMessage(id string, data []byte) {
	event := &util.ClientEvent{}
	if err := json.Unmarshal(data, event); err != nil {
		return
	}

	event.Client = id

	// fmt.Printf("category:'%s' name:'%s' client:'%s'\n", event.GetCategory(), event.GetName(), id)
	switch event.Category {
	case "game":
		g.events.Push(event)
	}
}

// GetPlayerCount func
func (game *Game) GetPlayerCount() int {
	return len(game.clients)
}

// Get the team with lowest number of players
func (game *Game) getNextTeam() *Team {
	var min *Team = nil
	for team := range game.teams {
		if min == nil {
			min = team
		} else if team.size < min.size {
			min = team
		}
	}
	return min
}

func (game *Game) getClientChampion(client string) *Champion {
	// Get the champion of the client
	return game.clients[client].champion
}

func (game *Game) getClientTeam(client string) *Team {
	// Get the team of the client
	return game.clients[client].team
}

func (game *Game) getClientScore(client string) *Score {
	// Get the info (team and champion) of the client
	return game.clients[client].score
}

func (game *Game) getClientInfo(client string) *ClientInfo {
	// Get the info (team and champion) of the client
	return game.clients[client]
}

func (game *Game) addClientInfo(id string, champion *Champion, team *Team) {
	// Set the info (team and champion) of the client
	score := NewScore(0, 0, 0)
	game.clients[id] = &ClientInfo{champion: champion, team: team, score: score}
}
