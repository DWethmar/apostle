package agent

const Type = "Agent"

type Goal uint

const (
	None Goal = iota
	MoveAdjacentToTarget
)

const NoTargetID = -1

// Agent represents an entity that can perform actions and has a state.
// It can target other entities and has a state machine for its behavior.
// ideas:
// - have memory: remembered facts about the world
type Agent struct {
	entityID                 int
	goal                     Goal
	targetEntityID           int                              // ID of the entity that this agent targets
	emitTargetEntitySetEvent func(*TargetEntityAcquiredEvent) // Event handler for when a target entity is acquired
}

type AgentOption func(*Agent)

func WithEmitTargetEntityAcquiredEvent(handler func(*TargetEntityAcquiredEvent)) AgentOption {
	return func(a *Agent) {
		a.emitTargetEntitySetEvent = handler
	}
}

func NewAgent(entityID int, opts ...AgentOption) *Agent {
	a := &Agent{
		goal:           None,
		targetEntityID: NoTargetID,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *Agent) EntityID() int {
	return a.entityID
}

func (a *Agent) ComponentType() string {
	return Type
}

func (a *Agent) SetGoal(goal Goal) {
	a.goal = goal
}

func (a *Agent) Goal() Goal {
	return a.goal
}

// SetTargetEntity sets the target entity for the agent.
func (a *Agent) SetTargetEntity(entityID int) {
	if entityID == NoTargetID {
		a.targetEntityID = NoTargetID
	} else {
		a.targetEntityID = entityID
	}
	if a.emitTargetEntitySetEvent != nil {
		a.emitTargetEntitySetEvent(&TargetEntityAcquiredEvent{
			EntityID: a.EntityID(),
		})
	}
}

func (a *Agent) TargetEntityID() int {
	return a.targetEntityID
}

func (a *Agent) HasTargetEntity() bool {
	return a.targetEntityID != NoTargetID
}

func (a *Agent) Reset() {
	a.goal = None
	a.SetTargetEntity(NoTargetID)
}
