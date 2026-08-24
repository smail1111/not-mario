package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smail1111/mario/internal/objects"
)

func main() {
	ebiten.SetWindowTitle("Mario")

	ebiten.SetWindowSize(640, 480)

	er := ebiten.RunGame(&objects.Game{})

	if er != nil {
		log.Fatal(er)
	}
}
