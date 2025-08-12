package agent

import (
	"github.com/dwethmar/apostle/component"
	"github.com/dwethmar/apostle/entity"
)

const Type = "Agent"

type Goal uint

const (
	None Goal = iota
	MoveAdjecentToTarget
)

const NoTargetID = -1

// Agent represents an entity that can perform actions and has a state.
// It can target other entities and has a state machine for its behavior.
// ideas:
// - have memory: remembered facts about the world
type Agent struct {
	*component.Component
	goal           Goal
	targetEntityID int // ID of the entity that this agent targets
	// systems
	entityStore *entity.Store
	// event handlers
	emitTargetEntitySetEvent func(*TargetEntityAcquiredEvent) // Event handler for when a target entity is acquired
}

type AgentOption func(*Agent)

func WithEmitTargetEntityAcquiredEvent(handler func(*TargetEntityAcquiredEvent)) AgentOption {
	return func(a *Agent) {
		a.emitTargetEntitySetEvent = handler
	}
}

func NewAgent(entityID int, entityStore *entity.Store, opts ...AgentOption) *Agent {
	a := &Agent{
		Component:      component.NewComponent(entityID, Type),
		goal:           None,
		targetEntityID: NoTargetID,
		entityStore:    entityStore,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *Agent) SetGoal(goal Goal) {
	a.goal = goal
}

func (a *Agent) Goal() Goal {
	return a.goal
}

// SetTargetEntity sets the target entity for the agent.
func (a *Agent) SetTargetEntity(e *entity.Entity) {
	if e == nil {
		a.targetEntityID = NoTargetID
	} else {
		a.targetEntityID = e.ID()
	}
	if a.emitTargetEntitySetEvent != nil {
		a.emitTargetEntitySetEvent(&TargetEntityAcquiredEvent{
			EntityID:     a.EntityID(),
			TargetEntity: e,
		})
	}
}

func (a *Agent) TargetEntity() (*entity.Entity, bool) {
	if a.targetEntityID == NoTargetID {
		return nil, false
	}
	return a.entityStore.Entity(a.targetEntityID)
}

func (a *Agent) TargetEntityID() int {
	return a.targetEntityID
}

func (a *Agent) Reset() {
	a.goal = None
	a.SetTargetEntity(nil)
}
