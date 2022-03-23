package engine

import (
	"context"
	"fmt"
	"time"
)

type TickFunc func()

type Base struct {
	tasks chan func()
	tps   int
	tick  TickFunc
}

func NewBase() *Base {
	return &Base{
		tasks: make(chan func()),
		tick:  func() {},
	}
}

func (b *Base) Run(ctx context.Context, tps int) {
	// create the ticker
	d := time.Duration(int64(float64(time.Second) / float64(tps)))
	ticker := time.NewTicker(d)

	// stop the ticker
	defer func() {
		ticker.Stop()
		fmt.Println("game closed")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.tick()
		case task := <-b.tasks:
			task()
		}
	}
}

func (b *Base) SetTickHandler(tick TickFunc) {
	b.tick = tick
}
