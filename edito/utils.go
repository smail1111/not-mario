package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smail1111/mario/assets"
	"github.com/smail1111/mario/internal/objects"
	"github.com/smail1111/mario/internal/utils"
)

var Assets = assets.Assets

func MustLoadDir(dirpath string) []*ebiten.Image {
	entries, er := Assets.Images.ReadDir(dirpath)
	if er != nil {
		log.Fatal(er)
	}

	images := []*ebiten.Image{}

	for _, entry := range entries {
		image := utils.MustLoadImage(dirpath + "/" + entry.Name())
		images = append(images, image)
	}

	return images
}

func MustLoadDirs(dirpath string) map[string][]*ebiten.Image {
	entries, er := Assets.Images.ReadDir(dirpath)
	if er != nil {
		log.Fatal(er)
	}

	dirs := map[string][]*ebiten.Image{}

	for _, entry := range entries {
		dirs[entry.Name()] = MustLoadDir(dirpath + "/" + entry.Name())
	}

	return dirs
}

func Save(tilemap objects.TileMap, path string) {
	file, er := os.Create(path)
	if er != nil {
		log.Fatal(er)
	}

	defer file.Close()

	bytes, er := json.Marshal(tilemap)
	if er != nil {
		log.Fatal(er)
	}

	_, er = file.Write(bytes)
	if er != nil {
		log.Fatal(er)
	}
}

func AutoTile(tm *objects.TileMap) {
	type Dirs struct {
		left  bool
		right bool
		up    bool
		down  bool
	}

	autoTileMap := map[Dirs]int{
		{down: true, right: true}:           0,
		{down: true, left: true}:            1,
		{up: true, right: true}:             2,
		{up: true, left: true}:              3,
		{up: true, down: true, right: true}: 2,
		{up: true, down: true, left: true}:  3,
	}

	autoTileTypes := utils.NewSet("pipe")

	for key, tile := range tm.Tiles {
		if autoTileTypes.Has(tile.Type) {
			coord := tile.Pos
			connections := Dirs{}

			for _, offset := range []objects.Vector{{X: -1, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1}} {
				key := fmt.Sprintf("%d,%d", coord.X+offset.X, coord.Y+offset.Y)
				if neighbor, ok := tm.Tiles[key]; ok && neighbor.Type == tile.Type {
					if offset.X == -1 && offset.Y == 0 {
						connections.left = true
					} else if offset.X == 1 && offset.Y == 0 {
						connections.right = true
					} else if offset.X == 0 && offset.Y == -1 {
						connections.up = true
					} else if offset.X == 0 && offset.Y == 1 {
						connections.down = true
					}
				}
			}

			tm.Tiles[key] = objects.Tile{
				Type:    tile.Type,
				Variant: autoTileMap[connections],
				Pos:     coord,
			}
		}
	}
}
