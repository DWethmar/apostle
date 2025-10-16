package movement

import (
	"github.com/dwethmar/apostle/point"
)

const Type = "Movement"

type Movement struct {
	entityID       int
	origin         point.P // Origin point for the movement
	hasDestination bool    // Indicates if a destination is set
	dest           point.P // Destination point for the movement
	steps          int     // Number of steps to reach the destination
	currentStep    int     // Current step in the movement

	// event handling
	movedEvent func(MovedEvent) // Callback for moved events
}

func NewComponent(entityID int) *Movement {
	return &Movement{
		entityID:       entityID,
		hasDestination: false,
	}
}

func (m *Movement) EntityID() int         { return m.entityID }
func (m *Movement) ComponentType() string { return Type }

func (m *Movement) Origin() point.P      { return m.origin }
func (m *Movement) Destination() point.P { return m.dest }
func (m *Movement) HasDestination() bool { return m.hasDestination }
func (m *Movement) Steps() int           { return m.steps }
func (m *Movement) CurrentStep() int     { return m.currentStep }

func (m *Movement) SetDestination(dest point.P, steps int) {
	m.origin = m.dest
	m.dest = dest
	m.steps = steps
	m.currentStep = 0
	m.hasDestination = true
}

func (m *Movement) ClearDestination() {
	m.hasDestination = false
	m.currentStep = 0
	m.steps = 0
}

func (m *Movement) AtDestination() bool {
	return m.HasDestination() && m.currentStep >= m.steps
}

func (m *Movement) AdvanceStep() {
	if m.AtDestination() {
		return
	}
	m.currentStep++
	if m.currentStep > m.steps {
		m.currentStep = m.steps // Ensure we don't exceed steps
	}
}
