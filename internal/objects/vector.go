package objects

import (
	"fmt"

	"github.com/smail1111/mario/internal/utils"
)

type Vector struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type FVector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector) FVector() FVector {
	return FVector{float64(v.X), float64(v.Y)}
}

func (fv FVector) Vector() Vector {
	return Vector{utils.Floor(fv.X), utils.Floor(fv.Y)}
}

func (v Vector) String() string {
	return fmt.Sprintf("%d,%d", v.X, v.Y)
}
