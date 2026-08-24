package objects

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smail1111/mario/internal/utils"
)

type Mario struct {
	PhysicsEntity
	jumps       int
	crouching   bool
	shrinking   bool
	turning     bool
	jumping     bool
	running     bool
	stompStreak int
}

func (m *Mario) Update(movement FVector) (alive bool) {
	m.Collisions.Clear()

	switch m.Death {
	case 0:
		if !m.shrinking {
			rects := m.FRectsAround(m.Game.TileMap)

			for _, goomba := range m.Game.Goombas {
				if goomba.Death == 0 && m.Overlaps(goomba.FRect) {
					m.Shrink()
				}
			}

			for _, koopa := range m.Game.Koopas {
				if koopa.Death == 0 && koopa.Speed != 0.0 && m.Overlaps(koopa.FRect) {
					m.Shrink()
				}
			}

			m.UpdateY(movement.Y, rects)

			for _, block := range m.BlocksAround(m.Game.TileMap) {
				if m.Left() < block.Right() && m.Right() > block.Left() {
					if m.Top() == block.Bottom() {
						block.Hit(m.Entity)
						m.Game.TileMap.Blocks[block.Coord().String()] = block
					}
				}
			}

			for i, goomba := range m.Game.Goombas {
				if m.Death == 0 && goomba.Death == 0 && m.Overlaps(goomba.FRect) && m.Velocity.Y > minimumVelocityY {
					m.Stomp(&goomba.Entity, enemyStompScore, goombaStompAmt)

					m.Game.Goombas[i] = goomba
				}
			}

			for i, koopa := range m.Game.Koopas {
				if m.Death == 0 && koopa.Death == 0 && m.Overlaps(koopa.FRect) {
					switch koopa.Anim.Name {
					case "koopa/walk":
						if m.Velocity.Y > minimumVelocityY {
							m.Stomp(&koopa.Entity, enemyStompScore, koopaShellAmt)
						}

					case "koopa/shell":
						if koopa.Speed == 0.0 {
							if m.Left() < koopa.Left() {
								koopa.Speed = koopaShellMoveSpeed
							} else {
								koopa.Speed = -koopaShellMoveSpeed
							}
						} else if m.Velocity.Y > minimumVelocityY {
							koopa.Speed = 0.0
						}

						if m.Velocity.Y > minimumVelocityY {
							m.Stomp(&koopa.Entity, 0, 0)
						}
					}

					m.Game.Koopas[i] = koopa
				}
			}

			if !ebiten.IsKeyPressed(ebiten.KeyUp) || m.Velocity.Y < minimumJumpingVelocity || m.Collisions.Has(up) {
				m.jumping = false
			} else if m.jumping {
				m.Velocity.Y += marioJumpingAcceleration
			}

			if m.Collisions.Has(down) && !m.jumping {
				if ebiten.IsKeyPressed(ebiten.KeyShift) {
					m.running = true
				} else {
					m.running = false
				}

				m.stompStreak = 0
			}

			m.UpdateX(movement.X, rects)

			if m.Pos.X < m.Game.Offset.X {
				m.Pos.X = m.Game.Offset.X
				m.Velocity.X = 0.0
			}

			m.turning = false

			if m.Anim.Name != "mario/big/crouch" && m.Anim.Name != "mario/big/shrink" {
				if !m.Collisions.Has(down) {
					m.SetAnimation("jump")
				} else if movement.X+m.Velocity.X != 0.0 {
					if m.Velocity.X != 0.0 && (m.Velocity.X < 0.0 != m.Flip[0]) {
						m.SetAnimation("turn")
						m.turning = true
					} else {
						m.SetAnimation("run")
					}
				} else {
					m.SetAnimation("idle")
				}
			}

			if m.Pos.Y > 240 {
				m.Death = marioDeathAmt
			}
		} else if m.Anim.Done {
			m.Type = "mario/small"
			m.SetAnimation("idle")
			m.shrinking = false
			m.Height = 14
			m.Pos.Y += 16
		}

		m.Anim.Update()

	case marioDeathAmt:
		if m.Game.Music.IsPlaying() {
			go m.Game.Music.Close()

			m.Fling(flingVelocityX, flingVelocityY, false, marioDeathAmt)

			m.Game.lives--
			if m.Game.lives > 1 {
				go utils.PlaySound("died")
			} else {
				go utils.PlaySound("game_over")
			}

			m.Type = "mario/small"
			m.SetAnimation("die")
		}

		m.Collisions.Clear()

		m.UpdateX(0.0, nil)
		m.UpdateY(0.0, nil)

		m.Anim.Update()

		return !m.Anim.Done
	}

	return true
}

