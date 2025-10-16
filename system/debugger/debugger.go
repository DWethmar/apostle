package debugger

import (
	"fmt"
	"image"
	"log/slog"
	"slices"

	"github.com/dwethmar/apostle/component"
	"github.com/dwethmar/apostle/component/agent"
	"github.com/dwethmar/apostle/component/kind"
	"github.com/dwethmar/apostle/component/movement"
	"github.com/dwethmar/apostle/component/path"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/propagation"
	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Debugger struct {
	logger         *slog.Logger
	debugui        debugui.DebugUI
	count          int
	entityStore    *entity.Store
	componentStore *component.Store
	windowBounds   image.Rectangle
	pointerPressed bool // whether the pointer is currently pressed within the debugger UI we dont want to propagate events outside the debugger UI
}

func New(logger *slog.Logger, entityStore *entity.Store, componentStore *component.Store) *Debugger {
	return &Debugger{
		logger:         logger.With(slog.String("system", "debugger")),
		entityStore:    entityStore,
		componentStore: componentStore,
	}
}

func (d *Debugger) Update() error {
	entities := d.entityStore.Entities()
	slices.SortFunc(entities, func(a, b *entity.Entity) int {
		if a.ID() < b.ID() {
			return -1
		}
		if a.ID() > b.ID() {
			return 1
		}
		return 0
	})
	if _, err := d.debugui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Test", image.Rect(50, 50, 500, 700), func(layout debugui.ContainerLayout) {
			d.windowBounds = layout.Bounds
			ctx.TreeNode("entities", func() {
				ctx.Loop(len(entities), func(i int) {
					entity := entities[i]
					ctx.TreeNode(fmt.Sprintf("entity: %d", entity.ID()), func() {
						ctx.Button("destroy").On(func() {
							d.entityStore.RemoveEntity(entity.ID())
						})
						ctx.Text(fmt.Sprintf("Pos: %d, %d", entity.Pos().X, entity.Pos().Y))
						// components
						if agemt := entity.Components().Agent(); agemt != nil {
							ctx.TreeNode(agemt.ComponentType(), func() {
								d.DebugAgentComponent(ctx, agemt)
							})
						}
						if kind := entity.Components().Kind(); kind != nil {
							ctx.TreeNode(kind.ComponentType(), func() {
								d.DebugKindComponent(ctx, kind)
							})
						}
						if path := entity.Components().Path(); path != nil {
							ctx.TreeNode(path.ComponentType(), func() {
								d.DebugPathComponent(ctx, path)
							})
						}
						if movement := entity.Components().Movement(); movement != nil {
							ctx.TreeNode(movement.ComponentType(), func() {
								d.DebugMovementComponent(ctx, movement)
							})
						}
					})
				})
			})
		})
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (d *Debugger) DebugAgentComponent(ctx *debugui.Context, a *agent.Agent) {
	ctx.Text(fmt.Sprintf("goal: %s", a.Goal()))
	if a.HasTargetEntity() {
		ctx.Text(fmt.Sprintf("target entity ID: %d", a.TargetEntityID()))
	} else {
		ctx.Text("no target entity")
	}
	ctx.Button("reset").On(func() {
		a.Reset()
	})
	ctx.Button("clear target").On(func() {
		a.SetTargetEntity(agent.NoTargetID)
	})
}

func (d *Debugger) DebugKindComponent(ctx *debugui.Context, k *kind.Kind) {
	ctx.Text(fmt.Sprintf("value: %s", k.Value()))
}

func (d *Debugger) DebugPathComponent(ctx *debugui.Context, p *path.Path) {
	cells := p.Cells()
	ctx.Loop(len(cells), func(i int) {
		cell := cells[i]
		if p.CurrentCell().Equal(cell) {
			ctx.Text(fmt.Sprintf("cell %d: %d, %d (current)", i, cell.X, cell.Y))
		} else {
			ctx.Text(fmt.Sprintf("cell %d: %d, %d", i, cell.X, cell.Y))
		}
	})
}

func (d *Debugger) DebugMovementComponent(ctx *debugui.Context, m *movement.Movement) {
	ctx.Text(fmt.Sprintf("has destination: %t", m.HasDestination()))
	ctx.Text(fmt.Sprintf("origin: %d, %d", m.Origin().X, m.Origin().Y))
	ctx.Text(fmt.Sprintf("destination: %d, %d", m.Destination().X, m.Destination().Y))
	ctx.Text(fmt.Sprintf("at destination: %t", m.AtDestination()))
	ctx.Text(fmt.Sprintf("steps: %d/%d", m.CurrentStep(), m.Steps()))
}

func (d *Debugger) Draw(screen *ebiten.Image) {
	d.debugui.Draw(screen)
}

func (d *Debugger) OnPointerPressed(x, y int) propagation.Event {
	if d.windowBounds.Overlaps(image.Rect(x, y, x+1, y+1)) {
		d.pointerPressed = true
		return propagation.Stop
	}
	if d.pointerPressed {
		return propagation.Stop
	}
	return propagation.Propagate
}

func (d *Debugger) OnPointerReleased(x, y int) propagation.Event {
	if d.pointerPressed {
		d.pointerPressed = false
		return propagation.Stop
	}
	return propagation.Propagate
}
