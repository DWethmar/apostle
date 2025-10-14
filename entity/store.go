package entity

import (
	"github.com/dwethmar/apostle/component"
	"github.com/dwethmar/apostle/point"
)

type Store struct {
	entities       map[int]*Entity
	componentStore *component.Store
}

func NewStore(componentCollection *component.Store) *Store {
	return &Store{
		entities:       make(map[int]*Entity),
		componentStore: componentCollection,
	}
}

func (s *Store) CreateEntity(pos point.P) *Entity {
	var id int
	for {
		if _, exists := s.entities[id]; !exists {
			break
		}
		id++
	}
	entity := &Entity{
		id:         id,
		pos:        pos,
		components: component.NewComponents(s.componentStore),
	}
	s.entities[id] = entity
	return entity
}

func (s *Store) Entity(id int) (*Entity, bool) {
	entity, exists := s.entities[id]
	return entity, exists
}

func (s *Store) RemoveEntity(id int) {
	delete(s.entities, id)
}

func (s *Store) Entities() []*Entity {
	entities := make([]*Entity, 0, len(s.entities))
	for _, entity := range s.entities {
		entities = append(entities, entity)
	}
	return entities
}
