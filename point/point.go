package point

// P represents a point in 2D space.
type P struct {
	X, Y int
}

// New creates a new point with the given coordinates.
func New(x, y int) P {
	return P{X: x, Y: y}
}

func (p P) Equal(other P) bool {
	return p.X == other.X && p.Y == other.Y
}

// Neighboring checks if two points are adjacent (horizontally, vertically and diagonally).
func (p P) Neighboring(other P) bool {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return (dx == 0 && (dy == 1 || dy == -1)) || (dy == 0 && (dx == 1 || dx == -1)) || (dx == 1 && (dy == 1 || dy == -1)) || (dx == -1 && (dy == 1 || dy == -1))
}

func (p P) Divide(scalar int) P {
	return P{X: p.X / scalar, Y: p.Y / scalar}
}
