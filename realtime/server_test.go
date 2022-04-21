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
			Metrics:             &EmptyMetrics{},
			ConnectionLimit:     1,
			SynchronousMessages: true,
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
	server := NewServer[mockid](
		game,
		&Config{
			Metrics:             &EmptyMetrics{},
			ConnectionLimit:     1,
			SynchronousMessages: true,
		})

	t.Run("server closed", func(t *testing.T) {
		err := server.Connect(mockid{"1"}, &mockconn{conn: make(chan []byte)})
		assert.ErrorIs(t, err, ErrServerClosed)
	})

	ctx, cancel := context.WithCancel(context.Background())
	server.Open(ctx)

	t.Run("ok", func(t *testing.T) {
		// ok
		err := server.Connect(mockid{"1"}, &mockconn{conn: make(chan []byte)})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(server.connections))
	})

	t.Run("server full", func(t *testing.T) {
		// connection limit is 1 so the server is already full
		err := server.Connect(mockid{"1"}, &mockconn{conn: make(chan []byte)})
		assert.ErrorIs(t, err, ErrServerFull)
	})

	t.Run("already connected", func(t *testing.T) {
		// attempt to connect twice with the same id
		server.connectionLimit = 2
		err := server.Connect(mockid{"1"}, &mockconn{conn: make(chan []byte)})
		assert.ErrorIs(t, err, ErrAlreadyConnected)
	})

	t.Run("game.OnConnect error", func(t *testing.T) {
		game.OnConnectErr = errors.New("game connection rejected")
		err := server.Connect(mockid{"2"}, &mockconn{conn: make(chan []byte)})
		assert.ErrorIs(t, err, game.OnConnectErr)
	})

	cancel()
	time.Sleep(time.Millisecond*5)
	assert.Equal(t, 0, len(server.connections))

}

func TestInterval(t *testing.T) {
	game := &mockgame{}
	server := NewServer[mockid](
		game,
		&Config{
			Metrics:             &EmptyMetrics{},
			ConnectionLimit:     1,
			SynchronousMessages: true,
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
			ConnectionLimit:     1,
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
