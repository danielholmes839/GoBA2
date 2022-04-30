package realtime

import (
	"context"
	"errors"
)

type mockid struct {
	id string
}

func (m mockid) ID() string {
	return m.id
}

type mockgame struct {
	OnConnectErr error
}

func (m *mockgame) HandleMessage(id string, data []byte) {

}

func (m *mockgame) HandleConnect(identity mockid, conn Connection) error {
	return m.OnConnectErr
}

func (m *mockgame) HandleDisconnect(id string) {

}

func (m *mockgame) HandleOpen(ctx context.Context, engine Engine) {

}

func (m *mockgame) HandleClose() {

}

type mockconn struct {
	conn chan []byte
}

func (m *mockconn) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockconn) Close() error {
	defer recover() // closing a closed channel
	close(m.conn)
	return nil
}

func (m *mockconn) Receive() ([]byte, error) {
	data, ok := <-m.conn
	if !ok {
		return data, errors.New("mockconn closed")
	}
	return data, nil
}
