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

const defaultStepSize = 100 // Default step size for movement

// calculateSteps calculates the number of steps needed to move from start to end
// Uses Euclidean distance to ensure diagonal movement isn't faster than axis-aligned movement
func calculateSteps(start, end point.P, stepsPerUnit int) int {
	dx := float64(end.X - start.X)
	dy := float64(end.Y - start.Y)
	distance := math.Sqrt(dx*dx + dy*dy)
	return int(math.Ceil(distance * float64(stepsPerUnit)))
}

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
	for _, c := range l.entityStore.Components("Movement") {
		e := l.entityStore.Entity(c.EntityID())
		m, ok := c.(*movement.Movement)
		if !ok {
			return fmt.Errorf("component %T is not a Movement component", c)
		}
		logger := l.logger.With("component", c.Type(), "entityID", c.EntityID(), "currentStep", m.CurrentStep(), "steps", m.Steps(), "destination", m.Destination())

		hasAdvanced := false
		if !m.AtDestination() {
			// logger.Debug("Advancing step for entity", "currentPos", e.Pos)
			m.AdvanceStep()
			hasAdvanced = true
		}

		if hasAdvanced && m.AtDestination() { // Reached destination
			logger.Debug("Entity reached destination")
			e.Pos.X = m.Destination().X
			e.Pos.Y = m.Destination().Y
			if err := l.entityStore.UpdateEntity(e); err != nil {
				return fmt.Errorf("failed to update entity %d: %w", e.ID, err)
			}
		}

		// check if the entity has a path component
		if c, ok := l.entityStore.GetComponent(e.ID, "Path"); ok {
			p, ok := c.(*path.Path)
			if !ok {
				return fmt.Errorf("component %T is not a Path component", c)
			}

			// If the entity has no destination, set it from the path component
			if !m.HasDestination() {
				steps := calculateSteps(e.Pos, p.CurrentCell(), defaultStepSize)
				logger.Debug("Entity has no path destination, setting new destination", "currentPos", e.Pos, "destination", p.CurrentCell(), "steps", steps)
				m.SetDestination(p.CurrentCell(), steps) // Set new destination with calculated steps
			} else {
				// If the entity is at its destination and the path has more cells, move to the next cell
				if m.AtDestination() {
					if p.AtDestination() {
						cells := p.Cells()
						logger.Debug("Path completed, resetting path", "entityID", e.ID, "cells", cells)
						p.Reset()            // Reset path if no more cells
						p.AddCells(cells...) // Re-add cells to path
						steps := calculateSteps(e.Pos, p.CurrentCell(), defaultStepSize)
						m.SetDestination(p.CurrentCell(), steps) // Set new destination with calculated steps
					} else if p.NextCell() {
						steps := calculateSteps(e.Pos, p.CurrentCell(), defaultStepSize)
						logger.Debug("Entity at destination, moving to next path cell", "entityID", e.ID, "nextCell", p.CurrentCell(), "steps", steps)
						m.SetDestination(p.CurrentCell(), steps) // Set new destination with calculated steps
					}
				}
			}
		}
	}
	return nil
}
