package component

// Component represents a game entity component.
type Component struct {
	EID int    // ID of the entity this component belongs to
	T   string // Type of the component, e.g., "Position", "Health", etc.
}

func NewComponent(entityID int, componentType string) *Component {
	return &Component{
		EID: entityID,
		T:   componentType,
	}
}

func (c *Component) EntityID() int {
	return c.EID
}

func (c *Component) Type() string {
	return c.T
}
