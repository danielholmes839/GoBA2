package gameplay

import (
	"goba2/games/goba/gameplay/geometry"
	"goba2/games/goba/util"
	"goba2/realtime"
)

// Team struct
type Team struct {
	name        string
	color       string
	size        int
	respawn     *geometry.Point
	events      *realtime.Subscription
	projectiles map[*Projectile]struct{}
	players     map[string]*Champion
}

// NewTeam func
func NewTeam(name string, color string, respawn *geometry.Point) *Team {
	return &Team{
		name:        name,
		color:       color,
		size:        0,
		respawn:     respawn,
		events:      realtime.NewSubscription("team-events"),
		projectiles: make(map[*Projectile]struct{}),
		players:     make(map[string]*Champion),
	}
}

func (team *Team) tick(game *Game) {
	// Tick
	visibleChampions := []*ChampionJSON{}
	visibleProjectiles := []*ProjectileJSON{}

	// Ally players are visible
	for client, champ := range team.players {
		visibleChampions = append(visibleChampions, NewChampionJSON(client, champ))
	}

	// Ally projectiles are visible
	for projectile := range team.projectiles {
		visibleProjectiles = append(visibleProjectiles, NewProjectileJSON(projectile))
	}

	for otherTeam := range game.teams {
		if otherTeam == team {
			continue
		}

		// Adding visible champions on other teams
		for client, otherChampion := range otherTeam.players {
			vision := false
			p2 := otherChampion.hitbox

			// Check players with vision
			for _, teamChampion := range team.players {
				p1 := teamChampion.hitbox
				// Line of sight champion to enemy
				line := geometry.NewLine(p1.GetX(), p1.GetY(), p2.GetX(), p2.GetY())
				if game.hasLineOfSight(line) {
					vision = true
					break
				}
			}

			if vision {
				visibleChampions = append(visibleChampions, NewChampionJSON(client, otherChampion))
				continue
			}

			for teamProjectile := range team.projectiles {
				p1 := teamProjectile.hitbox
				// Line of sight champion to enemy
				line := geometry.NewLine(p1.GetX(), p1.GetY(), p2.GetX(), p2.GetY())
				if game.hasLineOfSight(line) {
					vision = true
				}
			}

			if vision {
				visibleChampions = append(visibleChampions, NewChampionJSON(client, otherChampion))
				continue
			}
		}

		// Adding visible projectiles from other teams
		for otherProjectile := range otherTeam.projectiles {
			p2 := otherProjectile.hitbox
			for _, teamChampion := range team.players {
				p1 := teamChampion.hitbox

				// Line of sight champion to projectile
				line := geometry.NewLine(p1.GetX(), p1.GetY(), p2.GetX(), p2.GetY())
				if game.hasLineOfSight(line) {
					visibleProjectiles = append(visibleProjectiles, NewProjectileJSON(otherProjectile))
					break
				}
			}
		}
	}

	util.BroadcastMessage("tick", NewTickUpdate(visibleChampions, visibleProjectiles), team.events)
}

func (team *Team) addClient(id string, conn realtime.Connection, champion *Champion) {
	// Add a client to the game
	team.events.Subscribe(id, conn)
	team.players[id] = champion
	team.size++
}

func (team *Team) removeClient(id string) {
	// Remove a client to the game
	team.events.Unsubscribe(id)
	delete(team.players, id)
	team.size--
}
