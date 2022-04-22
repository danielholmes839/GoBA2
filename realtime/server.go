package realtime

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrServerClosed = errors.New("server closed")
var ErrServerAlreadyOpen = errors.New("server already open")

type Config struct {
	Room                *Room
	Metrics             Metrics
	SynchronousMessages bool
}

type Server[I Identity] struct {
	sync.Mutex
	ctx                 context.Context
	app                 Application[I]
	metrics             Metrics
	room                *Room
	synchronousMessages bool
	open                bool
}

func NewServer[I Identity](app Application[I], conf *Config) *Server[I] {
	return &Server[I]{
		Mutex:               sync.Mutex{},
		app:                 app,
		room:                conf.Room,
		metrics:             conf.Metrics,
		synchronousMessages: conf.SynchronousMessages,
		open:                false,
	}
}

func (s *Server[I]) Connect(identity I, conn Connection) error {
	s.Lock()
	defer s.Unlock()

	// check that the server is open
	if !s.open {
		return ErrServerClosed
	}

	id := identity.ID()

	// check that the id is not already connected or the server is full
	if err := s.room.Connect(id, conn); err != nil {
		return err
	}

	// check that the identity successfully connected to the game
	if err := s.app.HandleConnect(identity, conn); err != nil {
		return err
	}

	// ctx to close the connection / disconnect the player
	ctx, cancel := context.WithCancel(s.ctx)

	go func() {
		<-ctx.Done()
		s.Lock()
		defer s.Unlock()

		// disconnect
		s.room.Disconnect(id)
		s.app.HandleDisconnect(id)
	}()

	go func() {
		defer cancel()
		// start processing messages
		if s.synchronousMessages {
			s.processMessages(id, conn, s.processMessagesSynchronously)
		} else {
			s.processMessages(id, conn, s.processMessagesAsynchronously)
		}
	}()

	return nil
}

func (s *Server[I]) Open(ctx context.Context) error {
	s.Lock()
	defer s.Unlock()

	// check the server isn't already open
	if s.open {
		return ErrServerAlreadyOpen
	}

	go func() {
		<-ctx.Done()
		s.Lock()
		defer s.Unlock()
		s.app.HandleClose()
		s.open = false
	}()

	s.open = true
	s.ctx = ctx
	s.app.HandleOpen(s.ctx, s)
	return nil
}

func (s *Server[I]) Do(f func()) {
	// measure the wait and execution time of tasks
	var ready, start, done time.Time

	ready = time.Now() // measure when the waiting started
	s.Lock()           // acquire the lock
	start = time.Now() // measure when the task started
	f()                // execute the function
	done = time.Now()  // measure when the task ended
	s.Unlock()         // release the lock

	s.metrics.RecordLockContention(start, start.Sub(ready), done.Sub(start))
}

func (s *Server[I]) After(d time.Duration, f func()) {
	// execute function after a delay
	go func() {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(d):
			s.Do(f)
		}
	}()
}

func (s *Server[I]) At(t time.Time, f func()) {
	// execute function at a specific time
	if now := time.Now(); now.Before(t) {
		s.After(t.Sub(now), f)
	}
}

func (s *Server[I]) Interval(d time.Duration, f func()) {
	// execute function on an interval
	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				s.Do(f)
			}
		}
	}()
}

/* processMessages from a connection until there's an error (connection ends)*/
func (s *Server[I]) processMessages(id string, conn Connection, process func(id string, data []byte)) {
	for {
		data, err := conn.Receive()
		if err != nil {
			return
		}
		process(id, data)
	}
}

/* processMessagesAsynchronously
call game.OnMessage without acquiring the lock. better for real-time games
better performance since messages can be deserialized without acquiring the lock
*/
func (s *Server[I]) processMessagesAsynchronously(id string, data []byte) {
	s.app.HandleMessage(id, data)
}

/* processMessagesSynchronously
the lock must be acquired before calling game.OnMessage(). better for games that process a lower number of events
*/
func (s *Server[I]) processMessagesSynchronously(id string, data []byte) {
	s.Lock()
	s.app.HandleMessage(id, data)
	s.Unlock()
}
