package blueprint

import (
	"errors"

	"github.com/dwethmar/apostle/component/factory"
	"github.com/dwethmar/apostle/component/kind"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/point"
)

func NewHuman(p point.P, s *entity.Store, componentFactory *factory.Factory) (*entity.Entity, error) {
	e := s.CreateEntity(p)
	k := componentFactory.NewKindComponent(e.ID())
	k.SetValue(kind.Human)
	return e, errors.Join(
		e.Components().SetKind(k),
		e.Components().SetMovement(componentFactory.NewMovementComponent(e.ID())),
		e.Components().SetPath(componentFactory.NewPathComponent(e.ID())),
		e.Components().SetAgent(componentFactory.NewAgentComponent(e.ID())),
	)
}
