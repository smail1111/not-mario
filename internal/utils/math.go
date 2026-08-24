package utils

import (
	"math"
	"math/rand/v2"
)

// Floors f and returns an int.
func Floor(f float64) int {
	if f < 0.0 {
		return int(f) - 1
	} else {
		return int(f)
	}
}

// Returns the remainder of f floored divided by mod as a float.
func ModFloat(f float64, mod int) float64 {
	return float64(Floor(f) % mod)
}

// Returns a as a float divided by b as a float floored.
func FloorDiv(f float64, div int) int {
	return Floor(f / float64(div))
}

// Returns the x and y coordinates of a random velocity with a speed in between min and max.
// Max should be greater than or equal to min.
func RandVelocity(min float64, max float64) [2]float64 {
	angle := rand.Float64() * math.Pi * 2

	speed := rand.Float64()*(max-min) + min

	return [2]float64{math.Cos(angle) * speed, math.Sin(angle) * speed}
}
