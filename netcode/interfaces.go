package netcode

import (
	"context"
	"time"
)

type ServerMetrics interface {
	RecordTask(start time.Time, wait, execution time.Duration)
	RecordTick(start time.Time, wait, execution time.Duration)
}

type Identity interface {
	ID() string
}

/* ServerHooks interface
implemented by games
*/
type ServerHooks[I Identity] interface {
	Tick()
	OnMessage(id string, data []byte)
	OnConnect(identity I, conn Connection) error
	OnDisconnect(identity I)
	OnOpen(ctx context.Context, engine Engine)
	OnClose()
}

type Engine interface {
	After(d time.Duration, f func())
	At(t time.Time, f func())
	Interval(t time.Duration, f func())
}
