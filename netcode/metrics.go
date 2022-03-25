package netcode

import (
	"fmt"
	"time"
)

const interval = 5

type LocalServerMetrics struct {
	ticks              int
	ticksExecutionTime time.Duration
	ticksWaitTime      time.Duration
	ticksRecordedAt    time.Time
}

func (m *LocalServerMetrics) RecordTask(start time.Time, wait, execution time.Duration) {
	fmt.Printf("task: %s wait: %s execution: %s\n", start.UTC().Format("2006-01-02T15:04:05Z"), wait, execution)
}

func (m *LocalServerMetrics) RecordTick(start time.Time, wait, execution time.Duration) {
	if m.ticks > 1 && time.Since(m.ticksRecordedAt) > time.Second*5 {
		avg := time.Duration(int64(m.ticksExecutionTime) / int64(m.ticks)).Milliseconds()
		fmt.Printf("tps: %d tick-average-execution: %dms ticks: %d ticks-total-wait: %dns\n", m.ticks/5, avg, m.ticks, m.ticksWaitTime)
		m.ticks = 0
		m.ticksExecutionTime = 0
		m.ticksWaitTime = 0
		m.ticksRecordedAt = time.Now()
	}

	m.ticksExecutionTime += execution
	m.ticks++
	m.ticksWaitTime += wait
}

type EmptyMetrics struct {
}

func (m *EmptyMetrics) RecordTask(start time.Time, wait, execution time.Duration) {
}

func (m *EmptyMetrics) RecordTick(start time.Time, wait, execution time.Duration) {
}
