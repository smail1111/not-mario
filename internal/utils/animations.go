package utils

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	Images   []*ebiten.Image
	Name     string
	Duration int
	Frame    int
	Loop     bool
	Done     bool
}

// Increases an animation's frame by one, looping or stopping based on the animation's loop.
func (a *Animation) Update() {
	if a.Loop {
		a.Frame = (a.Frame + 1) % (len(a.Images) * a.Duration)
	} else {
		a.Frame = min(a.Frame+1, len(a.Images)*a.Duration-1)
		if a.Frame >= len(a.Images)*a.Duration-1 {
			a.Done = true
		}
	}
}

// Gets the current image from an animation based on the animation's current frame.
func (a Animation) Image() *ebiten.Image {
	return a.Images[a.Frame/a.Duration]
}

// Returns a new animation with the given images, name, duration, and loop.
func NewAnimation(images []*ebiten.Image, name string, duration int, loop bool) Animation {
	return Animation{
		Images:   images,
		Name:     name,
		Duration: duration,
		Loop:     loop,
	}
}

// Returns a new animation with images loaded from the given dirpath with the given name, duration, and loop.
func MustLoadAnimation(dirpath string, name string, duration int, loop bool) Animation {
	entries, er := Assets.Images.ReadDir(dirpath)
	if er != nil {
		log.Fatal(er)
	}

	images := []*ebiten.Image{}

	dirpath += "/"

	for _, entry := range entries {
		images = append(images, MustLoadImage(dirpath+entry.Name()))
	}

	return NewAnimation(images, name, duration, loop)
}
