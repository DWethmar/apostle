package behavior

import (
	"fmt"
	"log/slog"

	"github.com/dwethmar/apostle/component"
	"github.com/dwethmar/apostle/component/agent"
	"github.com/dwethmar/apostle/component/factory"
	"github.com/dwethmar/apostle/component/kind"
	"github.com/dwethmar/apostle/drawer"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/entity/blueprint"
	"github.com/dwethmar/apostle/point"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	var newTargetEntity *entity.Entity
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mX, mY := ebiten.CursorPosition()
		p := point.P{X: mX / drawer.CellSize, Y: mY / drawer.CellSize}
		e, err := blueprint.NewApple(p, b.entityStore, b.componentFactory)
		if err != nil {
			return fmt.Errorf("failed to create apple entity at %v: %w", p, err)
		}
		newTargetEntity = e
		// delete all other apples
		for _, k := range b.componentStore.KindEntries() {
			if k.Value() == kind.Apple && k.EntityID() != e.ID() {
				b.logger.Info("Removing old apple", "entityID", k.EntityID())
				b.entityStore.RemoveEntity(k.EntityID())
			}
		}
	}

	for _, a := range b.componentStore.AgentEntries() {
		if newTargetEntity != nil {
			a.SetTargetEntity(newTargetEntity.ID())
			a.SetGoal(agent.MoveAdjacentToTarget)
			continue
		}
		b.clearTargetIfEntityRemoved(a)
		switch a.Goal() {
		case agent.None:
			if err := b.lookForTargets(a); err != nil {
				return fmt.Errorf("failed to look for targets for agent %d: %w", a.EntityID(), err)
			}
		case agent.MoveAdjacentToTarget:
			if err := b.moveToTarget(a); err != nil {
				return fmt.Errorf("failed to move agent %d to target: %w", a.EntityID(), err)
			}
		}
	}
	return nil
}

// clearTargetIfEntityRemoved checks if the agent's target is removed and clears it if so.
func (b *Behavior) clearTargetIfEntityRemoved(a *agent.Agent) {
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
			if k.EntityID() == a.EntityID() { // don't target self
				continue
			}
			if e, ok := b.entityStore.Entity(k.EntityID()); ok {
				a.SetTargetEntity(e.ID())
				a.SetGoal(agent.MoveAdjacentToTarget)
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
	p := e.Components().Path()
	if p == nil {
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

	if dest, hasDest := p.Destination(); hasDest {
		// check if the target position is still adjacent to destination
		if !dest.Neighboring(targetEntity.Pos()) {
			b.logger.Debug("Agent's target has moved, recalculating path", "entityID", a.EntityID(), "oldTargetPos", dest, "newTargetPos", targetEntity.Pos())
			p.ClearAfterNext()
			steps := b.pathfinder.Find(e.Pos(), targetEntity.Pos())
			if len(steps) == 0 {
				b.logger.Warn("No path found to target", "targetID", targetEntity, "entityID", a.EntityID())
				a.Reset()
				return nil // No path found, reset the agent
			}
			b.logger.Debug("Agent found new path to target", "targetID", targetEntity, "steps", steps)
			// Append new steps to the existing path, excluding the last step to avoid duplication
			p.AddCells(steps[:len(steps)-1]...)
		}
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
