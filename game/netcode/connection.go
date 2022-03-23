package netcode

import (
	"context"
	"io"
)

type Callback func()

type Connection interface {
	io.Writer // write data to the client
	io.Closer
	// Open method - start reading from the client on loop
	Open(ctx context.Context, handler io.Writer, close Callback) error
}
