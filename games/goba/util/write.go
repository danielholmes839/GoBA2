package util

import (
	"goba2/realtime"
	"io"
)

func WriteMessage(event string, data []byte, conn io.Writer) {
	msg := NewServerEvent("personal", event, data).Serialize()
	conn.Write(msg)
}

func BroadcastMessage(event string, data []byte, subscription *realtime.Subscription) {
	msg := NewServerEvent(subscription.Name(), event, data).Serialize()
	subscription.Broadcast(msg)
}
