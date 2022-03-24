package netcode

import "time"

type ServerMetrics interface {
	RecordTask(start time.Time, wait, execution time.Duration)
	RecordTick(start time.Time, wait, execution time.Duration)
}

type ServerHooks[T Token] interface {
	Tick()
	OnConnect(token T, conn Connection) error
	OnDisconnect(token T)
	OnStartup(engine Engine)
	OnShutdown()
}

type Engine interface {
	After(time.Duration, func())
	At(time.Time, func())
	Interval(time.Duration, func())
}

type Token interface {
	ID() string
}
