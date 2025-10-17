package world

import (
	"github.com/dwethmar/apostle/point"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func centerCellX(x int) float32 {
	return float32(x)*CellSize + CellSize/2
}

func centerCellY(y int) float32 {
	return float32(y)*CellSize + CellSize/2
}

func drawEntityDiamond(screen *ebiten.Image, x, y float32) {
	width := float32(CellSize) * 0.57 // slimmer than full cell width
	height := float32(CellSize) * 0.9 // taller diamond but bottom tip at cell base

	var path vector.Path
	// Bottom tip at the real entity position
	path.MoveTo(x+CellSize/2, y+CellSize)
	// Right tip
	path.LineTo(x+CellSize/2+width/2, y+CellSize-height/2)
	// Top tip
	path.LineTo(x+CellSize/2, y+CellSize-height)
	// Left tip
	path.LineTo(x+CellSize/2-width/2, y+CellSize-height/2)
	path.Close()

	dopt := &vector.DrawPathOptions{}
	dopt.AntiAlias = true
	dopt.ColorScale.ScaleWithColor(colorEntity)

	vector.FillPath(screen, &path, &vector.FillOptions{
		FillRule: vector.FillRuleEvenOdd,
	}, dopt)
}

func drawPath(screen *ebiten.Image, points []point.P) {
	for i := 0; i < len(points)-1; i++ {
		vector.StrokeLine(screen, centerCellX(points[i].X), centerCellY(points[i].Y),
			centerCellX(points[i+1].X), centerCellY(points[i+1].Y), 2, colorPath, false)
	}
}

func drawApple(screen *ebiten.Image, x, y float32) {
	x += CellSize / 2
	y += CellSize / 2
	vector.FillCircle(screen, x, y, float32(CellSize)*0.4, colorApple, true)
}
