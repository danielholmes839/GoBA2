package netcode

import (
	"context"
	"time"
)

type engine[T Token] struct {
	*Server[T]
}

func (e *engine[T]) After(task string, d time.Duration, f func()) {
	go func() {
		ctx, cancel := context.WithTimeout(e.ctx, d)
		defer cancel()

		<-ctx.Done()
		if ctx.Err() == context.Canceled {
			return
		}

		e.Do(task, f)
	}()
}

func (e *engine[T]) At(task string, t time.Time, f func()) {
	if now := time.Now(); now.Before(t) {
		e.After(task, t.Sub(now), f)
	}
}

func (e *engine[T]) Interval(task string, d time.Duration, f func()) {
	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for {
			select {
			case <-e.ctx.Done():
				return
			case <-ticker.C:
				e.Do(task, f)
			}
		}
	}()
}
