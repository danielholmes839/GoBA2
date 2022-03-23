package netcode

import (
	"fmt"
	"time"
)

type LocalServerMetrics struct {
	ticks              int
	ticksExecutionTime time.Duration
	ticksRecordedAt    time.Time
}

func (m *LocalServerMetrics) RecordTask(start time.Time, wait, execution time.Duration) {
	fmt.Printf("task: %s wait: %s execution: %s\n", start.UTC().Format("2006-01-02T15:04:05Z"), wait, execution)
}

func (m *LocalServerMetrics) RecordTick(start time.Time, execution time.Duration) {
	m.ticksExecutionTime += execution
	m.ticks++

	if time.Since(m.ticksRecordedAt) > time.Second*2 {
		avg := time.Duration(int64(m.ticksExecutionTime) / int64(m.ticks)).Milliseconds()
		fmt.Printf("tps: %d tick-duration: %dms\n", m.ticks/2, avg)
		m.ticks = 0
		m.ticksExecutionTime = 0
		m.ticksRecordedAt = time.Now()
	}
}
