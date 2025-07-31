package drawer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func drawEntityDiamond(screen *ebiten.Image, x, y float32) {
	width := float32(cellSize) * 0.5  // slimmer than full cell width
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

	vector.FillPath(screen, &path, colorEntity, true, vector.FillRuleNonZero)
}
