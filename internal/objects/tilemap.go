package objects

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smail1111/not-mario/internal/utils"
)

type Tile struct {
	Pos     Vector `json:"pos"`
	Type    string `json:"type"`
	Variant int    `json:"variant"`
}

// Returns an FRect based on the Tile.
func (t Tile) FRect(size int) FRect {
	return NewFRect(float64(t.Pos.X*size), float64(t.Pos.Y*size), float64(size), float64(size))
}

// Draws a Tile on the given image based on the given tile size and offset.
func (t *Tile) DrawGrid(screen, image *ebiten.Image, tileSize int, offset FVector) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(t.Pos.X)*float64(tileSize)-offset.X, float64(t.Pos.Y)*float64(tileSize)-offset.Y)
	screen.DrawImage(image, &opts)
}

// Draws a Tile on the given image based on the given offset.
func (t *Tile) DrawOffGrid(screen, image *ebiten.Image, offset FVector) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(t.Pos.X)-offset.X, float64(t.Pos.Y)-offset.Y)
	screen.DrawImage(image, &opts)
}

// Returns a formatted key based on the Tile's type and variant.
func (t *Tile) ImageKey() string {
	return fmt.Sprintf("tiles/%s/%d.png", t.Type, t.Variant)
}

type TileMap struct {
	Tiles    map[string]Tile  `json:"tiles"`
	Blocks   map[string]Block `json:"-"`
	Offgrid  []Tile           `json:"offgrid"`
	TileSize int              `json:"tile_size"`
}

// Draws a TileMap based on the given assets and offset.
func (tm TileMap) Draw(screen *ebiten.Image, assets map[string]*ebiten.Image, offset FVector) {
	for _, tile := range tm.Offgrid {
		tile.DrawOffGrid(screen, assets[tile.ImageKey()], offset)
	}

	topLeft := Vector{utils.Floor(offset.X) / tm.TileSize, utils.Floor(offset.Y) / tm.TileSize}
	for x := topLeft.X - 1; x < topLeft.X+screen.Bounds().Dx()/tm.TileSize+1; x++ {
		for y := topLeft.Y - 1; y < topLeft.Y+screen.Bounds().Dy()/tm.TileSize+1; y++ {
			coord := Vector{x, y}

			if tile, ok := tm.Tiles[coord.String()]; ok {
				tile.DrawGrid(screen, assets[tile.ImageKey()], tm.TileSize, offset)
			}

			if block, ok := tm.Blocks[coord.String()]; ok {
				block.Draw(screen, offset)
			}
		}
	}
}

// Loads a TileMap from the given name.
func (tm *TileMap) MustLoad(name string) {
	bytes := utils.MustReadMap(name)

	er := json.Unmarshal(bytes, tm)
	if er != nil {
		log.Fatal(er)
	}
}

type TileID struct {
	Type    string
	Variant int
}

// Extracts every Tile in a TileMap that has an ID in the given ids,
// removing or keeping the extracted tiles from the TileMap based on keep.
func (tm *TileMap) Extract(ids utils.Set, keep bool) (got []Tile) {
	for key, tile := range tm.Tiles {
		if ids.Has(TileID{tile.Type, tile.Variant}) {
			gridTile := tile

			gridTile.Pos.X *= tm.TileSize
			gridTile.Pos.Y *= tm.TileSize

			got = append(got, gridTile)
			if !keep {
				defer delete(tm.Tiles, key)
			}
		}
	}

	newOffgridTiles := []Tile{}
	for _, tile := range tm.Offgrid {
		if ids.Has(TileID{tile.Type, tile.Variant}) {
			got = append(got, tile)
			if keep {
				newOffgridTiles = append(newOffgridTiles, tile)
			}
		} else {
			newOffgridTiles = append(newOffgridTiles, tile)
		}
	}
	tm.Offgrid = newOffgridTiles

	return
}

// Returns whether there is a solid TIle at the given x and y position in the TileMap.
func (tm TileMap) CheckSolidTile(x float64, y float64) (exists bool) {
	coord := Vector{utils.FloorDiv(x, tm.TileSize), utils.FloorDiv(y, tm.TileSize)}

	tile, ok := tm.Tiles[coord.String()]
	if ok && collisionTypes.Has(tile.Type) {
		return true
	}

	block, ok := tm.Blocks[coord.String()]
	return ok && collisionTypes.Has(block.Type)
}

// Returns a new TileMap with the given tile size loaded from the given map name if the given map name is not an empty string.
func (g *Game) NewTileMap(tileSize int, mapName string) TileMap {
	tileMap := TileMap{
		Tiles:    map[string]Tile{},
		Blocks:   map[string]Block{},
		Offgrid:  []Tile{},
		TileSize: tileSize,
	}

	if mapName != "" {
		tileMap.MustLoad(mapName)
	}

	return tileMap
}
