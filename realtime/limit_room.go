package realtime

import "errors"

var ErrRoomFull = errors.New("server full")
var ErrAlreadyConnected = errors.New("id already connected")
var ErrNotConnected = errors.New("id not connected")

type LimitRoom struct {
	connectionLimit int
	connections     map[string]struct{}
}

func NewLimitRoom(connectionLimit int) *LimitRoom {
	return &LimitRoom{
		connectionLimit: connectionLimit,
		connections:     map[string]struct{}{},
	}
}

func (r *LimitRoom) Connect(id string) error {
	if len(r.connections) == r.connectionLimit {
		return ErrRoomFull
	}

	if _, exists := r.connections[id]; exists {
		return ErrAlreadyConnected
	}

	// add the connection
	r.connections[id] = struct{}{}
	return nil
}

func (r *LimitRoom) Disconnect(id string) {
	delete(r.connections, id)
}
