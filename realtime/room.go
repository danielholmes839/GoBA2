package realtime

import "errors"

var ErrRoomFull = errors.New("server full")
var ErrAlreadyConnected = errors.New("id already connected")
var ErrNotConnected = errors.New("id not connected")

type BasicRoom struct {
	connectionLimit int
	connections     map[string]struct{}
}

func NewBasicRoom(connectionLimit int) *BasicRoom {
	return &BasicRoom{
		connectionLimit: connectionLimit,
		connections:     map[string]struct{}{},
	}
}

func (r *BasicRoom) Connect(id string) error {
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

func (r *BasicRoom) Disconnect(id string) {
	delete(r.connections, id)
}
