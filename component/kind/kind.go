package kind

const Type = "Kind"

type Value uint

const (
	None Value = iota
	Human
	Apple
)

type Kind struct {
	entityID      int
	componentType string
	value         Value
}

func NewComponent(entityID int) *Kind {
	return &Kind{
		entityID:      entityID,
		componentType: Type,
		value:         0,
	}
}

func (k *Kind) EntityID() int {
	return k.entityID
}

func (k *Kind) ComponentType() string {
	return k.componentType
}

func (k *Kind) SetValue(value Value) {
	k.value = value
}

func (k *Kind) Value() Value {
	return k.value
}
