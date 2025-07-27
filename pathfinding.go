package main

// A* pathfinding node
type PathNode struct {
	pos    Point
	parent *PathNode
	g      int // Cost from start
	h      int // Heuristic cost to goal
	f      int // Total cost (g + h)
}

// Manhattan distance heuristic
func manhattanDistance(a, b Point) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// A* pathfinding algorithm
func (g *Game) findPath(start, goal Point) []Point {
	if start.X == goal.X && start.Y == goal.Y {
		return []Point{start}
	}

	openSet := []*PathNode{}
	closedSet := make(map[Point]bool)

	startNode := &PathNode{
		pos: start,
		g:   0,
		h:   manhattanDistance(start, goal),
		f:   manhattanDistance(start, goal),
	}

	openSet = append(openSet, startNode)

	// Directions: up, down, left, right
	directions := []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for len(openSet) > 0 {
		// Find node with lowest f cost
		currentIdx := 0
		for i, node := range openSet {
			if node.f < openSet[currentIdx].f {
				currentIdx = i
			}
		}

		current := openSet[currentIdx]

		// Remove current from open set
		openSet = append(openSet[:currentIdx], openSet[currentIdx+1:]...)
		closedSet[current.pos] = true

		// Check if we reached the goal
		if current.pos.X == goal.X && current.pos.Y == goal.Y {
			// Reconstruct path
			path := []Point{}
			for node := current; node != nil; node = node.parent {
				path = append([]Point{node.pos}, path...)
			}
			return path
		}

		// Check all neighbors
		for _, dir := range directions {
			neighbor := Point{
				X: current.pos.X + dir.X,
				Y: current.pos.Y + dir.Y,
			}

			// Skip if in closed set or can't move there
			if closedSet[neighbor] || !g.canMoveTo(current.pos, neighbor) {
				continue
			}

			tentativeG := current.g + 1

			// Check if this neighbor is already in open set
			var neighborNode *PathNode
			neighborInOpen := false
			for _, node := range openSet {
				if node.pos.X == neighbor.X && node.pos.Y == neighbor.Y {
					neighborNode = node
					neighborInOpen = true
					break
				}
			}

			if !neighborInOpen {
				// Add new node to open set
				neighborNode = &PathNode{
					pos:    neighbor,
					parent: current,
					g:      tentativeG,
					h:      manhattanDistance(neighbor, goal),
				}
				neighborNode.f = neighborNode.g + neighborNode.h
				openSet = append(openSet, neighborNode)
			} else if tentativeG < neighborNode.g {
				// Update existing node with better path
				neighborNode.parent = current
				neighborNode.g = tentativeG
				neighborNode.f = neighborNode.g + neighborNode.h
			}
		}
	}

	// No path found
	return nil
}
