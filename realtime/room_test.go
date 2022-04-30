package realtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomConnect(t *testing.T) {
	room := NewRoom(2)

	err := room.Connect("1")
	assert.Nil(t, err)

	err = room.Connect("1")
	assert.ErrorIs(t, err, ErrAlreadyConnected)

	err = room.Connect("2")
	assert.Nil(t, err)

	err = room.Connect("3")
	assert.ErrorIs(t, err, ErrRoomFull)
}

func TestRoomDisconnect(t *testing.T) {
	room := NewRoom(2)

	err := room.Connect("1")
	assert.Nil(t, err)

	err = room.Disconnect("1")
	assert.Nil(t, err)

	err = room.Disconnect("1")
	assert.ErrorIs(t, err, ErrNotConnected)
}
