package entity

import (
	"fmt"

	"github.com/dwethmar/apostle/point"
)

// Component represents a game entity component.
type Component struct {
	EntityID int    // ID of the entity this component belongs to
	Type     string // Type of the component, e.g., "Position", "Health", etc.
	Data     any    // Data associated with the component, can be any type
}

// Entity represents a game entity with a position and a set of components.
type Entity struct {
	ID         int
	Pos        point.P // Position of the entity
	components map[string]*Component
}

type Store struct {
	entities map[int]*Entity
}

func NewStore() *Store {
	return &Store{
		entities: make(map[int]*Entity),
	}
}

func (s *Store) CreateEntity(x, y int) Entity {
	var id int
	for {
		if _, exists := s.entities[id]; !exists {
			break
		}
		id++
	}
	entity := &Entity{
		ID:         id,
		Pos:        point.P{X: x, Y: y},
		components: make(map[string]*Component),
	}
	s.entities[id] = entity
	return *entity
}

func (s *Store) Entity(id int) Entity {
	return *s.entities[id]
}

func (s *Store) RemoveEntity(id int) {
	delete(s.entities, id)
}

func (s *Store) UpdateEntity(e Entity) error {
	if _, exists := s.entities[e.ID]; !exists {
		return fmt.Errorf("entity with ID %d does not exist", e.ID)
	}
	s.entities[e.ID] = &e
	return nil
}

func (s *Store) Entities() []*Entity {
	entities := make([]*Entity, 0, len(s.entities))
	for _, entity := range s.entities {
		entities = append(entities, entity)
	}
	return entities
}

func (s *Store) AddComponent(entityID int, c Component) error {
	if entity, exists := s.entities[entityID]; exists {
		c.EntityID = entityID
		entity.components[c.Type] = &c
		return nil
	}
	return fmt.Errorf("entity with ID %d does not exist", entityID)
}

func (s *Store) GetComponent(entityID int, componentType string) (Component, bool) {
	if entity, exists := s.entities[entityID]; exists {
		if comp, exists := entity.components[componentType]; exists {
			return *comp, true
		}
	}
	return Component{}, false
}

func (s *Store) RemoveComponent(entityID int, componentType string) error {
	if entity, exists := s.entities[entityID]; exists {
		delete(entity.components, componentType)
		return nil
	}
	return fmt.Errorf("entity with ID %d does not exist", entityID)
}

func (s *Store) Components(componentType string) []*Component {
	var components []*Component
	for _, entity := range s.entities {
		if comp, exists := entity.components[componentType]; exists {
			components = append(components, comp)
		}
	}
	return components
}
