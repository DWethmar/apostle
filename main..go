package main

import (
	"fmt"
	"log"

	"github.com/dwethmar/apostle/drawer"
	"github.com/dwethmar/apostle/entity"
	"github.com/dwethmar/apostle/entity/movement"
	"github.com/dwethmar/apostle/locomotion"
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
	tr := terrain.New()
	for s := range tr.Walk() {
		if s.X%2 == 0 && s.Y%2 == 0 {
			// Fill every second cell with solid terrain
			if err := tr.Fill(s.X, s.Y, terrain.Solid); err != nil {
				log.Fatalf("failed to fill cell (%d, %d): %v", s.X, s.Y, err)
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
	m := movement.NewComponent()

	m.Data.(*movement.Movement).Dest.X = 11
	m.Data.(*movement.Movement).Dest.Y = 11
	m.Data.(*movement.Movement).Steps = 50
	m.Data.(*movement.Movement).CurrentStep = 0

	entityStore.AddComponent(e.ID, *m)

	game := &Game{
		drawers: []Drawer{
			drawer.New(tr, entityStore),
		},
		systems: []System{
			locomotion.New(entityStore),
		},
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Grid Renderer")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
