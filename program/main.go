package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smail1111/mario/internal/objects"
	"github.com/smail1111/mario/internal/utils"
)

func main() {
	ebiten.SetWindowTitle("Not Mario")

	ebiten.SetWindowIcon([]image.Image{utils.MustLoadImage("images/icon16.png")})

	ebiten.SetWindowSize(640, 480)

	err := ebiten.RunGame(&objects.Game{})

	if err != nil {
		log.Fatal(err)
	}
}
