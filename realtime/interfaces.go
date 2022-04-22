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

type Identity interface {
	ID() string
}

type Application[I Identity] interface {
	HandleOpen(ctx context.Context, engine Engine)
	HandleClose()
	HandleMessage(id string, data []byte)
	HandleConnect(identity I, conn Connection) error
	HandleDisconnect(id string)
}

type Engine interface {
	After(d time.Duration, f func())
	At(t time.Time, f func())
	Interval(t time.Duration, f func())
}
