package netcode

import (
	"context"
	"time"
)

type ServerMetrics interface {
	RecordTask(start time.Time, wait, execution time.Duration)
	RecordTick(start time.Time, wait, execution time.Duration)
}

type ServerHooks[T Token] interface {
	Tick()
	OnMessage(user string, data []byte)
	OnConnect(token T, conn Connection) error
	OnDisconnect(token T)
	OnOpen(ctx context.Context, engine Engine)
	OnClose()
}

type Engine interface {
	After(d time.Duration, f func())
	At(t time.Time, f func())
	Interval(t time.Duration, f func())
}

type Token interface {
	ID() string
}
