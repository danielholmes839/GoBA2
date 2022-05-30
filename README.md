# GoBA2

Explanation of GoBA: https://github.com/danielholmes839/GoBA

The goal of GoBA2 was to build a generic framework for managing realtime games.
The way I managed the lifecylce of games and connections in the original GoBA was extremely clumsy.
This repo contains a refactored version of GoBA that uses the new [realtime framework](./realtime) I created.
I may create a polished version of the framework in the future.

## `realtime` Package

The realtime package can be used to run anything that implements the `App` interface.

```go
type Identity interface {
	ID() string
}

type App[I Identity] interface {
	HandleOpen(ctx context.Context, engine Engine)
	HandleClose()
	HandleMessage(id string, data []byte)
	HandleConnect(identity I, conn Connection) error
	HandleDisconnect(id string)
}
```
