package behavior

import (
	"fmt"
	"log/slog"

	"github.com/dwethmar/apostle/component/agent"
	"github.com/dwethmar/apostle/component/factory"
	"github.com/dwethmar/apostle/component/kind"
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

		b.resetTargetIfRemoved(a)

		switch a.Goal() {
		case agent.None:
			if err := b.lookForTargets(a); err != nil {
				return fmt.Errorf("failed to look for targets for agent %d: %w", a.EntityID(), err)
			}
		case agent.MoveAdjecentToTarget:
			if err := b.moveToTarget(a); err != nil {
				return fmt.Errorf("failed to move agent %d to target: %w", a.EntityID(), err)
			}
		}
	}
	return nil
}

// resetTargetIfRemoved checks if the agent's target is removed and clears it if so.
func (b *Behavior) resetTargetIfRemoved(a *agent.Agent) {
	if _, ok := a.TargetEntity(); !ok {
		a.SetTargetEntity(nil)
		a.SetGoal(agent.None)
	}
}

func (b *Behavior) lookForTargets(a *agent.Agent) error {
	if _, ok := a.TargetEntity(); !ok {
		for _, c := range b.entityStore.Components(kind.Type) {
			k, ok := c.(*kind.Kind)
			if !ok || k.Value() != kind.Apple || k.EntityID() == a.TargetEntityID() {
				continue
			}
			if e, ok := b.entityStore.Entity(k.EntityID()); ok {
				a.SetTargetEntity(e)
				a.SetGoal(agent.MoveAdjecentToTarget)
				break
			}
		}
	}
	return nil
}

func (b *Behavior) moveToTarget(a *agent.Agent) error {
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

	targetEntity, ok := a.TargetEntity()
	if !ok {
		b.logger.Warn("Agent has no target entities", "entityID", a.EntityID())
		a.Reset()
		return nil // No targets to move towards
	}

	if _, hasDest := p.Destination(); hasDest {

	} else {
		// calculate steps to the first target
		steps := b.pathfinder.Find(e.Pos(), targetEntity.Pos())
		if len(steps) == 0 {
			b.logger.Warn("No path found to target", "targetID", targetEntity, "entityID", a.EntityID())
			a.Reset()
			return nil // No path found, reset the agent
		}
		b.logger.Debug("Agent found path to target", "targetID", targetEntity, "steps", steps)
		// Set the path component with the calculated steps
		p.Reset()
		p.AddCells(steps[:len(steps)-1]...)
	}
	return nil
}
