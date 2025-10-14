package movement

import (
	"bytes"
	"encoding/gob"

	"github.com/dwethmar/apostle/point"
)

const Type = "Movement"

type Movement struct {
	entityID       int
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

func (m *Movement) Destination() point.P { return m.dest }
func (m *Movement) HasDestination() bool { return m.hasDestination }
func (m *Movement) Steps() int           { return m.steps }
func (m *Movement) CurrentStep() int     { return m.currentStep }

func (m *Movement) SetDestination(dest point.P, steps int) {
	m.dest = dest
	m.steps = steps
	m.currentStep = 0
	m.hasDestination = true
}

func (m *Movement) AtDestination() bool {
	return m.currentStep >= m.steps
}

func (m *Movement) AdvanceStep() {
	if m.AtDestination() {
		return // Already at destination
	}
	m.currentStep++
	if m.currentStep > m.steps {
		m.currentStep = m.steps // Ensure we don't exceed steps
	}
}

// For later use with gob encoding/decoding
func (m *Movement) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(m.Destination())
	enc.Encode(m.Steps())
	enc.Encode(m.CurrentStep())
	return buf.Bytes(), nil
}

func (m *Movement) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	dec.Decode(&m.dest)
	dec.Decode(&m.steps)
	dec.Decode(&m.currentStep)
	return nil
}
