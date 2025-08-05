package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/dwethmar/apostle/behavior"
	"github.com/dwethmar/apostle/component/factory"
	"github.com/dwethmar/apostle/drawer"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/event"
	"github.com/dwethmar/apostle/locomotion"
	"github.com/dwethmar/apostle/pathfinding/astar"
	"github.com/dwethmar/apostle/terrain"
	"github.com/hajimehoshi/ebiten/v2"
)

type Drawer interface {
	Draw(screen *ebiten.Image)
}

type System interface {
	Update() error
}

type Game struct {
	drawers []Drawer
	systems []System
}

func (g *Game) Update() error {
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
	return 320, 240
}

func main() {
	logger := slog.New(slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	tr := terrain.New()
	for s := range tr.Walk() {
		if s.X%2 == 0 && s.Y%2 == 0 {
			// Fill every second cell with solid terrain
			if err := tr.Fill(s.X, s.Y, terrain.Solid); err != nil {
				logger.Error("failed to fill cell", "x", s.X, "y", s.Y, "error", err)
			}
		}

		// Add borders to some cells
		if s.X == 5 && s.Y == 5 {
			if err := tr.Fill(s.X, s.Y, terrain.BorderNorth|terrain.BorderWest); err != nil {
				log.Fatalf("failed to fill borders at (%d, %d): %v", s.X, s.Y, err)
			}
		}
	}

	eventBus := event.NewBus(0)
	componentFactory := factory.NewFactory(eventBus)

	entityStore := entity.NewStore()
	entityStore.CreateEntity(1, 1)

	e := entityStore.CreateEntity(11, 11)
	entityStore.AddComponent(componentFactory.NewMovementComponent(e.ID()))
	entityStore.AddComponent(componentFactory.NewPathComponent(e.ID()))
	entityStore.AddComponent(componentFactory.NewAgentComponent(e.ID()))

	game := &Game{
		drawers: []Drawer{
			drawer.New(tr, entityStore),
		},
		systems: []System{
			locomotion.New(logger, entityStore),
			behavior.New(logger, componentFactory, entityStore, astar.New(tr)),
		},
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Apostle")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
