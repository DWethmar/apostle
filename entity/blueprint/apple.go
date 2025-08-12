package blueprint

import (
	"github.com/dwethmar/apostle/component/factory"
	"github.com/dwethmar/apostle/component/kind"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/point"
)

func NewApple(p point.P, s *entity.Store, componentFactory *factory.Factory) *entity.Entity {
	e := s.CreateEntity(p)
	k := componentFactory.NewKindComponent(e.ID())
	k.SetValue(kind.Apple)
	s.AddComponent(k)
	return e
}
