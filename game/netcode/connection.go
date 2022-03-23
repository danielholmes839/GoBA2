package netcode

import (
	"context"
	"io"
)

type Callback func()

type Connection interface {
	io.Writer
	Open(ctx context.Context, handler io.Writer, close Callback)
}
