package realtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitRoomConnect(t *testing.T) {
	room := NewLimitRoom(2)

	err := room.Connect("1")
	assert.Nil(t, err)

	err = room.Connect("1")
	assert.ErrorIs(t, err, ErrAlreadyConnected)

	err = room.Connect("2")
	assert.Nil(t, err)

	err = room.Connect("3")
	assert.ErrorIs(t, err, ErrRoomFull)
}

func TestLimitRoomDisconnect(t *testing.T) {
	room := NewLimitRoom(2)

	err := room.Connect("1")
	assert.Nil(t, err)

	room.Disconnect("1")
	room.Disconnect("1")
}
