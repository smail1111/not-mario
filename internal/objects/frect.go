package objects

type FRect struct {
	Pos    FVector
	Width  float64
	Height float64
}

// Return whether an FRect overlaps another FRect.
func (fr FRect) Overlaps(other FRect) bool {
	return fr.Left() < other.Right() && fr.Right() > other.Left() &&
		fr.Top() < other.Bottom() && fr.Bottom() > other.Top()
}

// Return whether an FRect overlaps a FVector2.
func (fr FRect) CollidesPoint(pt FVector) bool {
	return pt.X < fr.Right() && pt.X > fr.Left() &&
		pt.Y < fr.Bottom() && pt.Y > fr.Top()
}

// Return the leftmost position of an FRect.
func (fr FRect) Left() float64 {
	return fr.Pos.X
}

// Return the rightmost position of an FRect.
func (fr FRect) Right() float64 {
	return fr.Pos.X + fr.Width
}

// Return the uppermost position of an FRect.
func (fr FRect) Top() float64 {
	return fr.Pos.Y
}

// Return the bottommost position of an FRect.
func (fr FRect) Bottom() float64 {
	return fr.Pos.Y + fr.Height
}

// Return the center of an FRect.
func (fr FRect) Center() FVector {
	return FVector{fr.Left() + fr.Width/2.0, fr.Top() + fr.Height/2.0}
}

// Return the topleft point of an FRect.
func (fr FRect) TopLeft() FVector {
	return FVector{fr.Left(), fr.Top()}
}

// Return the topright point of an FRect.
func (fr FRect) BottomRight() FVector {
	return FVector{fr.Right(), fr.Bottom()}
}

// Increase an FRect's x and y coordinates by x and y.
func (fr *FRect) Add(x float64, y float64) {
	fr.Pos.X += x
	fr.Pos.Y += y
}

// Return a new FRect with the given x coordinate, y coordinate, width, and height.
func NewFRect(x, y, width, height float64) FRect {
	return FRect{
		Pos:    FVector{x, y},
		Width:  width,
		Height: height,
	}
}
