package gameplay

import (
	"goba2/games/goba/gameplay/geometry"
	"math"
)

// Projectile struct
type Projectile struct {
	hit    bool
	speedX int
	speedY int
	origin *geometry.Point
	hitbox *geometry.Circle
	team   *Team
	client string
}

// NewProjectile function
func NewProjectile(origin *geometry.Point, target *geometry.Point, game *Game, client string) *Projectile {
	x, y := origin.GetX(), origin.GetY()

	dx := (target.GetX() - origin.GetX())
	dy := (target.GetY() - origin.GetY())

	speedPerSecond := projectileSpeed
	speedPerTick := float64(speedPerSecond) / float64(game.tps)
	distance := math.Sqrt(float64((dx * dx) + (dy * dy)))

	speedX := int(math.RoundToEven((float64(dx) / distance) * speedPerTick)) // speed per tick
	speedY := int(math.RoundToEven((float64(dy) / distance) * speedPerTick)) // speed per tick

	return &Projectile{
		speedX: speedX,
		speedY: speedY,
		origin: geometry.NewPoint(x, y),
		hitbox: geometry.NewCircle(x, y, projectileRadius),
		team:   game.getClientTeam(client),
		client: client,
	}
}

func (projectile *Projectile) move() {
	projectile.hitbox.Shift(projectile.speedX, projectile.speedY)
}

// Check for collisions with other players
func (projectile *Projectile) collisions(game *Game) {
	for _, info := range game.clients {
		champ := info.champion
		team := info.team

		if team == projectile.team {
			continue
		}

		// The projectiles hit a champion
		if projectile.hitbox.HitsCircle(champ.hitbox) {
			projectile.hit = true
			champ.damage(projectileDamage, projectile.client)
		}
	}
}

// The projectile should be deleted
func (projectile *Projectile) done() bool {
	return projectile.hit || (projectile.hitbox.Distance2(projectile.origin) > (projectileRange * projectileRange))
}
