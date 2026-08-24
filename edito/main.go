package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/silbinarywolf/preferdiscretegpu"
)

func init() {
	for i, arg := range os.Args {
		if len(os.Args) > i+1 {
			switch arg {
			case "--to":
				saveTo = os.Args[i+1]

			case "--load":
				loadFrom = os.Args[i+1]
			}
		}
	}
}

func main() {
	ebiten.SetWindowSize(640, 480)

	ebiten.SetWindowTitle("Edito")

	g := Game{}

	er := ebiten.RunGame(&g)

	if er != nil {
		log.Fatal(er)
	}
}
