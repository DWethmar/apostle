package world

import (
	"image/color"
	"log/slog"

	"github.com/dwethmar/apostle/component"
	"github.com/dwethmar/apostle/component/kind"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/event"
	"github.com/dwethmar/apostle/input"
	"github.com/dwethmar/apostle/propagation"
	"github.com/dwethmar/apostle/terrain"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const CellSize = 16 // Size of each cell in pixels

var (
	colorSolid  = color.RGBA{0, 128, 0, 255}
	colorBorder = color.RGBA{255, 0, 0, 255}
	colorEntity = color.RGBA{255, 255, 0, 255} // Yellow for entities
	colorPath   = color.RGBA{0, 0, 255, 255}   // Blue for paths
	colorApple  = color.RGBA{255, 0, 0, 255}   // Red for apples
)

type World struct {
	logger        *slog.Logger
	terrain       *terrain.Terrain
	entityStore   *entity.Store
	componenStore *component.Store
	eventBus      *event.Bus
}

func New(logger *slog.Logger, t *terrain.Terrain, entityStore *entity.Store, componenStore *component.Store, eventBus *event.Bus) *World {
	return &World{
		logger:        logger.With(slog.String("system", "world")),
		terrain:       t,
		entityStore:   entityStore,
		componenStore: componenStore,
		eventBus:      eventBus,
	}
}

func (d *World) OnPointerPressed(x, y int) propagation.Event {
	if err := d.eventBus.Publish(&input.Click{X: x, Y: y}); err != nil {
		d.logger.Error("failed to publish click event", slog.Int("x", x), slog.Int("y", y), slog.Any("error", err))
	}
	return propagation.Propagate
}

func (d *World) OnPointerReleased(x, y int) propagation.Event { return propagation.Propagate }

func (d *World) Draw(screen *ebiten.Image) {
	for step := range d.terrain.Walk() {
		x, y := step.X, step.Y
		cell := step.Cell

		if cell&terrain.Solid != 0 {
			// Draw solid cells as filled rectangles
			vector.FillRect(screen, float32(x*CellSize), float32(y*CellSize), CellSize, CellSize, colorSolid, false)
		}

		for _, border := range d.terrain.Walls(x, y) {
			switch border {
			case terrain.BorderNorth:
				vector.StrokeLine(screen, float32(x*CellSize), float32(y*CellSize), float32((x+1)*CellSize), float32(y*CellSize), 2, colorBorder, false)
			case terrain.BorderSouth:
				vector.StrokeLine(screen, float32(x*CellSize), float32((y+1)*CellSize), float32((x+1)*CellSize), float32((y+1)*CellSize), 2, colorBorder, false)
			case terrain.BorderEast:
				vector.StrokeLine(screen, float32((x+1)*CellSize), float32(y*CellSize), float32((x+1)*CellSize), float32((y+1)*CellSize), 2, colorBorder, false)
			case terrain.BorderWest:
				vector.StrokeLine(screen, float32(x*CellSize), float32(y*CellSize), float32(x*CellSize), float32((y+1)*CellSize), 2, colorBorder, false)
			}
		}
	}

	for _, e := range d.entityStore.Entities() {
		pos := e.Pos()
		x := float32(pos.X) * CellSize
		y := float32(pos.Y) * CellSize

		if m := e.Components().Movement(); m != nil {
			if !m.AtDestination() {
				progress := float32(m.CurrentStep()) / float32(m.Steps())
				endX := float32(m.Destination().X) * CellSize
				endY := float32(m.Destination().Y) * CellSize
				x += (endX - x) * progress
				y += (endY - y) * progress
			}
		}

		if k := e.Components().Kind(); k != nil {
			switch k.Value() {
			case kind.Human:
				drawEntityDiamond(screen, x, y)
			case kind.Apple:
				drawApple(screen, x, y)
			}
		}
	}

	for _, p := range d.componenStore.PathEntries() {
		drawPath(screen, p.Cells())
	}
}
