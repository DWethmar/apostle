package behavior

import (
	"fmt"
	"log/slog"

	"github.com/dwethmar/apostle/component"
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
	componentStore   *component.Store
	pathfinder       PathFinder // Interface for pathfinding algorithms
}

func New(logger *slog.Logger, componentFactory *factory.Factory, entityStore *entity.Store, componentStore *component.Store, pathfinder PathFinder) *Behavior {
	return &Behavior{
		logger:           logger,
		componentFactory: componentFactory,
		entityStore:      entityStore,
		componentStore:   componentStore,
		pathfinder:       pathfinder,
	}
}

func (b *Behavior) Update() error {
	for _, a := range b.componentStore.AgentEntries() {
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
	if !a.HasTargetEntity() {
		return
	}
	if _, ok := b.entityStore.Entity(a.TargetEntityID()); !ok {
		b.logger.Info("Agent's target entity has been removed, resetting target", "entityID", a.EntityID(), "removedTargetID", a.TargetEntityID())
		a.Reset()
	}
}

func (b *Behavior) lookForTargets(a *agent.Agent) error {
	if !a.HasTargetEntity() {
		for _, k := range b.componentStore.KindEntries() {
			if e, ok := b.entityStore.Entity(k.EntityID()); ok {
				a.SetTargetEntity(e.ID())
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
	if e.Components().Path() != nil {
		p = e.Components().Path()
	} else {
		p = b.componentFactory.NewPathComponent(a.EntityID())
		if err := e.Components().SetPath(p); err != nil {
			return fmt.Errorf("failed to add path component to entity %d: %w", a.EntityID(), err)
		}
	}

	targetEntity, ok := b.entityStore.Entity(a.TargetEntityID())
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
