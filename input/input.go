package input

const ClickEvent = "click"

type Click struct {
	X int
	Y int
}

func (c Click) Event() string {
	return ClickEvent
}
