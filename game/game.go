package game

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goba2/game/netcode"
	"os"
	"strconv"
	"time"
)

type Event struct {
	Code      int             `json:"code"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type User struct {
	id string
}

type Game struct {
	netcode.Server

	id          int
	counter     int
	connections map[string]netcode.Connection
	tasks       chan func()
}

func NewGame(id int) *Game {
	game := &Game{
		id:          id,
		counter:     0,
		connections: make(map[string]netcode.Connection),
		tasks:       make(chan func()),
	}
	return game
}

func (g *Game) Connect(ctx context.Context, conn netcode.Connection, user *User) error {
	// add the connection
	if err := g.connect(conn, user); err != nil {
		return err
	}

	fName := fmt.Sprintf("./logs/user-%s.txt", user.id)
	os.Remove(fName)
	f, _ := os.Create(fName)

	// open connection
	go conn.Open(context.Background(), f, func() {
		g.disconnect(user)
	})

	return nil
}

func (g *Game) Run(ctx context.Context, tps int) {
	// create the ticker
	target := time.Duration(int64(float64(time.Second) / float64(tps)))
	ticker := time.NewTicker(target)

	metricInterval := time.Duration(5)
	metrics := time.NewTicker(time.Second * metricInterval)

	// stop the ticker
	defer func() {
		ticker.Stop()
		fmt.Printf("game: %d, stopped\n", g.id)
	}()

	fmt.Printf("game: %d, started\n", g.id)

	ticks := 0
	// last := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// diff := time.Since(last)
			// last = time.Now()
			// fmt.Println(diff)
			ticks++
			g.tick()

		case <-metrics.C:
			fmt.Printf("TPS: %d\n", int(ticks/int(metricInterval)))
			ticks = 0

		case task := <-g.tasks:
			task()
		}
	}
}

func (g *Game) connect(conn netcode.Connection, user *User) error {
	// connect a player
	var err error
	g.Do(func() {
		if _, taken := g.connections[user.id]; taken {
			err = errors.New("user already connected")
			return
		}
		g.connections[user.id] = conn
	})

	if err == nil {
		fmt.Printf("game: %d, user connected: %s\n", g.id, user.id)
	}

	return err
}

func (g *Game) disconnect(user *User) {
	// disconnect a player
	g.Do(func() {
		delete(g.connections, user.id)
	})

	fmt.Printf("game: %d, user disconnected: %s\n", g.id, user.id)
}

func (g *Game) tick() {
	g.counter++
	// fmt.Println(g.counter, len(g.connections))

	// send some data
	data := []byte(strconv.Itoa(len(g.connections)))

	for _, connection := range g.connections {
		connection.Write(data)
	}
}
