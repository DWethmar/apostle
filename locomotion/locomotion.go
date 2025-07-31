package locomotion

import (
	"fmt"

	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/entity/movement"
)

type Locomotion struct {
	entityStore *entity.Store
}

func New(entityStore *entity.Store) *Locomotion {
	return &Locomotion{
		entityStore: entityStore,
	}
}

func (l *Locomotion) Update() error {
	for _, c := range l.entityStore.Components("Movement") {
		e := l.entityStore.Entity(c.EntityID)
		m := c.Data.(*movement.Movement)
		if e.Pos.Equal(m.Dest) {
			continue // Already at destination
		}
		if m.CurrentStep < m.Steps {
			m.CurrentStep++
		}
		if m.CurrentStep == m.Steps {
			e.Pos.X = m.Dest.X
			e.Pos.Y = m.Dest.Y
			if err := l.entityStore.UpdateEntity(e); err != nil {
				return fmt.Errorf("failed to update entity %d: %w", e.ID, err)
			}
			m.CurrentStep = 0
			m.Steps = 0
			if err := l.entityStore.UpdateEntity(e); err != nil {
				return fmt.Errorf("failed to reset movement for entity %d: %w", e.ID, err)
			}
		}
	}
	return nil
}
