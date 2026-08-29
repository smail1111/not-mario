package objects

import (
	"math"

	"github.com/smail1111/not-mario/internal/utils"
)

type Koopa struct {
	PhysicsEntity
	Speed float64
}

func (k *Koopa) Update(movement FVector) (alive bool) {
	switch k.Death {
	case 0:
		movement.X += k.Speed

		k.PhysicsEntity.Update(movement)

		if movement.X != 0.0 {
			k.Flip[0] = !k.Flip[0]
		}

		if k.Collisions.Has(right) || k.Collisions.Has(left) || (k.Collisions.Has(down) &&
			(k.Anim.Name == "koopa/walk" && !k.Game.TileMap.CheckSolidTile(
				k.Pos.X+k.Width/2.0+math.Abs(k.Speed)/k.Speed*(tileSize/2.0),
				k.Pos.Y+k.Height/2.0+tileSize,
			))) {

			k.Speed *= -1.0
		}

		switch k.Anim.Name {
		case "koopa/walk":
			for _, goomba := range k.Game.Goombas {
				if goomba.Death == 0 && k.Overlaps(goomba.FRect) {
					if k.Left() < goomba.Left() {
						k.Speed = -math.Abs(k.Speed)
					} else {
						k.Speed = math.Abs(k.Speed)
					}
				}
			}

			for _, koopa := range k.Game.Koopas {
				if k.ID != koopa.ID && koopa.Death == 0 && k.Overlaps(koopa.FRect) {
					if k.Left() < koopa.Left() {
						k.Speed = -math.Abs(k.Speed)
					} else {
						k.Speed = math.Abs(k.Speed)
					}
				}
			}

		case "koopa/shell":
			if k.Speed != 0.0 {
				k.Anim.Frame = 0

				for key, goomba := range k.Game.Goombas {
					if goomba.Death == 0 && k.Overlaps(goomba.FRect) {
						goomba.Fling(flingVelocityX, flingVelocityY, true, goombaDeathAmount)
						go utils.PlaySound("kick")
						k.Game.Goombas[key] = goomba
					}
				}

				for key, koopa := range k.Game.Koopas {
					if k.ID != koopa.ID && koopa.Death == 0 && k.Overlaps(koopa.FRect) {
						koopa.Fling(flingVelocityX, flingVelocityY, true, koopaDeathAmt)
						go utils.PlaySound("kick")
						k.Game.Koopas[key] = koopa
					}
				}
			}

			if k.Anim.Done {
				k.SetAnimation("walk")
				k.AnimOffset = FVector{-2.0, -10.0}
				k.Speed = koopaWalkSpeed
			}
		}

		for _, block := range k.BlocksAround(k.Game.TileMap) {
			if k.Left() < block.Right() && k.Right() > block.Left() {
				if k.Bottom() == block.Top() && block.hit {
					go utils.PlaySound("kick")
					k.Fling(flingVelocityX, flingVelocityY, false, koopaShellAmt)
					k.Anim.Frame = 0
				}
			} else if k.Anim.Name == "koopa/shell" {
				if k.Top() < block.Bottom() && k.Bottom() > block.Top() {
					if k.Left() == block.Right() || k.Right() == block.Left() {
						block.Hit(k.Entity)
						k.Game.TileMap.Blocks[block.Coord().String()] = block
					}
				}
			}
		}

	case koopaShellAmt:
		k.SetAnimation("shell")
		k.AnimOffset = FVector{0.0, 0.0}
		k.Speed = 0.0
		k.Death = 0

	case koopaDeathAmt:
		k.Collisions.Clear()

		k.UpdateX(0.0, nil)
		k.UpdateY(0.0, nil)
	}

	return k.Pos.Y < 320
}

func (g *Game) NewKoopa(pos FVector) Koopa {
	koopa := Koopa{
		PhysicsEntity: g.NewPhysicsEntity("koopa", pos, 14.0, 14.0, "walk", FVector{-2.0, -10.0}),
		Speed:         -koopaWalkSpeed,
	}

	g.Koopas[koopa.ID] = koopa

	return koopa
}
