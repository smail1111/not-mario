package objects

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smail1111/mario/internal/utils"
)

type Entity struct {
	FRect
	Anim       utils.Animation
	AnimOffset FVector
	Game       *Game
	Type       string
	Flip       [2]bool
	ID         int
	Death      int
}

func (pe Entity) Draw(screen *ebiten.Image, offset FVector) {
	opts := ebiten.DrawImageOptions{}

	if pe.Flip[0] {
		opts.GeoM.Translate(-pe.Width/2+pe.AnimOffset.X, 0.0)
		opts.GeoM.Scale(-1.0, 1.0)
		opts.GeoM.Translate(pe.Width/2-pe.AnimOffset.X, 0.0)
	}

	if pe.Flip[1] {
		opts.GeoM.Translate(0.0, -pe.Height/2+pe.AnimOffset.Y)
		opts.GeoM.Scale(1.0, -1.0)
		opts.GeoM.Translate(0.0, pe.Height/2-pe.AnimOffset.Y)
	}

	opts.GeoM.Translate(pe.Pos.X-offset.X+pe.AnimOffset.X, pe.Pos.Y-offset.Y+pe.AnimOffset.Y)
	screen.DrawImage(pe.Anim.Image(), &opts)
}

// Sets the animation of an Entity based on its type and given action
// if the animation is not the Entity's current animation.
func (e *Entity) SetAnimation(action string) {
	animName := e.Type + "/" + action

	if animName != e.Anim.Name {
		e.Anim = e.Game.Anims[animName]
	}
}

// Returns a new Entity based on the given type, position, width, height, animation name, and animation offset.
func (g *Game) NewEntity(entityType string, pos FVector, width float64, height float64, animName string, animOffset FVector) Entity {
	entity := Entity{
		Type: entityType,
		FRect: FRect{
			Pos:    pos,
			Width:  width,
			Height: height,
		},
		AnimOffset: animOffset,
		Game:       g,
		ID:         g.IDCounter(),
	}

	entity.SetAnimation(animName)

	return entity
}

type PhysicsEntity struct {
	Entity
	Collisions   utils.Set
	LastMovement FVector
	Velocity     FVector
}

// Updates a Physics Entity based on the given movement.
func (pe *PhysicsEntity) Update(movement FVector) {
	pe.Collisions.Clear()

	rects := pe.FRectsAround(pe.Game.TileMap)

	pe.UpdateY(movement.Y, rects)

	pe.UpdateX(movement.X, rects)

	pe.Anim.Update()
}

// Updates a Physics Entity on the X dimension based on the given movement and rects.
func (pe *PhysicsEntity) UpdateX(movement float64, rects []FRect) {
	pe.LastMovement.X = movement

	frameMovement := pe.Velocity.X + movement

	pe.Pos.X += frameMovement
	for _, rect := range rects {
		if pe.Overlaps(rect) {
			if frameMovement > 0.0 {
				pe.Collisions.Add(right)
				pe.Pos.X = rect.Left() - pe.Width
			} else if frameMovement < 0.0 {
				pe.Collisions.Add(left)
				pe.Pos.X = rect.Right()
			}
		}
	}

	if pe.Collisions.Has(right) || pe.Collisions.Has(left) {
		pe.Velocity.X = 0.0
	} else if pe.Velocity.X != 0.0 {
		if pe.Velocity.X < 0.0 {
			pe.Velocity.X = min(pe.Velocity.X+normalizationX, 0.0)
		} else {
			pe.Velocity.X = max(pe.Velocity.X-normalizationX, 0.0)
		}
	}

	if movement != 0.0 {
		pe.Flip[0] = movement < 0.0
	}
}

