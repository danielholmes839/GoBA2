package realtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerOpen(t *testing.T) {
	server := NewServer[mockid](
		&mockgame{},
		&Config{
			Metrics: &EmptyMetrics{},
			Room:    NewBasicRoom(1),
		})

	ctx, cancel := context.WithCancel(context.Background())

	err := server.Open(ctx)
	assert.NoError(t, err)
	assert.True(t, server.open)

	err = server.Open(ctx)
	assert.Error(t, err)
	assert.True(t, server.open)

	id := mockid{"1"}
	conn := &mockconn{conn: make(chan []byte)}

	err = server.Connect(id, conn)
	assert.NoError(t, err)

	cancel()
	time.Sleep(time.Millisecond * 5)
	assert.False(t, server.open)
}

func TestServerConnect(t *testing.T) {
	game := &mockgame{}
	room := NewBasicRoom(1)
	server := NewServer[mockid](
		game,
		&Config{
			Metrics:             &EmptyMetrics{},
			Room:                room,
			SynchronousMessages: true,
		})

	t.Run("server closed", func(t *testing.T) {
		// server is closed
		err := server.Connect(mockid{"1"}, &mockconn{conn: make(chan []byte)})
		assert.ErrorIs(t, err, ErrServerClosed)
	})

	// open the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server.Open(ctx)

	t.Run("ok", func(t *testing.T) {
		// ok
		err := server.Connect(mockid{"1"}, &mockconn{conn: make(chan []byte)})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(room.connections))
	})

	t.Run("server full", func(t *testing.T) {
		// connection limit is 1 so the server is already full
		err := server.Connect(mockid{"1"}, &mockconn{conn: make(chan []byte)})
		assert.ErrorIs(t, err, ErrRoomFull)
	})

	t.Run("game.OnConnect error", func(t *testing.T) {
		game.OnConnectErr = errors.New("game connection rejected")
		room.connectionLimit++
		err := server.Connect(mockid{"2"}, &mockconn{conn: make(chan []byte)})
		assert.ErrorIs(t, err, game.OnConnectErr)
	})
}

func TestServerAfter(t *testing.T) {
	server := NewServer[mockid](&mockgame{}, &Config{
		Metrics: &EmptyMetrics{},
		Room:    NewBasicRoom(1),
	})

	ctx, cancel := context.WithCancel(context.Background())
	server.Open(ctx)

	counter := 0

	increment := func() {
		counter++
	}
	// execute an increment after 3 milliseconds
	server.After(time.Millisecond*3, increment)

	// after 2 milliseconds counter not incremented
	time.Sleep(time.Millisecond * 2)
	assert.Equal(t, 0, counter)

	// after 4 milliseconds total counter incremented
	time.Sleep(time.Millisecond * 2)
	assert.Equal(t, 1, counter)

	// execute another increment after 3 milliseconds
	server.After(time.Millisecond*3, increment)

	// cancel should stop the increment
	cancel()

	// after 4 milliseconds the counter is still 1
	time.Sleep(time.Millisecond * 4)
	assert.Equal(t, 1, counter)

	// execute another increment when the context is cancelled
	server.After(time.Millisecond*1, func() {
		counter++
	})

	// the counter is still 1
	time.Sleep(time.Millisecond * 2)
	assert.Equal(t, 1, counter)
}

func TestServerAt(t *testing.T) {
	server := NewServer[mockid](&mockgame{}, &Config{
		Metrics: &EmptyMetrics{},
		Room:    NewBasicRoom(1),
	})
	ctx, cancel := context.WithCancel(context.Background())
	server.Open(ctx)

	time.Now().Add(time.Millisecond * 1)
	cancel()
}
func TestInterval(t *testing.T) {
	game := &mockgame{}
	server := NewServer[mockid](
		game,
		&Config{
			Metrics: &EmptyMetrics{},
			Room:    NewBasicRoom(1),
		})

	counter := 0
	ctx, cancel := context.WithCancel(context.Background())

	server.Open(ctx)

	server.Interval(time.Millisecond*5, func() {
		counter++
		if counter == 3 {
			cancel()
		}
	})

	<-ctx.Done()

	assert.Equal(t, 3, counter)
}

func BenchmarkServer(b *testing.B) {
	server := NewServer[mockid](
		&mockgame{},
		&Config{
			Metrics:             &EmptyMetrics{},
			Room:                NewBasicRoom(1),
			SynchronousMessages: true,
		})

	b.Run("do", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			server.Do(func() {
				return
			})
		}
	})

	b.Run("func", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			func() {
				return
			}()
		}
	})
}
