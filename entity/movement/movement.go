package movement

import (
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/point"
)

const Type = "Movement"

type Movement struct {
	Dest        point.P // Destination point for the movement
	Steps       int     // Number of steps to reach the destination
	CurrentStep int     // Current step in the movement
}

func NewComponent() *entity.Component {
	return &entity.Component{
		Type: Type,
		Data: &Movement{},
	}
}
