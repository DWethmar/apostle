package agent

import "github.com/dwethmar/apostle/component"

const Type = "Agent"

type Goal uint

const (
	None Goal = iota
	MoveToTarget
)

// Agent represents an entity that can perform actions and has a state.
// It can target other entities and has a state machine for its behavior.
// ideas:
// - have memory: remembered facts about the world
type Agent struct {
	*component.Component
	goal            Goal
	targetEntityIDs []int // IDs of entities that this agent can target
}

func NewAgent(entityID int) *Agent {
	return &Agent{
		Component:       component.NewComponent(entityID, Type),
		goal:            None,
		targetEntityIDs: make([]int, 0),
	}
}

func (a *Agent) SetGoal(goal Goal) {
	a.goal = goal
}

func (a *Agent) Goal() Goal {
	return a.goal
}

func (a *Agent) AddTargetEntityID(entityID int) {
	a.targetEntityIDs = append(a.targetEntityIDs, entityID)
}

func (a *Agent) TargetEntityIDs() []int {
	return a.targetEntityIDs
}

func (a *Agent) ClearTargetEntityIDs() {
	a.targetEntityIDs = make([]int, 0)
}

func (a *Agent) Reset() {
	a.goal = None
	a.ClearTargetEntityIDs()
}
