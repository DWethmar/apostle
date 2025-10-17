package terrain

import (
	"fmt"
	"iter"

	"github.com/dwethmar/apostle/direction"
	"github.com/dwethmar/apostle/point"
)

type (
	Cell byte // 1 byte to store flags for solidity and borders
)

const (
	width  = 50
	height = 50
)

// Solidity, border and other flags for terrain cells
const (
	Solid       Cell = 1 << iota // 00000001
	BorderNorth                  // 00000010
	BorderSouth                  // 00000100
	BorderWest                   // 00001000
	BorderEast                   // 00010000
	Ceiling                      // 00100000
	Floor                        // 01000000
	_                            // 10000000 (unused)
)

type Terrain struct {
	cells [width][height]Cell
}

func New() *Terrain {
	return &Terrain{
		cells: [width][height]Cell{},
	}
}

func (t *Terrain) InBounds(x, y int) bool {
	// calc if x and y are within the bounds of the terrain
	return x >= 0 && x < width && y >= 0 && y < height
}

func (t *Terrain) Fill(x, y int, cell Cell) error {
	if !t.InBounds(x, y) {
		return fmt.Errorf("coordinates exceed bounds: (%d, %d) out of (%d, %d)", x, y, width, height)
	}
	t.cells[y][x] = cell
	return nil
}

func (t *Terrain) Solid(x, y int) bool {
	if !t.InBounds(x, y) {
		return false
	}
	return t.HasFlag(x, y, Solid)
}

func (t *Terrain) HasFlag(x, y int, flag Cell) bool {
	return t.InBounds(x, y) && (t.cells[y][x]&flag != 0)
}

func (t *Terrain) HasCeiling(x, y int) bool {
	return t.HasFlag(x, y, Ceiling)
}

func (t *Terrain) HasFloor(x, y int) bool {
	return t.HasFlag(x, y, Floor)
}

var borderFlags = []Cell{BorderNorth, BorderSouth, BorderEast, BorderWest}

func (t *Terrain) Walls(x, y int) []Cell {
	if !t.InBounds(x, y) {
		return nil
	}
	cell := t.cells[y][x]
	borders := make([]Cell, 0, 4)
	for _, flag := range borderFlags {
		if cell&flag != 0 {
			borders = append(borders, flag)
		}
	}
	return borders
}

func (t *Terrain) Width() int {
	return width
}

func (t *Terrain) Height() int {
	return height
}

type moveInfo struct {
	dx, dy                      int
	currentBorder, targetBorder Cell
}

var moves = map[direction.Direction]moveInfo{
	direction.North:     {0, -1, BorderNorth, BorderSouth},
	direction.South:     {0, 1, BorderSouth, BorderNorth},
	direction.East:      {1, 0, BorderEast, BorderWest},
	direction.West:      {-1, 0, BorderWest, BorderEast},
	direction.NorthEast: {1, -1, BorderNorth, BorderSouth},
	direction.NorthWest: {-1, -1, BorderNorth, BorderSouth},
	direction.SouthEast: {1, 1, BorderSouth, BorderNorth},
	direction.SouthWest: {-1, 1, BorderSouth, BorderNorth},
}

// Traversable checks if a point is traversable in a given direction.
func (t *Terrain) Traversable(p point.P, d direction.Direction) bool {
	move := moves[d]
	newX, newY := p.X+move.dx, p.Y+move.dy

	if !t.InBounds(newX, newY) {
		return false
	}

	if t.cells[p.Y][p.X]&move.currentBorder != 0 {
		return false
	}
	if t.cells[newY][newX]&move.targetBorder != 0 {
		return false
	}

	return t.cells[newY][newX]&Solid == 0
}

// Step represents a position and the cell at that position during a walk through the terrain.
type Step struct {
	X, Y int
	Cell Cell
}

func (t *Terrain) Walk() iter.Seq[Step] {
	return func(yield func(Step) bool) {
		for y := range height {
			for x := range width {
				if !yield(Step{X: x, Y: y, Cell: t.cells[y][x]}) {
					return
				}
			}
		}
	}
}
