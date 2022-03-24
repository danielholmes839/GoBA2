package netcode

import (
	"context"
	"time"
)

type engine[T Token] struct {
	*Server[T]
}

func (e *engine[T]) After(d time.Duration, task func()) {
	go func() {
		ctx, cancel := context.WithTimeout(e.ctx, d)
		defer cancel()

		<-ctx.Done()
		if ctx.Err() == context.Canceled {
			return
		}

		e.Do(task)
	}()
}

func (e *engine[T]) Interval(d time.Duration, task func()) {
	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for {
			select {
			case <-e.ctx.Done():
				return
			case <-ticker.C:
				e.Do(task)
			}
		}
	}()
}

func (e *engine[T]) At(t time.Time, task func()) {
	if now := time.Now(); now.Before(t) {
		e.After(t.Sub(now), task)
	}
}
