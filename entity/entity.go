package entity

import (
	"github.com/dwethmar/apostle/point"
)

// Component represents a game entity component.
type Component interface {
	EntityID() int
	Type() string
}

// Entity represents a game entity with a position and a set of components.
type Entity struct {
	id         int
	pos        point.P // Position of the entity
	components map[string]Component
}

func (e *Entity) ID() int {
	return e.id
}

func (e *Entity) Pos() point.P {
	return e.pos
}

func (e *Entity) SetPos(pos point.P) {
	e.pos = pos
}
