package path

import (
	"github.com/dwethmar/apostle/component"
	"github.com/dwethmar/apostle/point"
)

type Path struct {
	*component.Component
	cells   []point.P
	current int // Index of the current cell in the path
}

func NewComponent(entityID int) *Path {
	return &Path{
		Component: &component.Component{
			EID: entityID,
			T:   "Path",
		},
		cells:   make([]point.P, 0),
		current: 0,
	}
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

func (p *Path) NextCell() bool {
	if p.current+1 < len(p.cells) {
		p.current++
		return true
	}
	return false
}

func (p *Path) AtDestination() bool {
	return p.current >= len(p.cells)-1
}

func (p *Path) Reset() {
	p.current = 0
	p.cells = make([]point.P, 0) // Clear the path
}
