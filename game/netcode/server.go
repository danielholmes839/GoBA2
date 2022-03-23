package netcode

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

type ServerMetrics interface {
	RecordTask(start time.Time, wait, execution time.Duration)
	RecordTick(start time.Time, wait, execution time.Duration)
}

type HasID interface {
	ID() string
}

type ServerHooks[R HasID] interface {
	Tick()
	OnConnect(request R, conn Connection) error
	OnDisconnect(request R)
	OnShutdown()
	OnStartup()
}

type Server[R HasID] struct {
	sync.Mutex
	ServerHooks[R]
	ServerMetrics
	Name             string
	CONNECTION_LIMIT int
	Connections      map[string]Connection
	open             bool
	ctx              context.Context
}

func (s *Server[R]) Connect(ctx context.Context, request R, conn Connection, handler io.Writer) error {
	var err error

	s.Do(func() {
		id := request.ID()
		if !s.open {
			err = errors.New("server closed")
			return
		}

		if s.CONNECTION_LIMIT == len(s.Connections) {
			err = errors.New("server full")
			return
		}

		if _, found := s.Connections[id]; found {
			err = errors.New("connection id already exists")
			return
		}

		if err = s.OnConnect(request, conn); err != nil {
			return
		}

		s.Connections[id] = conn

		go conn.Open(ctx, handler, func() {
			// disconnect the client
			s.Do(func() {
				if s.open {
					s.OnDisconnect(request)
					delete(s.Connections, id)
				}
			})
		})
	})

	return err
}

func (s *Server[R]) Open(ctx context.Context, tps int) error {
	// calculate the delay to achieve the correct tps
	s.Lock()
	defer s.Unlock()

	if s.open {
		return errors.New("server already open")
	}

	s.open = true

	target := time.Duration(int64(float64(time.Second) / float64(tps)))
	ticker := time.NewTicker(target)

	shutdown := func() {
		s.Lock()
		defer s.Unlock()

		for _, connection := range s.Connections {
			connection.Close()
		}

		s.OnShutdown()
		s.open = false
	}

	go func() {
		var ready, start, done time.Time
		for {
			select {
			case <-ctx.Done():
				shutdown()
				return

			case <-ticker.C:
				// execute ticks
				start = time.Now()
				s.Lock()
				ready = time.Now()
				s.Tick()
				done = time.Now()
				s.Unlock()
				s.RecordTick(start, start.Sub(ready), done.Sub(ready))
			}
		}
	}()

	s.OnStartup()
	return nil
}

func (s *Server[R]) Do(task func()) {
	// measure the wait and execution time of tasks
	var ready, start, done time.Time

	// measure when the waiting started
	ready = time.Now()
	s.Lock()

	// measure when the task started
	start = time.Now()
	task()
	done = time.Now()
	s.Unlock()

	s.RecordTask(start, start.Sub(ready), done.Sub(ready))
}
