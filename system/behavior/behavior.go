package behavior

import (
	"fmt"
	"log/slog"

	"github.com/dwethmar/apostle/component"
	"github.com/dwethmar/apostle/component/agent"
	"github.com/dwethmar/apostle/component/factory"
	"github.com/dwethmar/apostle/component/kind"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/entity/blueprint"
	"github.com/dwethmar/apostle/event"
	"github.com/dwethmar/apostle/input"
	"github.com/dwethmar/apostle/point"
	"github.com/dwethmar/apostle/system/world"
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
	eventBus         *event.Bus

	// events
	subscriptions []int
	click         *point.P
}

func New(logger *slog.Logger, componentFactory *factory.Factory, entityStore *entity.Store, componentStore *component.Store, pathfinder PathFinder, eventBus *event.Bus) *Behavior {
	b := &Behavior{
		logger:           logger.With(slog.String("system", "behavior")),
		componentFactory: componentFactory,
		entityStore:      entityStore,
		componentStore:   componentStore,
		pathfinder:       pathfinder,
		eventBus:         eventBus,
	}
	b.subscriptions = []int{
		b.eventBus.Subscribe(event.MatcherFunc(func(e event.Event) bool {
			_, ok := e.(*input.Click)
			return ok
		}), func(e event.Event) error {
			clickEvent := e.(*input.Click)
			p := point.P{X: clickEvent.X, Y: clickEvent.Y}
			b.click = &p
			return nil
		}),
	}
	return b
}

func (b *Behavior) Update() error {
	var newTargetEntity *entity.Entity
	if b.click != nil {
		defer func() { b.click = nil }()
		p := point.P{
			X: b.click.X / world.CellSize,
			Y: b.click.Y / world.CellSize,
		}
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

	dest, hasDest := p.Destination()
	if !hasDest || !dest.Neighboring(targetEntity.Pos()) {
		b.logger.Debug("Agent's target has moved, recalculating path", "entityID", a.EntityID(), "oldTargetPos", dest, "newTargetPos", targetEntity.Pos())
		p.Clear()
		// Also clear the movement destination to stop the entity from continuing on the old path
		m := e.Components().Movement()
		// Immediately recalculate the path from current position
		steps := b.pathfinder.Find(m.Destination(), targetEntity.Pos())
		if len(steps) == 0 {
			b.logger.Warn("No path found to target after recalculation", "targetID", targetEntity.ID(), "entityID", a.EntityID())
			a.Reset()
			return nil // No path found, reset the agent
		}
		b.logger.Debug("Agent recalculated path to target", "targetID", targetEntity.ID(), "steps", steps)
		// Set the path component with the calculated steps
		p.Reset()
		p.AddCells(steps[:len(steps)-1]...)
	}
	return nil
}
