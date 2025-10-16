package direction

type (
	Direction int
)

// Direction flags for movement
const (
	None Direction = 1 << iota
	North
	South
	East
	West
	NorthEast
	NorthWest
	SouthEast
	SouthWest
)
