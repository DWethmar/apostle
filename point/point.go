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

func (p P) Neighboring(other P) bool {
	return (p.X == other.X && (p.Y == other.Y-1 || p.Y == other.Y+1)) ||
		(p.Y == other.Y && (p.X == other.X-1 || p.X == other.X+1))
}
