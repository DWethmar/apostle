package blueprint

import (
	"github.com/dwethmar/apostle/component/factory"
	"github.com/dwethmar/apostle/component/kind"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/point"
)

func NewHuman(p point.P, s *entity.Store, componentFactory *factory.Factory) *entity.Entity {
	e := s.CreateEntity(p)
	k := componentFactory.NewKindComponent(e.ID())
	k.SetValue(kind.Human)
	s.AddComponent(k)
	s.AddComponent(componentFactory.NewMovementComponent(e.ID()))
	s.AddComponent(componentFactory.NewPathComponent(e.ID()))
	s.AddComponent(componentFactory.NewAgentComponent(e.ID()))
	return e
}
