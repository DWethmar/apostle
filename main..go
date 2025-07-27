package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Point struct {
	X, Y int
}

type CellType uint8

const (
	Empty CellType = iota
	Solid
)

type Cell struct {
	Point
	Type CellType
}

type Border struct {
	Index, Low, High int
}

type Actor struct {
	PrevPos         Point
	Pos             Point
	Waypoints       []Point
	currentWaypoint int
	moveTimer       int
	moveSpeed       int     // frames between moves
	path            []Point // Current path to follow
	pathIndex       int     // Current position in path
	animTimer       int     // Animation timer
	animDuration    int     // How long animation takes
}

type Game struct {
	gridSize          int
	cellSize          float32
	gridColor         color.Color
	gridWidth         int
	gridHeight        int
	cells             [][]Cell
	horizontalBorders []Border
	verticalBorders   []Border
	actor             Actor
}

func (g *Game) canMoveTo(from, to Point) bool {
	// Check bounds
	if to.X < 0 || to.X >= g.gridWidth || to.Y < 0 || to.Y >= g.gridHeight {
		return false
	}

	// Check if target cell is solid
	if g.cells[to.Y][to.X].Type == Solid {
		return false
	}

	// Check horizontal borders (movement in X direction)
	if from.X != to.X {
		for _, border := range g.verticalBorders {
			borderX := border.Index
			// Check if moving across this border
			if (from.X < borderX && to.X >= borderX) || (from.X >= borderX && to.X < borderX) {
				// Check if current Y position is within border range
				if from.Y >= border.Low && from.Y < border.High {
					return false
				}
			}
		}
	}

	// Check vertical borders (movement in Y direction)
	if from.Y != to.Y {
		for _, border := range g.horizontalBorders {
			borderY := border.Index
			// Check if moving across this border
			if (from.Y < borderY && to.Y >= borderY) || (from.Y >= borderY && to.Y < borderY) {
				// Check if current X position is within border range
				if from.X >= border.Low && from.X < border.High {
					return false
				}
			}
		}
	}

	return true
}

