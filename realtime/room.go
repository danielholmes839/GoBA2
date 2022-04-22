package realtime

import "errors"

var ErrRoomFull = errors.New("server full")
var ErrAlreadyConnected = errors.New("id already connected")
var ErrNotConnected = errors.New("id not connected")

type Room struct {
	connectionLimit int
	connections     map[string]Connection
}

func NewRoom(connectionLimit int) *Room {
	return &Room{
		connectionLimit: connectionLimit,
		connections:     map[string]Connection{},
	}
}

func (r *Room) Connect(id string, conn Connection) error {
	if len(r.connections) == r.connectionLimit {
		return ErrRoomFull
	}

	if _, exists := r.connections[id]; exists {
		return ErrAlreadyConnected
	}

	// add the connection
	r.connections[id] = conn
	return nil
}

func (r *Room) Disconnect(id string) error {
	if _, exists := r.connections[id]; !exists {
		return ErrNotConnected
	}

	// delete the connection
	r.connections[id].Close()
	delete(r.connections, id)
	return nil
}
