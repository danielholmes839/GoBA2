package gameplay

import (
	"goba2/games/goba/gameplay/geometry"
)

// TeamJSON struct
type TeamJSON struct {
	Color string `json:"color"`
	Size  int    `json:"size"`
}

// RectangleJSON struct
type RectangleJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// NewRectangleJSON func
func NewRectangleJSON(rect *geometry.Rectangle) *RectangleJSON {
	return &RectangleJSON{
		X: rect.GetX(),
		Y: rect.GetY(),
		W: rect.GetWidth(),
		H: rect.GetHeight(),
	}
}

// ChampionJSON struct
type ChampionJSON struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Health int    `json:"health"`
	Radius int    `json:"r"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

// NewChampionJSON func
func NewChampionJSON(client string, champ *Champion) *ChampionJSON {
	return &ChampionJSON{
		ID:     champ.id,
		Name:   client,
		Health: champ.health,
		Radius: champ.hitbox.GetRadius(),
		X:      champ.hitbox.GetX(),
		Y:      champ.hitbox.GetY(),
	}
}

// ProjectileJSON struct
type ProjectileJSON struct {
	Team   string `json:"team"`
	Radius int    `json:"r"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

// NewProjectileJSON func
func NewProjectileJSON(projectile *Projectile) *ProjectileJSON {
	return &ProjectileJSON{
		Team:   projectile.team.name,
		Radius: projectile.hitbox.GetRadius(),
		X:      projectile.hitbox.GetX(),
		Y:      projectile.hitbox.GetY(),
	}
}

// ChampionMoveEvent struct
type ChampionMoveEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ChampionShootEvent struct
type ChampionShootEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}
