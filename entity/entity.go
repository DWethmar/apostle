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
	ID         int
	Pos        point.P // Position of the entity
	components map[string]Component
}
