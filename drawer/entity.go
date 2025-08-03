package drawer

import (
	"github.com/dwethmar/apostle/point"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func centerCellX(x int) float32 {
	return float32(x)*cellSize + cellSize/2
}

func centerCellY(y int) float32 {
	return float32(y)*cellSize + cellSize/2
}

func drawEntityDiamond(screen *ebiten.Image, x, y float32) {
	width := float32(cellSize) * 0.57 // slimmer than full cell width
	height := float32(cellSize) * 0.9 // taller diamond but bottom tip at cell base

	var path vector.Path
	// Bottom tip at the real entity position
	path.MoveTo(x+cellSize/2, y+cellSize)
	// Right tip
	path.LineTo(x+cellSize/2+width/2, y+cellSize-height/2)
	// Top tip
	path.LineTo(x+cellSize/2, y+cellSize-height)
	// Left tip
	path.LineTo(x+cellSize/2-width/2, y+cellSize-height/2)
	path.Close()

	vector.FillPath(screen, &path, colorEntity, true, vector.FillRuleEvenOdd)
}

func drawPath(screen *ebiten.Image, points []point.P) {
	var path vector.Path
	for _, p := range points {
		path.MoveTo(centerCellX(p.X), centerCellY(p.Y))
	}
	vector.StrokePath(screen, &path, colorPath, true, &vector.StrokeOptions{
		Width: 4,
	})
}
