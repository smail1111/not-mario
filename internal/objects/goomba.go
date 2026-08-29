package objects

import (
	"math"

	"github.com/smail1111/not-mario/internal/utils"
)

type Goomba struct {
	PhysicsEntity
	Speed float64
}

// Updates a Goomba.
func (g *Goomba) Update(movement FVector) bool {
	switch g.Death {
	case 0:
		movement.X += g.Speed

		g.PhysicsEntity.Update(movement)

		for _, goomba := range g.Game.Goombas {
			if g.ID != goomba.ID && goomba.Death == 0 && g.Overlaps(goomba.FRect) {
				if g.Left() < goomba.Left() {
					g.Speed = -math.Abs(g.Speed)
				} else {
					g.Speed = math.Abs(g.Speed)
				}
			}
		}

		for _, koopa := range g.Game.Koopas {
			if koopa.Death == 0 && g.Overlaps(koopa.FRect) {
				if koopa.Death == 0 && g.Overlaps(koopa.FRect) {
					if g.Left() < koopa.Left() {
						g.Speed = -math.Abs(g.Speed)
					} else {
						g.Speed = math.Abs(g.Speed)
					}
				}
			}
		}

		for _, block := range g.BlocksAround(g.Game.TileMap) {
			if g.Left() < block.Right() && g.Right() > block.Left() && g.Bottom() == block.Top() && block.hit {
				g.Fling(flingVelocityX, flingVelocityY, true, goombaDeathAmount)
				go utils.PlaySound("kick")
			}
		}

		if g.Collisions.Has(down) &&
			(g.Collisions.Has(right) || g.Collisions.Has(left) ||
				!g.Game.TileMap.CheckSolidTile(
					g.Pos.X+g.Width/2+math.Abs(g.Speed)/g.Speed*(tileSize/2),
					g.Pos.Y+g.Height/2+tileSize)) {

			g.Speed *= -1.0
		}

	case goombaStompAmt:
		g.SetAnimation("stomp")
		g.Anim.Update()

	case goombaDeathAmount:
		g.Collisions.Clear()

		g.UpdateX(0.0, nil)
		g.UpdateY(0.0, nil)
	}

	return g.Death == 0 || !(g.Anim.Done || g.Pos.Y > 320)
}

// Returns a new Goomba with the given position.
func (g *Game) NewGoomba(pos FVector) Goomba {
	goomba := Goomba{
		PhysicsEntity: g.NewPhysicsEntity("goomba", pos, 14.0, 14.0, "idle", FVector{-2.0, -2.0}),
		Speed:         -goombaSpeed,
	}

	g.Goombas[goomba.ID] = goomba

	return goomba
}
