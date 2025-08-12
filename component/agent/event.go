package agent

import "github.com/dwethmar/apostle/entity"

type TargetEntityAcquiredEvent struct {
	EntityID     int            // ID of the entity that acquired the target
	TargetEntity *entity.Entity // ID of the target entity that was acquired
}

func (e *TargetEntityAcquiredEvent) Event() string { return "TargetEntityAcquired" }
