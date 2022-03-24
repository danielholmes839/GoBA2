package netcode

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

type ServerEngine struct {
	sync.Mutex
	ctx context.Context
	do  func()
}

func NewServer[T Token](game ServerHooks[T], connectionLimit int) *Server[T] {
	return &Server[T]{
		Mutex:           sync.Mutex{},
		game:            game,
		metrics:         &EmptyMetrics{},
		connectionLimit: connectionLimit,
		connections:     make(map[string]Connection),
		open:            false,
	}
}

type Server[T Token] struct {
	sync.Mutex
	game            ServerHooks[T]
	metrics         ServerMetrics
	connectionLimit int
	connections     map[string]Connection
	open            bool
	ctx             context.Context
}

func (s *Server[T]) WithMetrics(metrics ServerMetrics) *Server[T] {
	s.metrics = metrics
	return s
}

func (s *Server[T]) Connect(ctx context.Context, token T, conn Connection, handler io.Writer) error {
	var err error

	s.Do(func() {
		id := token.ID()

		if !s.open {
			err = errors.New("server closed")
			return
		}

		if s.connectionLimit == len(s.connections) {
			err = errors.New("server full")
			return
		}

		if _, found := s.connections[id]; found {
			err = errors.New("id taken")
			return
		}

		if err = s.game.OnConnect(token, conn); err != nil {
			return
		}

		s.connections[id] = conn

		go conn.Open(ctx, handler, func() {
			// disconnect the client
			s.Do(func() {
				if s.open {
					s.game.OnDisconnect(token)
					delete(s.connections, id)
				}
			})
		})
	})

	return err
}

func (s *Server[T]) Open(ctx context.Context, tps int) error {
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

		s.game.OnShutdown()
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
				s.metrics.RecordTick(start, start.Sub(ready), done.Sub(ready))
			}
		}
	}()

	s.game.OnStartup(&engine[T]{s})
	return nil
}

func (s *Server[T]) Do(task func()) {
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

	s.metrics.RecordTask(start, start.Sub(ready), done.Sub(start))
}
