package agent

type TargetEntityAcquiredEvent struct {
	EntityID int // ID of the entity that acquired the target
}

func (e *TargetEntityAcquiredEvent) Event() string { return "TargetEntityAcquired" }
