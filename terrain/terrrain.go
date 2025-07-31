package terrain

import (
	"fmt"
	"iter"
)

type (
	Direction int
	Cell      byte // 1 byte to store flags for solidity and borders
)

const (
	width  = 20
	height = 20
)

// Direction flags for movement
const (
	North Direction = 1 << iota
	South
	East
	West
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

func inbounds(width, height, x, y int) bool {
	return x >= 0 && x < width && y >= 0 && y < height
}

func (t *Terrain) Fill(x, y int, cell Cell) error {
	if !inbounds(width, height, x, y) {
		return fmt.Errorf("coordinates exceed bounds: (%d, %d) out of (%d, %d)", x, y, width, height)
	}
	t.cells[y][x] = cell
	return nil
}

func (t *Terrain) Solid(x, y int) bool {
	if !inbounds(width, height, x, y) {
		return false
	}
	return t.cells[y][x]&Solid != 0
}

func (t *Terrain) HasFlag(x, y int, flag Cell) bool {
	return inbounds(width, height, x, y) && (t.cells[y][x]&flag != 0)
}

func (t *Terrain) HasCeiling(x, y int) bool {
	return t.HasFlag(x, y, Ceiling)
}

func (t *Terrain) HasFloor(x, y int) bool {
	return t.HasFlag(x, y, Floor)
}

var borderFlags = []Cell{BorderNorth, BorderSouth, BorderEast, BorderWest}

func (t *Terrain) Walls(x, y int) []Cell {
	if !inbounds(width, height, x, y) {
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

type moveInfo struct {
	dx, dy                      int
	currentBorder, targetBorder Cell
}

var moves = map[Direction]moveInfo{
	North: {0, -1, BorderNorth, BorderSouth},
	South: {0, 1, BorderSouth, BorderNorth},
	East:  {1, 0, BorderEast, BorderWest},
	West:  {-1, 0, BorderWest, BorderEast},
}

func (t *Terrain) Traversable(x, y int, d Direction) bool {
	move := moves[d]
	newX, newY := x+move.dx, y+move.dy

	if !inbounds(width, height, newX, newY) {
		return false
	}

	if t.cells[y][x]&move.currentBorder != 0 {
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
