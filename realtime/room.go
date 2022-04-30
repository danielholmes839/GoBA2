package realtime

import "errors"

var ErrRoomFull = errors.New("server full")
var ErrAlreadyConnected = errors.New("id already connected")
var ErrNotConnected = errors.New("id not connected")

type Room struct {
	connectionLimit int
	connections     map[string]struct{}
}

func NewRoom(connectionLimit int) *Room {
	return &Room{
		connectionLimit: connectionLimit,
		connections:     map[string]struct{}{},
	}
}

func (r *Room) Connect(id string) error {
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

func (r *Room) Disconnect(id string) error {
	if _, exists := r.connections[id]; !exists {
		return ErrNotConnected
	}

	// delete the connection
	delete(r.connections, id)
	return nil
}