// Updates a Physics Entity on the Y dimension based on the given movement and rects.
func (pe *PhysicsEntity) UpdateY(movement float64, rects []FRect) {
	pe.LastMovement.Y = movement

	frameMovement := pe.Velocity.Y + movement

	pe.Pos.Y += frameMovement
	for _, rect := range rects {
		if pe.Overlaps(rect) {
			if frameMovement > 0.0 {
				pe.Collisions.Add(down)
				pe.Pos.Y = rect.Top() - pe.Height
			} else if frameMovement < 0.0 {
				pe.Collisions.Add(up)
				pe.Pos.Y = rect.Bottom()
			}
		}
	}

	if pe.Collisions.Has(down) || pe.Collisions.Has(up) {
		pe.Velocity.Y = minimumVelocityY
	} else {
		pe.Velocity.Y = min(pe.Velocity.Y+gravity, terminalVelocityY)
	}
}

func (pe *PhysicsEntity) Fling(velocityX, velocityY float64, flip bool, death int) {
	pe.Flip[1] = flip
	pe.Velocity.X = velocityX
	pe.Velocity.Y = velocityY
	pe.Death = death
}

// Returns a slice of Tiles near the Physics Entity based on the given TileMap.
func (pe PhysicsEntity) TilesAround(tm TileMap) (tilesAround []Tile) {
	coord := Vector{utils.FloorDiv(pe.Pos.X+pe.Width/2, tm.TileSize), utils.FloorDiv(pe.Pos.Y+pe.Height/2, tm.TileSize)}

	checked := utils.NewSet()
	for sizeOffsetX := range utils.FloorDiv(pe.Width, tm.TileSize) + 2 {
		for sizeOffsetY := range utils.FloorDiv(pe.Height, tm.TileSize) + 2 {
			for _, flip := range []Vector{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}} {
				offset := Vector{coord.X + sizeOffsetX*flip.X, coord.Y + sizeOffsetY*flip.Y}
				if !checked.Has(offset) {
					if tile, ok := tm.Tiles[offset.String()]; ok {
						tilesAround = append(tilesAround, tile)
					}
					checked.Add(offset)
				}
			}
		}
	}

	return
}

// Returns a slice of Blocks near the Physics Entity based on the given TileMap.
func (pe PhysicsEntity) BlocksAround(tm TileMap) (BlocksAround []Block) {
	coord := Vector{utils.FloorDiv(pe.Pos.X+pe.Width/2, tm.TileSize), utils.FloorDiv(pe.Pos.Y+pe.Height/2, tm.TileSize)}

	checked := utils.NewSet()
	for sizeOffsetX := range utils.FloorDiv(pe.Width, tm.TileSize) + 2 {
		for sizeOffsetY := range utils.FloorDiv(pe.Height, tm.TileSize) + 2 {
			for _, flip := range []Vector{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}} {
				offset := Vector{coord.X + sizeOffsetX*flip.X, coord.Y + sizeOffsetY*flip.Y}
				if !checked.Has(offset) {
					if block, ok := tm.Blocks[offset.String()]; ok {
						BlocksAround = append(BlocksAround, block)
					}
					checked.Add(offset)
				}
			}
		}
	}

	return
}

// Returns a slice of FRects near the Physics Entity based on the given TileMap.
func (pe PhysicsEntity) FRectsAround(tm TileMap) (FRectsAround []FRect) {
	tiles := pe.TilesAround(tm)

	blocks := pe.BlocksAround(tm)

	for _, tile := range tiles {
		if collisionTypes.Has(tile.Type) {
			FRectsAround = append(FRectsAround, tile.FRect(tm.TileSize))
		}
	}

	for _, block := range blocks {
		if collisionTypes.Has(block.Type) {
			FRectsAround = append(FRectsAround, block.FRect)
		}
	}

	return
}

// Returns a new Physics Entity based on the given type, position, width, height, animation name, and animation offset.
func (g *Game) NewPhysicsEntity(entityType string, pos FVector, width float64, height float64, animName string, animOffset FVector) PhysicsEntity {
	return PhysicsEntity{
		Entity: g.NewEntity(entityType, pos, width, height, animName, animOffset),
	}
}

// Returns a new ID counter that increments by one each time it is called.
func NewEntityIDCounter() func() (id int) {
	currentId := 0

	return func() int {
		currentId++

		return currentId
	}
}
