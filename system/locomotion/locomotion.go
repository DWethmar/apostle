package locomotion

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/dwethmar/apostle/component"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/point"
	"github.com/dwethmar/apostle/system/world"
)

const defaultStepSize = 20 // Default step size for movement

// calculateSteps calculates the number of steps needed to move from start to end
// Uses Euclidean distance to ensure diagonal movement isn't faster than axis-aligned movement
func calculateSteps(start, end point.P, stepsPerUnit int) int {
	dx := float64(end.X - start.X)
	dy := float64(end.Y - start.Y)
	distance := math.Sqrt(dx*dx + dy*dy)
	return int(math.Ceil(distance * float64(stepsPerUnit)))
}

// Locomotion handles the movement of entities based on their paths and movement components.
type Locomotion struct {
	logger         *slog.Logger
	entityStore    *entity.Store
	componentStore *component.Store
}

func New(logger *slog.Logger, entityStore *entity.Store, componenStore *component.Store) *Locomotion {
	return &Locomotion{
		logger:         logger.With(slog.String("system", "locomotion")),
		entityStore:    entityStore,
		componentStore: componenStore,
	}
}

func (l *Locomotion) Update() error {
	for _, m := range l.componentStore.MovementEntries() {
		e, ok := l.entityStore.Entity(m.EntityID())
		if !ok {
			return fmt.Errorf("entity with ID %d does not exist", m.EntityID())
		}

		if !m.HasDestination() {
			m.SetDestinationCell(world.PXToCell(e.Pos()), 0) // Set current position as destination with 0 steps
		}

		// check if the entity has a path component
		if p := e.Components().Path(); p != nil {
			if _, hasDest := p.Destination(); !hasDest {
				break
			}

			// If the entity is at its destination and the path has more cells, move to the next cell
			if m.AtDestination() && p.Next() {
				steps := calculateSteps(m.OriginCell(), m.DestinationCell(), defaultStepSize)
				m.SetDestinationCell(p.CurrentCell(), steps) // Set new destination with calculated steps
			}
		}

		if !m.AtDestination() {
			m.AdvanceStep()
			cellSize := float32(world.CellSize)
			progress := float32(m.CurrentStep()) / float32(m.Steps())
			newX := float32(m.OriginCell().X)*(1-progress) + float32(m.DestinationCell().X)*progress
			newY := float32(m.OriginCell().Y)*(1-progress) + float32(m.DestinationCell().Y)*progress
			e.SetPos(point.P{
				X: int(newX*cellSize) + world.CellSize/2,
				Y: int(newY*cellSize) + world.CellSize/2,
			})
		}
	}
	return nil
}
