package astar

import (
	"container/heap"
	"math"

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

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].fCost < pq[j].fCost }
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
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}

// AStar implements the PathFinder interface using A* search.
type AStar struct {
	terrain *terrain.Terrain
}

func New(t *terrain.Terrain) *AStar {
	return &AStar{terrain: t}
}

func moves() map[terrain.Direction]struct{ Dx, Dy int } {
	return map[terrain.Direction]struct{ Dx, Dy int }{
		terrain.North:     {0, -1},
		terrain.South:     {0, 1},
		terrain.East:      {1, 0},
		terrain.West:      {-1, 0},
		terrain.NorthEast: {1, -1},
		terrain.NorthWest: {-1, -1},
		terrain.SouthEast: {1, 1},
		terrain.SouthWest: {-1, 1},
	}
}

func (a *AStar) Find(start, end point.P) []point.P {
	startNode := &node{x: start.X, y: start.Y}
	goalX, goalY := end.X, end.Y

	openSet := &priorityQueue{}
	heap.Init(openSet)
	heap.Push(openSet, startNode)

	closedSet := make(map[[2]int]bool) // Using a map for closed set to track visited nodes

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*node)
		closedSet[[2]int{current.x, current.y}] = true

		if current.x == goalX && current.y == goalY {
			return reconstructPath(current)
		}

		for dir, move := range moves() {
			newX := current.x + move.Dx
			newY := current.y + move.Dy

			if !a.terrain.InBounds(newX, newY) || closedSet[[2]int{newX, newY}] {
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
			if diagonal { // Diagonal moves have a higher cost
				stepCost = math.Sqrt2
			}

			gCost := current.gCost + stepCost
			hCost := heuristic(newX, newY, goalX, goalY)
			neighbor := &node{
				x:      newX,
				y:      newY,
				gCost:  gCost,
				hCost:  hCost,
				fCost:  gCost + hCost,
				parent: current,
			}
			heap.Push(openSet, neighbor)
		}
	}
	return nil // no path found
}

func reconstructPath(n *node) []point.P {
	var path []point.P
	for n != nil {
		path = append([]point.P{{X: n.x, Y: n.y}}, path...)
		n = n.parent
	}
	return path
}

func isDiagonal(dir terrain.Direction) bool {
	return dir == terrain.NorthEast ||
		dir == terrain.NorthWest ||
		dir == terrain.SouthEast ||
		dir == terrain.SouthWest
}

func canMoveDiagonally(t *terrain.Terrain, x, y int, dir terrain.Direction) bool {
	switch dir {
	case terrain.NorthEast:
		return t.Traversable(point.New(x, y), terrain.North) &&
			t.Traversable(point.New(x, y), terrain.East)
	case terrain.NorthWest:
		return t.Traversable(point.New(x, y), terrain.North) &&
			t.Traversable(point.New(x, y), terrain.West)
	case terrain.SouthEast:
		return t.Traversable(point.New(x, y), terrain.South) &&
			t.Traversable(point.New(x, y), terrain.East)
	case terrain.SouthWest:
		return t.Traversable(point.New(x, y), terrain.South) &&
			t.Traversable(point.New(x, y), terrain.West)
	default:
		return true
	}
}
