package utils

import (
	"image"
	"log"
	"strings"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smail1111/not-mario/assets"
)

var Assets = assets.Assets

// The base path where image assets will be contained.
const ImagesPath = "images/"

// Returns a new *ebitien.Image from the given filepath.
// If an error occurs while loading, log the error and exit with status 1.
func MustLoadImage(filepath string) *ebiten.Image {
	f, er := Assets.Images.Open(filepath)
	if er != nil {
		log.Fatal(er)
	}

	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(er)
	}

	return ebiten.NewImageFromImage(img)
}

// Recursively loads every png image from a given directory and returns a map mapping each file's path to an *ebiten.Image.
// If an error occurs while loading, log the error and exit with status 1.
func MustLoadImages(dirpath string, imagesMap map[string]*ebiten.Image) map[string]*ebiten.Image {
	entries, er := Assets.Images.ReadDir(dirpath)
	if er != nil {
		log.Fatal(er)
	}

	dirpath += "/"

	for _, entry := range entries {
		if entry.IsDir() {
			MustLoadImages(dirpath+entry.Name(), imagesMap)
		} else {
			image := MustLoadImage(dirpath + entry.Name())
			key := strings.TrimPrefix(dirpath, ImagesPath) + entry.Name()
			imagesMap[key] = image
		}
	}

	return imagesMap
}

// Manually changes the alpha value of each pixel in the given image which color matches the given RGB values to 0.
// Must be called while the game is running.
func SetColorKey(image *ebiten.Image, red, green, blue byte) {
	colors := []byte{red, green, blue}

	imgBytes := make([]byte, 4*image.Bounds().Dx()*image.Bounds().Dy())
	image.ReadPixels(imgBytes)

	i := 0
	for i <= len(imgBytes)-4 {
		for j := range 4 {
			if j == 3 {
				for j = range 4 {
					imgBytes[i+j] = 0
				}
				i += 4
			} else if imgBytes[i+j] != colors[j] {
				i += 4
				break
			}
		}
	}

	image.WritePixels(imgBytes)
}

// Changes the alpha value of each pixel in the given image with an alpha value not equal to zero to alpha.
// Must be called while the game is running.
func SetAlpha(image *ebiten.Image, aplha byte) {
	imgBytes := make([]byte, 4*image.Bounds().Dx()*image.Bounds().Dy())
	image.ReadPixels(imgBytes)

	i := 3
	for i < len(imgBytes) {
		if imgBytes[i] != 0 {
			imgBytes[i] = aplha
		}
		i += 4
	}

	image.WritePixels(imgBytes)
}

// Returns a copy of the image with every color value set to zero.
func Silhouette(image *ebiten.Image) *ebiten.Image {
	imgBytes := make([]byte, 4*image.Bounds().Dx()*image.Bounds().Dy())
	image.ReadPixels(imgBytes)

	for i := range imgBytes {
		if i%4 != 3 {
			imgBytes[i] = 0
		}
	}

	silhouette := ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
	silhouette.WritePixels(imgBytes)

	return silhouette
}
