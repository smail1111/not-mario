package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/smail1111/not-mario/internal/objects"
	"github.com/smail1111/not-mario/internal/utils"
)

const scale = 2.0

var (
	loadFrom string
	saveTo   = "map"
)

type Game struct {
	TileMap    objects.TileMap
	ListAssets map[string][]*ebiten.Image
	MapAssets  map[string]*ebiten.Image
	Types      []string
	GridTile   objects.Tile
	Offset     objects.FVector
	Tile       objects.Vector
	Movement   objects.Vector
	Display    *ebiten.Image
	offGrid    bool
	collisions bool
	saved      bool
	init       bool
}

// Initiates the editor.
func (g *Game) Init() {
	g.ListAssets = MustLoadDirs("images/tiles")
	for _, group := range g.ListAssets {
		for _, image := range group {
			utils.SetColorKey(image, 146, 144, 255)
			utils.SetColorKey(image, 92, 148, 252)
		}
	}

	g.MapAssets = utils.MustLoadImages("images/tiles", map[string]*ebiten.Image{})
	for _, image := range g.MapAssets {
		utils.SetColorKey(image, 146, 144, 255)
		utils.SetColorKey(image, 92, 148, 252)
	}

	g.Display = ebiten.NewImage(320, 240)

	g.TileMap = objects.TileMap{
		Tiles:    map[string]objects.Tile{},
		Offgrid:  []objects.Tile{},
		TileSize: 16,
	}

	if loadFrom != "" {
		g.TileMap.MustLoad(loadFrom)
	}

	g.Types = []string{}
	for tileType := range g.ListAssets {
		g.Types = append(g.Types, tileType)
	}

	g.init = true
}

// Updates the editor.
func (g *Game) Update() error {
	// Initiate the editor.
	if !g.init {
		g.Init()
	}

	// Save the TileMap.
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyS) {
		Save(g.TileMap, fmt.Sprintf("assets/maps/%s.json", saveTo))

		g.saved = true
	}

	// Swap offgrid mode.
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.offGrid = !g.offGrid
	}

	// Scaled mouse position.
	cursorX, cursorY := ebiten.CursorPosition()
	cursorX /= scale
	cursorY /= scale

	if g.offGrid {
		// Create an offgrid tile.
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			g.TileMap.Offgrid = append(g.TileMap.Offgrid, objects.Tile{
				Type:    g.Types[g.Tile.X],
				Variant: g.Tile.Y,
				Pos:     objects.Vector{X: cursorX + int(g.Offset.X), Y: cursorY + int(g.Offset.Y)},
			})
			g.saved = false
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
			// Delete an offgrid tile.
			newOffgrid := []objects.Tile{}
			for _, tile := range g.TileMap.Offgrid {
				tileImage := g.ListAssets[tile.Type][tile.Variant]
				if frect := objects.NewFRect(
					float64(tile.Pos.X),
					float64(tile.Pos.Y),
					float64(tileImage.Bounds().Dx()),
					float64(tileImage.Bounds().Dy()),
				); !frect.CollidesPoint(objects.FVector{X: float64(cursorX) + g.Offset.X, Y: float64(cursorY) + g.Offset.Y}) {
					newOffgrid = append(newOffgrid, tile)
				}
			}
			g.TileMap.Offgrid = newOffgrid
		}
	} else {
		// Calculate the current grid coordinate based on the mouse position and offset.
		coord := objects.Vector{X: (utils.FloorDiv(float64(cursorX+utils.Floor(g.Offset.X)), g.TileMap.TileSize)),
			Y: utils.FloorDiv(float64(cursorY+utils.Floor(g.Offset.Y)), g.TileMap.TileSize)}

		// Create a tile from the coordinate and selected tile.
		g.GridTile = objects.Tile{
			Type:    g.Types[g.Tile.X],
			Variant: g.Tile.Y,
			Pos:     coord,
		}

		// Create a grid tile.
		if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			g.TileMap.Tiles[coord.String()] = g.GridTile
			AutoTile(&g.TileMap)
			g.saved = false
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
			// Delete a grid tile.
			delete(g.TileMap.Tiles, coord.String())
			AutoTile(&g.TileMap)
			g.saved = false
		}
	}

	// Mouse wheel movement.
	_, mouseWheelMov := ebiten.Wheel()

	// Change tile variant.
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		g.Tile.Y = (g.Tile.Y + int(mouseWheelMov)) % len(g.ListAssets[g.Types[g.Tile.X]])

		if g.Tile.Y < 0 {
			g.Tile.Y = len(g.ListAssets[g.Types[g.Tile.X]]) - 1
		}
	} else if mouseWheelMov != 0.0 {
		// Change tile type.
		g.Tile.X = (g.Tile.X + int(mouseWheelMov)) % len(g.Types)

		if g.Tile.X < 0 {
			g.Tile.X = len(g.Types) + g.Tile.X
		}

		g.Tile.Y = 0
	}

	// Scroll screen.
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Movement.X = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Movement.X = 1
	} else {
		g.Movement.X = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Movement.Y = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyS) && !ebiten.IsKeyPressed(ebiten.KeyControl) {
		g.Movement.Y = 1
	} else {
		g.Movement.Y = 0
	}

	// Update offset.
	g.Offset.X += float64(g.Movement.X * 5.0)
	g.Offset.Y += float64(g.Movement.Y * 5.0)

	// End.
	return nil
}

// Draws the editor.
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background.
	g.Display.Fill(color.RGBA{92, 148, 252, 255})

	// Draw tilemap.
	g.TileMap.Draw(g.Display, g.MapAssets, g.Offset)

	// Get current tile image.
	tile_img := g.ListAssets[g.Types[g.Tile.X]][g.Tile.Y]
	utils.SetAlpha(tile_img, 120)

	// Draw selected tile.
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(1.5, 1.5)
	opts.GeoM.Translate(5.0, 5.0)
	g.Display.DrawImage(tile_img, &opts)

	// Draw preview tile.
	opts = ebiten.DrawImageOptions{}
	if g.offGrid {
		opts = ebiten.DrawImageOptions{}
		x, y := ebiten.CursorPosition()
		opts.GeoM.Translate(float64(x)/scale, (float64(y))/scale)
		g.Display.DrawImage(tile_img, &opts)
	} else {
		g.GridTile.DrawGrid(g.Display, tile_img, 16, g.Offset)
	}

	if g.saved {
		ebitenutil.DebugPrintAt(g.Display, "Saved!", 5, 5)
	}

	// Draw display.
	opts = ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scale, scale)

	screen.DrawImage(g.Display, &opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
