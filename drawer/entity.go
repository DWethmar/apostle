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
	// vector.StrokeLine(screen, centerCellX(points[0].X), centerCellY(points[0].Y),
	// 	centerCellX(points[len(points)-1].X), centerCellY(points[len(points)-1].Y), 2, colorPath, false)

	for i := 0; i < len(points)-1; i++ {
		vector.StrokeLine(screen, centerCellX(points[i].X), centerCellY(points[i].Y),
			centerCellX(points[i+1].X), centerCellY(points[i+1].Y), 2, colorPath, false)
	}
}

func drawApple(screen *ebiten.Image, x, y float32) {
	x += cellSize / 2
	y += cellSize / 2
	vector.FillCircle(screen, x, y, float32(cellSize)*0.4, colorApple, true)
}
