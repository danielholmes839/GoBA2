package gameplay

import (
	"encoding/json"
	"goba2/games/goba/gameplay/geometry"
	"goba2/games/goba/util"
	"math"
	"sync"
	"time"
)

// Champion struct
type Champion struct {
	id        string
	maxHealth int
	health    int
	stop      int
	speed     int

	hitbox *geometry.Circle
	target *geometry.Point // moving towards this position

	movementLock  *sync.Mutex
	shootCooldown *Cooldown
	dashCooldown  *Cooldown

	lastHit  string               // The last client to hit this champion
	lastHits map[string]time.Time // The times hit by any other clients
}

// NewChampion func
func NewChampion(id string) *Champion {
	return &Champion{
		id:           id,
		hitbox:       geometry.NewCircle(championStartX, championStartY, championRadius),
		maxHealth:    championMaxHealth,
		health:       championMaxHealth,
		speed:        championSpeed, // units per second
		movementLock: &sync.Mutex{},

		shootCooldown: NewCooldown(shootCooldown),
		dashCooldown:  NewCooldown(dashCooldown),

		lastHit:  "",
		lastHits: make(map[string]time.Time),
	}
}

func (champ *Champion) damage(damage int, enemy string) {
	champ.health -= damage
	champ.lastHit = enemy
	champ.lastHits[enemy] = time.Now()
}

func (champ *Champion) respawn(point *geometry.Point) {
	x, y := point.GetX(), point.GetY()
	champ.health = champ.maxHealth
	champ.hitbox.Move(x, y)
}

func (champ *Champion) death() (string, []string) {
	now := time.Now()
	assists := make([]string, 0)

	for client, timeHit := range champ.lastHits {
		if client != champ.lastHit && now.Before(timeHit.Add(time.Second*10)) {
			assists = append(assists, client)
		}
	}

	return champ.lastHit, assists
}

func (champ *Champion) shoot(event *util.ClientEvent, game *Game) {
	if err := champ.shootCooldown.use(); err != nil {
		return
	}

	data := &ChampionShootEvent{}
	if err := json.Unmarshal(event.GetData(), data); err != nil {
		return
	}

	team := game.getClientTeam(event.Client)
	origin := champ.hitbox.Copy()
	target := geometry.NewPoint(data.X, data.Y)

	projectile := NewProjectile(origin, target, game, event.Client)
	team.projectiles[projectile] = struct{}{}
}

func (champ *Champion) dash() {
	if err := champ.dashCooldown.use(); err != nil {
		return
	}

	champ.speed *= dashSpeedMultiplier

	go func() {
		time.Sleep(dashDuration)
		champ.movementLock.Lock()
		defer champ.movementLock.Unlock()
		champ.speed /= dashSpeedMultiplier
	}()
}

func (champ *Champion) move(game *Game) {
	champ.movementLock.Lock()
	defer champ.movementLock.Unlock()

	// The champion isn't moving
	if champ.target == nil {
		return
	}

	// Calculate the difference between current and target position
	dx := champ.target.GetX() - champ.hitbox.GetX()
	dy := champ.target.GetY() - champ.hitbox.GetY()

	// The target is the current position
	if dx == 0 && dy == 0 {
		return
	}

	// Calculate the speed per tick in each direction
	distance := math.Sqrt(float64(dx*dx + dy*dy))
	speed := float64(champ.speed) / float64(game.tps)               // speed per tick
	speedX := int(math.RoundToEven(float64(dx) / distance * speed)) // speed per tick X axis
	speedY := int(math.RoundToEven(float64(dy) / distance * speed)) // speed per tick Y axis

	// Move the champion
	champ.moveX(game, speedX)
	champ.moveY(game, speedY)

	if distance < speed {
		champ.target = nil
	}
}

func (champ *Champion) moveX(game *Game, speedX int) {
	dirX := direction(speedX)
	position := champ.hitbox

	position.Shift(speedX, 0)
	for _, wall := range game.walls {
		if !wall.HitsCircle(champ.hitbox) {
			continue
		}

		position.Shift(-speedX, 0)
		for i := 0; i < (speedX * dirX); i++ {
			position.Shift(dirX, 0)
			if wall.HitsCircle(champ.hitbox) {
				position.Shift(-dirX, 0)
				break
			}
		}
		break
	}
}

func (champ *Champion) moveY(game *Game, speedY int) {
	dirY := direction(speedY)

	position := champ.hitbox
	position.Shift(0, speedY)
	for _, wall := range game.walls {
		if !wall.HitsCircle(champ.hitbox) {
			continue
		}

		position.Shift(0, -speedY)
		for i := 0; i < (speedY * dirY); i++ {
			position.Shift(0, dirY)
			if wall.HitsCircle(champ.hitbox) {
				position.Shift(0, -dirY)
				break
			}
		}
		break
	}
}

func (champ *Champion) setMovementDirection(event *util.ClientEvent) {
	champ.movementLock.Lock()
	defer champ.movementLock.Unlock()

	movement := &ChampionMoveEvent{}
	if err := json.Unmarshal(event.GetData(), movement); err != nil {
		return
	}

	champ.target = geometry.NewPoint(movement.X, movement.Y)
}
