package netcode

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Config struct {
	Metrics         ServerMetrics
	ConnectionLimit int
}

type Server[I Identity] struct {
	sync.Mutex
	game            ServerHooks[I]
	metrics         ServerMetrics
	connectionLimit int
	connections     map[string]Connection
	open            bool
	ctx             context.Context
}

func NewServer[I Identity](game ServerHooks[I], conf *Config) *Server[I] {
	return &Server[I]{
		Mutex:           sync.Mutex{},
		game:            game,
		metrics:         conf.Metrics,
		connectionLimit: conf.ConnectionLimit,
		connections:     make(map[string]Connection),
		open:            false,
	}
}

func (s *Server[I]) Connect(ctx context.Context, identity I, conn Connection) error {
	var err error
	s.Lock()
	defer s.Unlock()

	id := identity.ID()

	if !s.open {
		return errors.New("server closed")
	}

	if s.connectionLimit == len(s.connections) {
		return errors.New("server full")
	}

	if _, exists := s.connections[id]; exists {
		return errors.New("already connected")
	}

	if err = s.game.OnConnect(identity, conn); err != nil {
		return err
	}

	s.connections[id] = conn

	go func() {
		for {
			data, err := conn.Receive()
			if err != nil {
				return
			}

			s.Lock()
			s.game.OnMessage(id, data)
			s.Unlock()
		}
	}()

	return err
}

func (s *Server[I]) Open(ctx context.Context, tps int) error {
	// calculate the delay to achieve the correct tps
	s.Lock()
	defer s.Unlock()

	if s.open {
		return errors.New("server already open")
	}

	s.open = true
	s.ctx = ctx

	target := time.Duration(int64(float64(time.Second) / float64(tps)))
	ticker := time.NewTicker(target)

	shutdown := func() {
		s.Lock()
		defer s.Unlock()

		for _, connection := range s.connections {
			connection.Close()
		}

		s.game.OnClose()
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
				s.game.Tick()
				done = time.Now()
				s.Unlock()
				s.metrics.RecordTick(start, ready.Sub(start), done.Sub(ready))
			}
		}
	}()

	s.game.OnOpen(ctx, s)
	return nil
}

func (s *Server[I]) Do(f func()) {
	// measure the wait and execution time of tasks
	var ready, start, done time.Time

	// measure when the waiting started
	ready = time.Now()
	s.Lock()

	// measure when the task started
	start = time.Now()
	f()
	done = time.Now()
	s.Unlock()

	s.metrics.RecordTask(start, start.Sub(ready), done.Sub(start))
}

func (s *Server[I]) After(d time.Duration, f func()) {
	// execute function after a delay
	go func() {
		ctx, cancel := context.WithTimeout(s.ctx, d)
		defer cancel()

		<-ctx.Done()
		if ctx.Err() == context.Canceled {
			return
		}
		s.Do(f)
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
