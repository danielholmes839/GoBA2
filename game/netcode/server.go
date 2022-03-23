package netcode

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ServerMetrics interface {
	RecordTask(start time.Time, wait, execution time.Duration)
	RecordTick(start time.Time, execution time.Duration)
}

// TODO
// type ServerHooks interface {
// 	Opened()
// 	Closed()
// }

type Server struct {
	ServerMetrics
	Name     string
	Tasks    chan func()
	TickFunc func()
}

func (s *Server) Run(ctx context.Context, tps int) {
	// calculate the
	target := time.Duration(int64(float64(time.Second) / float64(tps)))
	ticker := time.NewTicker(target)

	defer func() {
		ticker.Stop()
		fmt.Printf("server: %s, stopped\n", s.Name)
	}()

	fmt.Printf("server: %s, started\n", s.Name)
	for {
		select {
		case <-ctx.Done():
			// context cancelled
			return
		case <-ticker.C:
			// execute ticks
			now := time.Now()
			s.TickFunc()
			s.RecordTick(now, time.Since(now))
		case task := <-s.Tasks:
			// execute task
			task()
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
