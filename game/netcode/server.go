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
	RecordTick(start time.Time, execution time.Duration)
}

type ServerHooks interface {
	Tick()
	OnConnect(id string)
	OnConnectError(id string, err error)
	OnDisconnect(id string)
}

type Engine interface {
	Do(task func())
}

type Identity interface {
	ID() string
}

type Server struct {
	ServerHooks
	ServerMetrics
	Name            string
	Tasks           chan func()
	Connections     map[string]Connection
	MAX_CONNECTIONS int
}

func (s *Server) Connect(ctx context.Context, conn Connection, handler io.Writer, id string) error {
	var err error

	s.Do(func() {
		if len(s.Connections) >= s.MAX_CONNECTIONS {
			// maximum connections reached
			err = errors.New("already connected")
			s.OnConnectError(id, err)
			return
		}

		if _, found := s.Connections[id]; found {
			// already connected
			err = errors.New("maximum connections reached")
			s.OnConnectError(id, err)
			return
		}

		go conn.Open(ctx, handler, func() {
			// disconnect the client
			s.Do(func() {
				delete(s.Connections, id)
				s.OnDisconnect(id)
			})
		})

		s.Connections[id] = conn
		s.OnConnect(id)
	})

	return err
}

func (s *Server) Open(ctx context.Context, tps int) {
	// calculate the delay to achieve the correct tps
	target := time.Duration(int64(float64(time.Second) / float64(tps)))
	ticker := time.NewTicker(target)
	open := true

	go func() {
		<-ctx.Done()
		time.Sleep(time.Second * 5)

		// stop the server
		open = false
		ticker.Stop()
		close(s.Tasks)
	}()

	for open {
		select {
		case <-ticker.C:
			// execute ticks
			now := time.Now()
			s.Tick()
			s.RecordTick(now, time.Since(now))
			break

		case task, ok := <-s.Tasks:
			// execute task
			if ok {
				task()
			}
		}
	}
}

func (s *Server) Do(task func()) {
	// wait group to block until the task has completed
	wg := sync.WaitGroup{}
	wg.Add(1)

	// measure the wait and execution time of tasks
	var ready, started time.Time
	var wait, execution time.Duration
	ready = time.Now()

	// execute the task
	s.Tasks <- func() {
		started = time.Now()
		task()
		wg.Done()
	}

	// block
	wg.Wait()

	// record metrics
	execution = time.Now().Sub(started)
	wait = started.Sub(ready)
	s.RecordTask(ready, wait, execution)
}
