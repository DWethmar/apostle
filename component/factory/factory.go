package factory

import (
	"github.com/dwethmar/apostle/component/agent"
	"github.com/dwethmar/apostle/component/kind"
	"github.com/dwethmar/apostle/component/movement"
	"github.com/dwethmar/apostle/component/path"
	"github.com/dwethmar/apostle/event"
)

type Factory struct {
	eventBus *event.Bus
}

func NewFactory(eventBus *event.Bus) *Factory {
	return &Factory{
		eventBus: eventBus,
	}
}

func (f *Factory) NewAgentComponent(entityID int) *agent.Agent {
	return agent.NewAgent(
		entityID,
		agent.WithEmitTargetEntityAcquiredEvent(func(e *agent.TargetEntityAcquiredEvent) {
			f.eventBus.Publish(e)
		}),
	)
}

func (f *Factory) NewMovementComponent(entityID int) *movement.Movement {
	return movement.NewComponent(entityID)
}

func (f *Factory) NewPathComponent(entityID int) *path.Path {
	return path.NewComponent(entityID)
}

func (f *Factory) NewKindComponent(entityID int) *kind.Kind {
	return kind.NewComponent(entityID)
}
