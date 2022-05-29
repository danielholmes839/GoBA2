package goba

import (
	"context"
	"goba2/games/goba/gameplay"
)

type Goba struct {
	*gameplay.Game
}

func NewGoba(tps int, shutdown context.CancelFunc) *Goba {
	return &Goba{
		Game: gameplay.NewGame(tps, shutdown),
	}
}
