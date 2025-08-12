package kind

import "github.com/dwethmar/apostle/component"

const Type = "Kind"

type Value uint

const (
	None Value = iota
	Human
	Apple
)

type Kind struct {
	*component.Component
	value Value
}

func NewComponent(entityID int) *Kind {
	return &Kind{
		Component: component.NewComponent(entityID, Type),
		value:     0,
	}
}

func (k *Kind) SetValue(value Value) {
	k.value = value
}

func (k *Kind) Value() Value {
	return k.value
}
