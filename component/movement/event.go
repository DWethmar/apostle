package movement

import "github.com/dwethmar/apostle/point"

type MovedEvent struct {
	EntityID int     // ID of the entity that moved
	From     point.P // Previous position of the entity
	To       point.P // New position of the entity
}
