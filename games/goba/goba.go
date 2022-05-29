package goba

import "goba2/games/goba/gameplay"

type Goba struct {
	*gameplay.Game
}

func New(tps int) *Goba {
	return &Goba{
		Game: gameplay.NewGame(tps),
	}
}