package drawer

import (
	"image/color"

	"github.com/dwethmar/apostle/component/movement"
	"github.com/dwethmar/apostle/component/path"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/terrain"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const cellSize = 16 // Size of each cell in pixels

var (
	colorSolid  = color.RGBA{0, 128, 0, 255}
	colorBorder = color.RGBA{255, 0, 0, 255}
	colorEntity = color.RGBA{255, 255, 0, 255} // Yellow for entities
	colorPath   = color.RGBA{0, 0, 255, 255}   // Blue for paths
)

type Drawer struct {
	terrain     *terrain.Terrain
	entityStore *entity.Store
}

func New(t *terrain.Terrain, entityStore *entity.Store) *Drawer {
	return &Drawer{
		terrain:     t,
		entityStore: entityStore,
	}
}

func (d *Drawer) Draw(screen *ebiten.Image) {
	for step := range d.terrain.Walk() {
		x, y := step.X, step.Y
		cell := step.Cell

		if cell&terrain.Solid != 0 {
			// Draw solid cells as filled rectangles
			vector.FillRect(screen, float32(x*cellSize), float32(y*cellSize), cellSize, cellSize, colorSolid, false)
		}

		for _, border := range d.terrain.Walls(x, y) {
			switch border {
			case terrain.BorderNorth:
				vector.StrokeLine(screen, float32(x*cellSize), float32(y*cellSize), float32((x+1)*cellSize), float32(y*cellSize), 2, colorBorder, false)
			case terrain.BorderSouth:
				vector.StrokeLine(screen, float32(x*cellSize), float32((y+1)*cellSize), float32((x+1)*cellSize), float32((y+1)*cellSize), 2, colorBorder, false)
			case terrain.BorderEast:
				vector.StrokeLine(screen, float32((x+1)*cellSize), float32(y*cellSize), float32((x+1)*cellSize), float32((y+1)*cellSize), 2, colorBorder, false)
			case terrain.BorderWest:
				vector.StrokeLine(screen, float32(x*cellSize), float32(y*cellSize), float32(x*cellSize), float32((y+1)*cellSize), 2, colorBorder, false)
			}
		}
	}

	for _, e := range d.entityStore.Entities() {
		pos := e.Pos()
		x := float32(pos.X) * cellSize
		y := float32(pos.Y) * cellSize

		if c, ok := d.entityStore.GetComponent(e.ID(), movement.Type); ok {
			m, ok := c.(*movement.Movement)
			if ok && !m.AtDestination() {
				progress := float32(m.CurrentStep()) / float32(m.Steps())
				endX := float32(m.Destination().X) * cellSize
				endY := float32(m.Destination().Y) * cellSize
				x += (endX - x) * progress
				y += (endY - y) * progress
			}
		}

		drawEntityDiamond(screen, x, y)
	}

	for _, p := range d.entityStore.Components("Path") {
		if pathComp, ok := p.(*path.Path); ok {
			drawPath(screen, pathComp.Cells())
		}
	}
}
