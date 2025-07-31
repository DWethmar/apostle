package path

import (
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/point"
)

type Path struct {
	Cells   []point.P
	Current int // Index of the current cell in the path
}

func NewComponent() *entity.Component {
	return &entity.Component{
		Type: "Path",
		Data: &Path{
			Cells:   make([]point.P, 0),
			Current: 0,
		},
	}
}
