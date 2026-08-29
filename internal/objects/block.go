package objects

import "github.com/smail1111/not-mario/internal/utils"

type Block struct {
	Entity
	hit bool
}

func (b *Block) Hit(e Entity) {
	switch b.Type {
	case "block/question":
		if b.Anim.Name != "block/question/hit" {
			b.AnimOffset.Y = blockHitAnimOffsetY
			b.hit = true

			b.SetAnimation("hit")

			particle := Particle{
				Entity:   b.Game.NewEntity("coin", FVector{b.Left() + 6.0, b.Top() - tileSize}, 8.0, 14.0, "spin", FVector{0.0, -2.0}),
				Movement: FVector{0.0, coinParticleSpeedY},
				LifeSpan: coinParticleLifespan,
			}

			b.Game.Particles[particle.ID] = particle

			b.Game.coins++
			b.Game.score += coinCollectScore

			go utils.PlaySound("coin")
		}

	case "block/brick":
		b.AnimOffset.Y = blockHitAnimOffsetY
		b.hit = true

		if b.Anim.Name != "block/brick/break" && (e.Anim.Name == "mario/big/jump" || e.Anim.Name == "koopa/shell") {
			b.SetAnimation("break")
			go utils.PlaySound("brick_break")
		}
	}
}

func (b *Block) Update() (alive bool) {
	b.Anim.Update()

	b.AnimOffset.Y = min(0.0, b.AnimOffset.Y+blockNormalizationY)

	b.hit = false

	switch b.Type {
	case "block/brick":
		if b.Anim.Name == "block/brick/break" && b.Anim.Done {
			b.Death++
		}
	}

	return b.Death == 0
}

func (b Block) Coord() Vector {
	return Vector{utils.FloorDiv(b.Pos.X, tileSize), utils.FloorDiv(b.Pos.Y, tileSize)}
}
