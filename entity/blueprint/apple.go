package blueprint

import (
	"github.com/dwethmar/apostle/component/factory"
	"github.com/dwethmar/apostle/component/kind"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/point"
)

func NewApple(p point.P, s *entity.Store, componentFactory *factory.Factory) (*entity.Entity, error) {
	e := s.CreateEntity(p)
	k := componentFactory.NewKindComponent(e.ID())
	k.SetValue(kind.Apple)
	if err := e.Components().SetKind(k); err != nil {
		return nil, err
	}
	return e, nil
}
