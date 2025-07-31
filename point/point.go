package point

// P represents a point in 2D space.
type P struct {
	X, Y int
}

func (p P) Equal(other P) bool {
	return p.X == other.X && p.Y == other.Y
}
