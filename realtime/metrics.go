package realtime

import (
	"time"
)

type Metrics interface {
	RecordLockContention(start time.Time, wait, execution time.Duration)
}

type EmptyMetrics struct {
}

func (m *EmptyMetrics) RecordLockContention(start time.Time, wait, execution time.Duration) {
}
