package behavior

import (
	"fmt"
	"log/slog"

	"github.com/dwethmar/apostle/component/agent"
	"github.com/dwethmar/apostle/component/factory"
	"github.com/dwethmar/apostle/component/path"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/point"
)

// PathFinder defines the behavior for pathfinding algorithms.
type PathFinder interface {
	Find(start, end point.P) []point.P
}

type Behavior struct {
	logger           *slog.Logger
	componentFactory *factory.Factory
	entityStore      *entity.Store
	pathfinder       PathFinder // Interface for pathfinding algorithms
}

func New(logger *slog.Logger, componentFactory *factory.Factory, entityStore *entity.Store, pathfinder PathFinder) *Behavior {
	return &Behavior{
		logger:           logger,
		componentFactory: componentFactory,
		entityStore:      entityStore,
		pathfinder:       pathfinder,
	}
}

func (b *Behavior) Update() error {
	for _, c := range b.entityStore.Components(agent.Type) {
		a, ok := c.(*agent.Agent)
		if !ok {
			return fmt.Errorf("component %T is not an Agent component", c)
		}
		b.logger.Debug("Updating agent behavior", "entityID", a.EntityID(), "goal", a.Goal())

		switch a.Goal() {
		case agent.None:
			if len(a.TargetEntityIDs()) == 0 {
				// If the agent is idle and has no targets, it can look for targets
				targets := b.entityStore.Entities() // Get all entities as potential targets
				for _, target := range targets {
					if target.ID() != a.EntityID() { // Avoid targeting itself
						a.AddTargetEntityID(target.ID())
						a.SetGoal(agent.MoveToTarget) // Change state to moving towards the target
						b.logger.Debug("Agent found target", "targetID", target.ID(), "entityID", a.EntityID())
						break // Stop after finding the first target
					}
				}
			}
		case agent.MoveToTarget:
			e, ok := b.entityStore.Entity(a.EntityID())
			if !ok {
				return fmt.Errorf("entity with ID %d does not exist", a.EntityID())
			}
			// check if the entity has a path
			var p *path.Path
			c, ok := b.entityStore.GetComponent(a.EntityID(), path.Type)
			if ok {
				p, ok = c.(*path.Path)
				if !ok {
					return fmt.Errorf("component %T is not a Path component", c)
				}
			} else {
				p = b.componentFactory.NewPathComponent(a.EntityID())
				if err := b.entityStore.AddComponent(p); err != nil {
					return fmt.Errorf("failed to add path component to entity %d: %w", a.EntityID(), err)
				}
			}

			targetEntities := a.TargetEntityIDs()
			if len(targetEntities) == 0 {
				b.logger.Warn("Agent has no target entities", "entityID", a.EntityID())
				a.Reset()
				continue
			}

			if _, hasDest := p.Destination(); hasDest {
				// b.logger.Debug("Agent moving towards destination", "entityID", a.EntityID(), "destination", d)
				// Here you would implement the logic to move the agent towards the destination
				// For example, using the locomotion system to advance the movement component
			} else {
				// calculate steps to the first target
				// TODO: neighbor the target entity. Do not step into the same cell.
				targetID := targetEntities[0]
				targetEntity, ok := b.entityStore.Entity(targetID)
				if !ok {
					return fmt.Errorf("target entity %d not found", targetID)
				}
				steps := b.pathfinder.Find(e.Pos(), targetEntity.Pos())
				if len(steps) == 0 {
					b.logger.Warn("No path found to target", "targetID", targetID, "entityID", a.EntityID())
					a.Reset()
					continue
				}
				b.logger.Debug("Agent found path to target", "targetID", targetID, "steps", steps)
				// Set the path component with the calculated steps
				p.Reset()
				p.AddCells(steps...)
			}
		}
	}
	return nil
}
