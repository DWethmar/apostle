package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/dwethmar/apostle/component/movement"
	"github.com/dwethmar/apostle/component/path"
	"github.com/dwethmar/apostle/drawer"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/locomotion"
	"github.com/dwethmar/apostle/point"
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

	entityStore := entity.NewStore()
	entityStore.CreateEntity(1, 1)

	e := entityStore.CreateEntity(10, 10)

	m := movement.NewComponent(e.ID)
	entityStore.AddComponent(m)

	p := path.NewComponent(e.ID)
	p.AddCells(point.P{X: 12, Y: 11}, point.P{X: 13, Y: 11}, point.P{X: 14, Y: 11}, point.P{X: 14, Y: 12}, point.P{X: 14, Y: 13})
	entityStore.AddComponent(p)

	game := &Game{
		drawers: []Drawer{
			drawer.New(tr, entityStore),
		},
		systems: []System{
			locomotion.New(logger, entityStore),
		},
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Apostle")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
