package astar

import (
	"container/heap"
	"math"

	"github.com/dwethmar/apostle/direction"
	"github.com/dwethmar/apostle/point"
	"github.com/dwethmar/apostle/terrain"
)

type node struct {
	x, y   int
	gCost  float64
	hCost  float64 // hCost is the heuristic cost to the goal
	fCost  float64 // fCost is the total cost (gCost + hCost)
	parent *node
	index  int
}

type priorityQueue []*node

func (pq priorityQueue) Len() int { return len(pq) }

// Less prioritizes lower fCost, and uses hCost as a tiebreaker
func (pq priorityQueue) Less(i, j int) bool {
	if pq[i].fCost == pq[j].fCost {
		return pq[i].hCost < pq[j].hCost
	}
	return pq[i].fCost < pq[j].fCost
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	n := x.(*node)
	n.index = len(*pq)
	*pq = append(*pq, n)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := old[len(old)-1]
	*pq = old[:len(old)-1]
	return n
}

func heuristic(x1, y1, x2, y2 int) float64 {
	dx := math.Abs(float64(x1 - x2))
	dy := math.Abs(float64(y1 - y2))
	const D = 1.0         // cost for orthogonal
	const D2 = math.Sqrt2 // cost for diagonal
	return D*(dx+dy) + (D2-2*D)*math.Min(dx, dy)
}

// AStar implements the PathFinder interface using A* search.
type AStar struct {
	terrain *terrain.Terrain
}

func New(t *terrain.Terrain) *AStar {
	return &AStar{terrain: t}
}

func moves() map[direction.Direction]struct{ Dx, Dy int } {
	return map[direction.Direction]struct{ Dx, Dy int }{
		direction.North:     {0, -1},
		direction.South:     {0, 1},
		direction.East:      {1, 0},
		direction.West:      {-1, 0},
		direction.NorthEast: {1, -1},
		direction.NorthWest: {-1, -1},
		direction.SouthEast: {1, 1},
		direction.SouthWest: {-1, 1},
	}
}

func (a *AStar) Find(start, end point.P) []point.P {
	startNode := &node{
		x:     start.X,
		y:     start.Y,
		gCost: 0,
	}
	goalX, goalY := end.X, end.Y

	startNode.hCost = heuristic(startNode.x, startNode.y, goalX, goalY)
	startNode.fCost = startNode.hCost

	openSet := &priorityQueue{}
	heap.Init(openSet)
	heap.Push(openSet, startNode)

	closedSet := make(map[[2]int]bool) // visited
	bestG := make(map[[2]int]float64)  // best known gCost per coord
	bestG[[2]int{startNode.x, startNode.y}] = 0.0

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*node)

		ck := [2]int{current.x, current.y}
		// skip nodes that were already closed (outdated heap entries)
		if closedSet[ck] {
			continue
		}
		closedSet[ck] = true

		if current.x == goalX && current.y == goalY {
			return reconstructPath(current)
		}

		for dir, move := range moves() {
			newX := current.x + move.Dx
			newY := current.y + move.Dy
			nk := [2]int{newX, newY}

			if !a.terrain.InBounds(newX, newY) || closedSet[nk] {
				continue
			}
			if !a.terrain.Traversable(point.New(current.x, current.y), dir) {
				continue
			}

			diagonal := isDiagonal(dir)
			if diagonal && !canMoveDiagonally(a.terrain, current.x, current.y, dir) {
				continue
			}

			stepCost := 1.0
			if diagonal {
				stepCost = math.Sqrt2
			}

			newG := current.gCost + stepCost

			// if we've seen a better or equal gCost for this cell, skip
			if prevG, ok := bestG[nk]; ok && newG >= prevG {
				continue
			}

			bestG[nk] = newG
			hCost := heuristic(newX, newY, goalX, goalY)
			neighbor := &node{
				x:      newX,
				y:      newY,
				gCost:  newG,
				hCost:  hCost,
				fCost:  newG + hCost,
				parent: current,
			}
			heap.Push(openSet, neighbor)
		}
	}
	return nil // no path
}

func reconstructPath(n *node) []point.P {
	var path []point.P
	for n != nil {
		path = append([]point.P{{X: n.x, Y: n.y}}, path...)
		n = n.parent
	}
	return path
}

func isDiagonal(dir direction.Direction) bool {
	return dir == direction.NorthEast ||
		dir == direction.NorthWest ||
		dir == direction.SouthEast ||
		dir == direction.SouthWest
}

func canMoveDiagonally(t *terrain.Terrain, x, y int, dir direction.Direction) bool {
	// Map of directions to deltas (reuse your moves() table)
	mv := moves()

	// Identify the two orthogonal directions that compose the diagonal.
	orth1, orth2 := func(d direction.Direction) (direction.Direction, direction.Direction) {
		switch d {
		case direction.NorthEast:
			return direction.North, direction.East
		case direction.NorthWest:
			return direction.North, direction.West
		case direction.SouthEast:
			return direction.South, direction.East
		case direction.SouthWest:
			return direction.South, direction.West
		default:
			return direction.None, direction.None // if you have a None; otherwise handle as needed
		}
	}(dir)

	p := point.New(x, y)

	// 1) From the current cell, both orthogonal exits must be OK
	if !t.Traversable(p, orth1) || !t.Traversable(p, orth2) {
		return false
	}

	// 2) From each orthogonal mid cell, the second leg must also be OK
	m1 := point.New(x+mv[orth1].Dx, y+mv[orth1].Dy) // after taking orth1
	m2 := point.New(x+mv[orth2].Dx, y+mv[orth2].Dy) // after taking orth2

	// Bounds checks (optional if Traversable already does it)
	if !t.InBounds(m1.X, m1.Y) || !t.InBounds(m2.X, m2.Y) {
		return false
	}

	// From mid cell reached via orth1, must be able to go orth2.
	// From mid cell reached via orth2, must be able to go orth1.
	if !t.Traversable(m1, orth2) || !t.Traversable(m2, orth1) {
		return false
	}

	return true
}
