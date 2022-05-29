package realtime

import "io"

type Subscription struct {
	name        string
	subscribers map[string]io.Writer
}

func NewSubscription(name string) *Subscription {
	return &Subscription{
		name:        name,
		subscribers: make(map[string]io.Writer),
	}
}

func (s *Subscription) Name() string {
	return s.name
}

func (s *Subscription) Subscribe(id string, conn io.Writer) {
	s.subscribers[id] = conn
}

func (s *Subscription) Unsubscribe(id string) {
	delete(s.subscribers, id)
}

func (s *Subscription) Broadcast(data []byte) {
	for _, conn := range s.subscribers {
		conn.Write(data)
	}
}
