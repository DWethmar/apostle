package path

import (
	"github.com/dwethmar/apostle/point"
)

const Type = "Path"

type Path struct {
	entityID int
	cells    []point.P
	current  int // Index of the current cell in the path
}

func NewComponent(entityID int) *Path {
	return &Path{
		entityID: entityID,
		current:  0,
	}
}

func (p *Path) EntityID() int {
	return p.entityID
}

func (p *Path) ComponentType() string {
	return Type
}

func (p *Path) SetCells(cells []point.P) {
	p.cells = cells
	p.current = 0 // Reset current index when setting new cells
}

func (p *Path) AddCells(cells ...point.P) {
	p.cells = append(p.cells, cells...)
}

func (p *Path) Cells() []point.P {
	return p.cells
}

func (p *Path) CurrentCell() point.P {
	if p.current < len(p.cells) {
		return p.cells[p.current]
	}
	return point.P{} // Return zero value if no current cell
}

func (p *Path) Next() bool {
	if p.current+1 < len(p.cells) {
		p.current++
		return true
	}
	return false
}

func (p *Path) AtDestination() bool {
	return p.current >= len(p.cells)-1
}

// Destination returns the last cell in the path and a boolean indicating if a destination is set
func (p *Path) Destination() (point.P, bool) {
	if len(p.cells) == 0 {
		return point.P{}, false
	}
	return p.cells[len(p.cells)-1], true
}

// Reset resets the path to the beginning
func (p *Path) Reset() {
	p.current = 0
}

// Clear clears the path, removing all cells and resetting the current index
func (p *Path) Clear() {
	p.cells = make([]point.P, 0)
	p.current = 0
}

func (p *Path) ClearAfterNext() {
	if p.current+1 < len(p.cells) {
		p.cells = p.cells[:p.current+1]
	} else {
		p.cells = make([]point.P, 0)
	}
}
