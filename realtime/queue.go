package realtime

import (
	"container/list"
	"sync"
)

type Queue[T any] struct {
	sync.Mutex
	queue *list.List
}

func (q *Queue[T]) ReadAll() []T {
	q.Lock()
	defer q.Unlock()

	events := make([]T, q.queue.Len())
	for i := 0; q.queue.Len() > 0; i++ {
		event := q.queue.Front()
		events[i] = q.queue.Remove(event).(T)
	}

	return events
}

// Push an event to the queue.
func (q *Queue[T]) Push(data T) {
	q.Lock()
	defer q.Unlock()
	q.queue.PushBack(data)
}

// NewClientEventQueue func
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		queue: list.New(),
	}
}
