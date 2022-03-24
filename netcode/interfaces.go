package netcode

import "time"

type ServerMetrics interface {
	RecordTask(task string, start time.Time, wait, execution time.Duration)
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
	After(task string, d time.Duration, f func())
	At(task string, t time.Time, f func())
	Interval(task string, t time.Duration, f func())
}

type Token interface {
	ID() string
}
