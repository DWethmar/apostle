package main

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand/v2"

	"github.com/dwethmar/apostle/component"
	"github.com/dwethmar/apostle/component/factory"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/entity/blueprint"
	"github.com/dwethmar/apostle/event"
	"github.com/dwethmar/apostle/pathfinding/astar"
	"github.com/dwethmar/apostle/point"
	"github.com/dwethmar/apostle/propagation"
	"github.com/dwethmar/apostle/system/behavior"
	"github.com/dwethmar/apostle/system/debugger"
	"github.com/dwethmar/apostle/system/locomotion"
	"github.com/dwethmar/apostle/system/world"
	"github.com/dwethmar/apostle/terrain"
	"github.com/dwethmar/apostle/terrain/generate"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:generate go run ./genout -config=components.yaml

type Drawer interface {
	Draw(screen *ebiten.Image)
}

type System interface {
	Update() error
}

type InputListener interface {
	OnPointerPressed(x, y int) propagation.Event
	OnPointerReleased(x, y int) propagation.Event
}

type Game struct {
	drawers        []Drawer
	systems        []System
	inputListeners []InputListener
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for _, d := range g.inputListeners {
			if d.OnPointerPressed(x, y) == propagation.Stop {
				break
			}
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for _, d := range g.inputListeners {
			if d.OnPointerReleased(x, y) == propagation.Stop {
				break
			}
		}
	}

	for _, s := range g.systems {
		if err := s.Update(); err != nil {
			return fmt.Errorf("failed to update system %T: %w", s, err)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, d := range g.drawers {
		d.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	logger := slog.New(slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	tr := terrain.New()
	// for s := range tr.Walk() {
	// 	if s.X%2 == 0 && s.Y%2 == 0 {
	// 		// Fill every second cell with solid terrain
	// 		// except the starting area around (10,10)
	// 		if (s.X >= 8 && s.X <= 12) && (s.Y >= 8 && s.Y <= 12) {
	// 			continue
	// 		}
	// 		if err := tr.Fill(s.X, s.Y, terrain.Solid); err != nil {
	// 			logger.Error("failed to fill cell", "x", s.X, "y", s.Y, "error", err)
	// 		}
	// 	}

	// 	// Add borders to some cells
	// 	if s.X == 5 && s.Y == 5 {
	// 		if err := tr.Fill(s.X, s.Y, terrain.BorderNorth|terrain.BorderWest); err != nil {
	// 			log.Fatalf("failed to fill borders at (%d, %d): %v", s.X, s.Y, err)
	// 		}
	// 	}
	// }
	generate.Generate(tr)

	componentCollection := component.NewStore()
	entityStore := entity.NewStore(componentCollection)
	eventBus := event.NewBus(0)
	componentFactory := factory.NewFactory(eventBus)

	{
		var x, y int
		for range 1000 {
			x = rand.IntN(tr.Width() - 1)
			y = rand.IntN(tr.Height() - 1)
			if !tr.Solid(x, y) {
				break
			}
		}
		blueprint.NewHuman(world.CellToCenterPX(point.P{
			X: x,
			Y: y,
		}), entityStore, componentFactory)
	}
	{
		var x, y int
		for range 1000 {
			x = rand.IntN(tr.Width() - 1)
			y = rand.IntN(tr.Height() - 1)
			if !tr.Solid(x, y) {
				break
			}
		}
		blueprint.NewApple(world.CellToCenterPX(point.P{
			X: x,
			Y: y,
		}), entityStore, componentFactory)
	}

	debugger := debugger.New(logger, entityStore, componentCollection)
	w := world.New(logger, tr, entityStore, componentCollection, eventBus)
	l := locomotion.New(logger, entityStore, componentCollection)
	b := behavior.New(logger, tr, componentFactory, entityStore, componentCollection, astar.New(tr), eventBus)

	game := &Game{
		drawers: []Drawer{
			w,
			debugger,
		},
		systems: []System{
			l,
			b,
			debugger,
		},
		inputListeners: []InputListener{
			debugger,
			w,
		},
	}

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Apostle")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