func (g *Game) Update() error {
	// Update actor movement
	if len(g.actor.Waypoints) > 0 {
		// Update animation timer
		if g.actor.animTimer < g.actor.animDuration {
			g.actor.animTimer++
		}

		g.actor.moveTimer++

		if g.actor.moveTimer >= g.actor.moveSpeed {
			g.actor.moveTimer = 0

			// Get current target waypoint
			if g.actor.currentWaypoint < len(g.actor.Waypoints) {
				target := g.actor.Waypoints[g.actor.currentWaypoint]

				// If we don't have a path or reached the end of current path, find new path
				if len(g.actor.path) == 0 || g.actor.pathIndex >= len(g.actor.path) {
					g.actor.path = g.findPath(g.actor.Pos, target)
					g.actor.pathIndex = 0

					// If no path found, skip to next waypoint
					if g.actor.path == nil {
						g.actor.currentWaypoint++
						if g.actor.currentWaypoint >= len(g.actor.Waypoints) {
							g.actor.currentWaypoint = 0
						}
						return nil
					}
				}

				// Follow the path
				if g.actor.pathIndex < len(g.actor.path)-1 {
					nextPos := g.actor.path[g.actor.pathIndex+1]

					// Double-check we can still move there (in case obstacles changed)
					if g.canMoveTo(g.actor.Pos, nextPos) {
						g.actor.PrevPos = g.actor.Pos // Store previous position
						g.actor.Pos = nextPos
						g.actor.pathIndex++
						g.actor.animTimer = 0 // Reset animation timer
					} else {
						// Path is blocked, recalculate
						g.actor.path = nil
						g.actor.pathIndex = 0
					}
				} else {
					// Reached current waypoint, move to next
					g.actor.currentWaypoint++
					g.actor.path = nil // Clear path for next waypoint
					g.actor.pathIndex = 0

					// Loop back to first waypoint when done
					if g.actor.currentWaypoint >= len(g.actor.Waypoints) {
						g.actor.currentWaypoint = 0
					}
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a dark background
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw cells first (so borders appear on top)
	for y := 0; y < g.gridHeight; y++ {
		for x := 0; x < g.gridWidth; x++ {
			cell := g.cells[y][x]
			cellX := float32(x * g.gridSize)
			cellY := float32(y * g.gridSize)
			cellWidth := float32(g.gridSize)
			cellHeight := float32(g.gridSize)

			// Choose color based on cell type
			var cellColor color.Color
			switch cell.Type {
			case Empty:
				cellColor = color.RGBA{40, 40, 50, 255} // Dark gray for empty
			case Solid:
				cellColor = color.RGBA{80, 60, 40, 255} // Brown for solid
			default:
				cellColor = color.RGBA{60, 60, 60, 255} // Default gray
			}

			// Fill the cell rectangle
			vector.DrawFilledRect(screen, cellX, cellY, cellWidth, cellHeight, cellColor, false)

			// Draw coordinate text for debugging with smaller appearance
			coordText := fmt.Sprintf("%d,%d", x, y)

			// Create a temporary image for the text
			textImg := ebiten.NewImage(int(g.cellSize), int(g.cellSize)) // Small image for text
			ebitenutil.DebugPrint(textImg, coordText)

			// Draw the text image scaled down to make it appear smaller
			op := &ebiten.DrawImageOptions{
				Filter: ebiten.FilterLinear, // Use nearest neighbor for pixel art style
			}
			op.GeoM.Scale(0.5, 0.5) // Scale down to 50% size
			op.GeoM.Translate(float64(cellX+1), float64(cellY+2))
			screen.DrawImage(textImg, op)
		}
	}

	// Draw horizontal borders
	for _, border := range g.horizontalBorders {
		y := float32(border.Index * g.gridSize)
		startX := float32(border.Low * g.gridSize)
		endX := float32(border.High * g.gridSize)
		vector.StrokeLine(screen, startX, y, endX, y, 1, g.gridColor, false)
	}

	// Draw vertical borders
	for _, border := range g.verticalBorders {
		x := float32(border.Index * g.gridSize)
		startY := float32(border.Low * g.gridSize)
		endY := float32(border.High * g.gridSize)
		vector.StrokeLine(screen, x, startY, x, endY, 1, g.gridColor, false)
	}

	// Draw actor with animation between PrevPos and Pos
	var actorX, actorY float32

	// Calculate animation progress (0.0 to 1.0)
	progress := float32(g.actor.animTimer) / float32(g.actor.animDuration)
	if progress > 1.0 {
		progress = 1.0
	}

	// Smooth easing (ease-out)
	progress = 1.0 - (1.0-progress)*(1.0-progress)

	// Interpolate between previous and current position
	prevX := float32(g.actor.PrevPos.X * g.gridSize)
	prevY := float32(g.actor.PrevPos.Y * g.gridSize)
	currX := float32(g.actor.Pos.X * g.gridSize)
	currY := float32(g.actor.Pos.Y * g.gridSize)

	actorX = prevX + (currX-prevX)*progress
	actorY = prevY + (currY-prevY)*progress

	actorSize := float32(g.gridSize) * 0.6   // Make actor slightly smaller than cell
	actorOffset := float32(g.gridSize) * 0.2 // Center the actor in the cell

	vector.DrawFilledCircle(screen, actorX+actorOffset+actorSize/2, actorY+actorOffset+actorSize/2, actorSize/2, color.RGBA{255, 100, 100, 255}, false)

	// Draw waypoints
	for i, waypoint := range g.actor.Waypoints {
		wpX := float32(waypoint.X * g.gridSize)
		wpY := float32(waypoint.Y * g.gridSize)
		wpSize := float32(g.gridSize) * 0.3
		wpOffset := float32(g.gridSize) * 0.35

		// Different color for current target waypoint
		var wpColor color.Color
		if i == g.actor.currentWaypoint {
			wpColor = color.RGBA{255, 255, 100, 255} // Yellow for current target
		} else {
			wpColor = color.RGBA{100, 255, 100, 255} // Green for other waypoints
		}

		vector.DrawFilledCircle(screen, wpX+wpOffset+wpSize/2, wpY+wpOffset+wpSize/2, wpSize/2, wpColor, false)
	}

	// Draw current path for debugging
	if len(g.actor.path) > 1 {
		for i := 0; i < len(g.actor.path)-1; i++ {
			current := g.actor.path[i]
			next := g.actor.path[i+1]

			currentX := float32(current.X*g.gridSize + g.gridSize/2)
			currentY := float32(current.Y*g.gridSize + g.gridSize/2)
			nextX := float32(next.X*g.gridSize + g.gridSize/2)
			nextY := float32(next.Y*g.gridSize + g.gridSize/2)

			// Different color for completed vs remaining path
			var lineColor color.Color
			if i < g.actor.pathIndex {
				lineColor = color.RGBA{100, 100, 200, 128} // Blue for completed path
			} else {
				lineColor = color.RGBA{200, 100, 200, 128} // Purple for remaining path
			}

			vector.StrokeLine(screen, currentX, currentY, nextX, nextY, 2, lineColor, false)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	game := &Game{
		gridSize:   20, // Fixed to match the calculation
		cellSize:   20.0,
		gridColor:  color.RGBA{100, 100, 100, 255},
		gridWidth:  16, // 320/20 = 16 cells wide
		gridHeight: 12, // 240/20 = 12 cells tall
	}

	// Initialize the grid with cells
	game.cells = make([][]Cell, game.gridHeight)
	for y := 0; y < game.gridHeight; y++ {
		game.cells[y] = make([]Cell, game.gridWidth)
		for x := 0; x < game.gridWidth; x++ {
			game.cells[y][x] = Cell{
				Point: Point{X: x, Y: y},
				Type:  Empty, // Default to empty
			}
		}
	}

	// Example: Set some cells to solid
	game.cells[1][1].Type = Solid
	game.cells[1][2].Type = Solid
	game.cells[2][1].Type = Solid
	game.cells[2][2].Type = Solid

	// Create a wall
	for x := 5; x < 10; x++ {
		game.cells[4][x].Type = Solid
	}

	// Create another solid area
	for y := 7; y < 10; y++ {
		for x := 12; x < 15; x++ {
			game.cells[y][x].Type = Solid
		}
	}

	// Initialize actor at (0,0) with waypoints
	game.actor = Actor{
		PrevPos: Point{X: 0, Y: 0},
		Pos:     Point{X: 0, Y: 0},
		Waypoints: []Point{
			{X: 3, Y: 3},  // First waypoint
			{X: 8, Y: 2},  // Second waypoint
			{X: 10, Y: 8}, // Third waypoint
			{X: 5, Y: 10}, // Fourth waypoint
			{X: 0, Y: 0},  // Return to start
		},
		currentWaypoint: 0,
		moveTimer:       0,
		moveSpeed:       20, // Slightly slower to see animation better
		path:            nil,
		pathIndex:       0,
		animTimer:       0,
		animDuration:    15, // Animation takes 15 frames
	}

	// Example: Create some border segments
	// Horizontal borders (row borders)
	game.horizontalBorders = []Border{
		{Index: 0, Low: 0, High: game.gridWidth},               // Top edge of grid
		{Index: 3, Low: 2, High: 8},                            // Horizontal border from x=2 to x=8 at y=3
		{Index: 6, Low: 1, High: 5},                            // Horizontal border from x=1 to x=5 at y=6
		{Index: game.gridHeight, Low: 0, High: game.gridWidth}, // Bottom edge of grid
	}

	// Vertical borders (column borders)
	game.verticalBorders = []Border{
		{Index: 0, Low: 0, High: game.gridHeight},              // Left edge of grid
		{Index: 4, Low: 1, High: 7},                            // Vertical border from y=1 to y=7 at x=4
		{Index: 10, Low: 2, High: 9},                           // Vertical border from y=2 to y=9 at x=10
		{Index: game.gridWidth, Low: 0, High: game.gridHeight}, // Right edge of grid
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Grid Renderer")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
