package direction

type (
	Direction int
)

// Direction flags for movement
const (
	North Direction = 1 << iota
	South
	East
	West
	NorthEast
	NorthWest
	SouthEast
	SouthWest
)