func (m *Mario) Move(direction int) (movement FVector) {
	if !directions.Has(direction) {
		log.Fatal("Error: Use constant left, right, up, or down.")
	}

	if m.Death == 0 && !m.shrinking && !m.crouching {
		maxVelocity := marioMaxWalkSpeed
		var acceleration float64

		if m.turning {
			acceleration = marioTurnAcceletation
		} else {
			acceleration = marioWalkAcceleration
		}

		if m.running {
			acceleration *= marioRunSpeedMult
			maxVelocity *= marioRunSpeedMult
		}

		switch direction {
		case left:
			movement.X = -0.1
			m.Velocity.X -= acceleration
			m.Velocity.X = max(m.Velocity.X, -maxVelocity)

		case right:
			movement.X = 0.1
			m.Velocity.X += acceleration
			m.Velocity.X = min(m.Velocity.X, maxVelocity)
		}
	}

	return
}

func (m *Mario) Jump() (jumped bool) {
	if m.Collisions.Has(down) {
		m.jumps = marioJumpRefreshAmount
	}

	if m.jumps > 0 && !m.crouching && !m.shrinking {
		m.jumping = true
		m.jumps--

		go utils.PlaySound("jump")

		jumped = true
	}

	return
}

func (m *Mario) Stomp(entity *Entity, score, death int) {
	if score != 0 {
		m.Game.score += score << m.stompStreak
		m.stompStreak = min(m.stompStreak+1, maxEnemyStompStreak)
	}

	m.jumps = marioJumpRefreshAmount - 1
	m.Velocity.Y = marioStompVelocity
	m.Pos.Y = entity.Top() - m.Height
	m.jumping = true

	entity.Death = death

	go utils.PlaySound("stomp")
}

func (m *Mario) Shrink() {
	switch m.Type {
	case "mario/small":
		m.Death = marioDeathAmt

	case "mario/big":
		if m.crouching {
			m.UnCrouch()
		}

		m.SetAnimation("shrink")
		go utils.PlaySound("shrink")

		m.Velocity = FVector{0.0, minimumVelocityY}
		m.shrinking = true
	}
}

func (m *Mario) Crouch() {
	if !m.crouching && !m.shrinking {
		if m.Type == "mario/big" {
			m.SetAnimation("crouch")
			m.crouching = true
			m.Height = 22
			m.Pos.Y += 8
		}
	}
}

func (m *Mario) UnCrouch() {
	if m.crouching && !m.shrinking {
		if m.Type == "mario/big" {
			m.SetAnimation("idle")
			m.crouching = false
			m.Height = 30
			m.Pos.Y -= 8
		}
	}
}

func (m Mario) Draw(screen *ebiten.Image, offset FVector) {
	if m.shrinking {
		if m.Game.frame%marioShrinkHideAmt != 0 {
			m.PhysicsEntity.Draw(screen, offset)
		}
	} else {
		m.PhysicsEntity.Draw(screen, offset)
	}
}

func (g *Game) NewMario(pos FVector) Mario {
	return Mario{PhysicsEntity: g.NewPhysicsEntity("mario/big", pos, 14.0, 30.0, "idle", FVector{-2.0, -2.0})}
}
