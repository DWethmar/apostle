package locomotion

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/dwethmar/apostle/component/movement"
	"github.com/dwethmar/apostle/component/path"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/point"
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
	logger      *slog.Logger
	entityStore *entity.Store
}

func New(logger *slog.Logger, entityStore *entity.Store) *Locomotion {
	return &Locomotion{
		logger:      logger,
		entityStore: entityStore,
	}
}

func (l *Locomotion) Update() error {
	for _, c := range l.entityStore.Components(movement.Type) {
		e, ok := l.entityStore.Entity(c.EntityID())
		if !ok {
			return fmt.Errorf("entity with ID %d does not exist", c.EntityID())
		}
		m, ok := c.(*movement.Movement)
		if !ok {
			return fmt.Errorf("component %T is not a Movement component", c)
		}
		logger := l.logger.With("component", c.Type(), "entityID", c.EntityID(), "currentStep", m.CurrentStep(), "steps", m.Steps(), "destination", m.Destination())

		// check if the entity has a path component
		if c, ok := l.entityStore.GetComponent(e.ID(), path.Type); ok {
			p, ok := c.(*path.Path)
			if !ok {
				return fmt.Errorf("component %T is not a Path component", c)
			}
			// If the entity has no destination, set it from the path component if ther path has cells
			if !m.HasDestination() {
				if len(p.Cells()) > 0 {
					steps := calculateSteps(e.Pos(), p.CurrentCell(), defaultStepSize)
					logger.Debug("Entity has no path destination, setting new destination", "currentPos", e.Pos, "destination", p.CurrentCell(), "steps", steps)
					m.SetDestination(p.CurrentCell(), steps) // Set new destination with calculated steps
				} else {
					m.SetDestination(e.Pos(), 0) // No path cells, stay at current position
				}
			} else {
				// If the entity is at its destination and the path has more cells, move to the next cell
				if m.AtDestination() && p.Next() {
					steps := calculateSteps(e.Pos(), p.CurrentCell(), defaultStepSize)
					logger.Debug("Entity at destination, moving to next path cell", "entityID", e.ID, "nextCell", p.CurrentCell(), "steps", steps)
					m.SetDestination(p.CurrentCell(), steps) // Set new destination with calculated steps
				}
			}
		}

		if !m.AtDestination() {
			m.AdvanceStep()
		}

		if m.AtDestination() && !e.Pos().Equal(m.Destination()) { // Reached destination and we didn't update the entity position yet
			logger.Debug("Entity reached destination")
			e.SetPos(m.Destination())
		}
	}
	return nil
}
