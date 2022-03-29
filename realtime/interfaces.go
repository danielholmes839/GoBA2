package realtime

import (
	"context"
	"io"
	"time"
)

type Connection interface {
	io.Writer
	io.Closer
	Receive() ([]byte, error)
}

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
	OnConnect(identity I, connection io.Writer) error
	OnDisconnect(identity I)
	OnOpen(scheduler Scheduler)
	OnClose()
}

type Scheduler interface {
	After(d time.Duration, f func())
	At(t time.Time, f func())
	Interval(t time.Duration, f func())
	Context() context.Context
	Cancel()
}
