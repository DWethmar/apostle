package generate

import (
	"math/rand"
	"time"

	"github.com/dwethmar/apostle/terrain"
)

// Generate fills the terrain with rooms + connecting corridors and then
// carves a maze in the remaining solid areas. Uses terrain.Solid to mark
// walls/solid blocks and t.Fill to carve passages (0).
func Generate(t *terrain.Terrain) error {
	rand.Seed(time.Now().UnixNano())

	w := t.Width()
	h := t.Height()

	// safety: if terrain too small just clear it
	if w < 3 || h < 3 {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if err := t.Fill(x, y, 0); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Fill everything with walls
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if err := t.Fill(x, y, terrain.Solid); err != nil {
				return err
			}
		}
	}

	// Room placement parameters
	const (
		minRoomSize  = 3
		maxRoomSize  = 9
		roomAttempts = 200
		roomMargin   = 1 // keep one tile margin around rooms
	)

	type room struct{ x, y, rw, rh, cx, cy int }
	var rooms []room

	// helper to make odd numbers (rooms sizes and positions on odd grid)
	odd := func(v int) int {
		if v%2 == 0 {
			return v - 1
		}
		return v
	}

	// try to place rooms
	for i := 0; i < roomAttempts; i++ {
		rw := minRoomSize + rand.Intn(maxRoomSize-minRoomSize+1)
		rh := minRoomSize + rand.Intn(maxRoomSize-minRoomSize+1)
		// ensure odd sizes for nicer maze connectivity
		if rw%2 == 0 {
			rw++
		}
		if rh%2 == 0 {
			rh++
		}
		// pick position leaving a margin
		maxX := w - rw - roomMargin - 1
		maxY := h - rh - roomMargin - 1
		if maxX <= roomMargin || maxY <= roomMargin {
			continue
		}
		rx := rand.Intn(maxX-roomMargin) + roomMargin
		ry := rand.Intn(maxY-roomMargin) + roomMargin
		// snap to odd coordinates
		rx = odd(rx)
		ry = odd(ry)
		if rx < roomMargin {
			rx = roomMargin
		}
		if ry < roomMargin {
			ry = roomMargin
		}
		// ensure within bounds
		if rx+rw >= w-roomMargin {
			rx = w - roomMargin - rw - 1
			rx = odd(rx)
		}
		if ry+rh >= h-roomMargin {
			ry = h - roomMargin - rh - 1
			ry = odd(ry)
		}

		// check overlap (with a 1-cell buffer)
		overlap := false
		for yy := ry - 1; yy <= ry+rh; yy++ {
			for xx := rx - 1; xx <= rx+rw; xx++ {
				if xx >= 0 && xx < w && yy >= 0 && yy < h {
					if !t.Solid(xx, yy) {
						overlap = true
						break
					}
				}
			}
			if overlap {
				break
			}
		}
		if overlap {
			continue
		}

		// carve room interior
		for yy := ry; yy < ry+rh; yy++ {
			for xx := rx; xx < rx+rw; xx++ {
				if err := t.Fill(xx, yy, 0); err != nil {
					return err
				}
			}
		}
		cx := rx + rw/2
		cy := ry + rh/2
		rooms = append(rooms, room{x: rx, y: ry, rw: rw, rh: rh, cx: cx, cy: cy})
	}

	// connect rooms with simple straight corridors (L-shaped)
	carveHoriz := func(x1, x2, y int) error {
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		for x := x1; x <= x2; x++ {
			if err := t.Fill(x, y, 0); err != nil {
				return err
			}
		}
		return nil
	}
	carveVert := func(y1, y2, x int) error {
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for y := y1; y <= y2; y++ {
			if err := t.Fill(x, y, 0); err != nil {
				return err
			}
		}
		return nil
	}

	for i := 1; i < len(rooms); i++ {
		a := rooms[i-1]
		b := rooms[i]
		// random order for L-shape
		if rand.Intn(2) == 0 {
			if err := carveHoriz(a.cx, b.cx, a.cy); err != nil {
				return err
			}
			if err := carveVert(a.cy, b.cy, b.cx); err != nil {
				return err
			}
		} else {
			if err := carveVert(a.cy, b.cy, a.cx); err != nil {
				return err
			}
			if err := carveHoriz(a.cx, b.cx, b.cy); err != nil {
				return err
			}
		}
	}

	// Now carve a maze in the remaining solid areas using recursive backtracker
	type p struct{ x, y int }
	// find a starting solid odd cell
	var start p
	found := false
	for yy := 1; yy < h-1 && !found; yy += 2 {
		for xx := 1; xx < w-1; xx += 2 {
			if t.Solid(xx, yy) {
				start = p{xx, yy}
				found = true
				break
			}
		}
	}
	if !found {
		// nothing left to maze
		return nil
	}

	// carve start cell
	if err := t.Fill(start.x, start.y, 0); err != nil {
		return err
	}

	var stack []p
	stack = append(stack, start)

	dirs := []struct{ dx, dy int }{
		{0, -2}, {0, 2}, {2, 0}, {-2, 0},
	}

	for len(stack) > 0 {
		cur := stack[len(stack)-1]

		// gather solid neighbors 2 cells away
		neighbors := make([]p, 0, 4)
		for _, d := range dirs {
			nx := cur.x + d.dx
			ny := cur.y + d.dy
			if nx >= 1 && nx < w-1 && ny >= 1 && ny < h-1 && t.Solid(nx, ny) {
				neighbors = append(neighbors, p{nx, ny})
			}
		}

		if len(neighbors) == 0 {
			// backtrack
			stack = stack[:len(stack)-1]
			continue
		}

		nb := neighbors[rand.Intn(len(neighbors))]
		mx := (cur.x + nb.x) / 2
		my := (cur.y + nb.y) / 2

		if err := t.Fill(mx, my, 0); err != nil {
			return err
		}
		if err := t.Fill(nb.x, nb.y, 0); err != nil {
			return err
		}

		stack = append(stack, nb)
	}

	// Build a snapshot of solidity and then set border flags based on adjacent solids.
	// This ensures rooms/corridors (non-solid) have border flags where they touch walls,
	// and wall cells can also carry border flags towards adjacent floors.
	solid := make([][]bool, h)
	for y := 0; y < h; y++ {
		solid[y] = make([]bool, w)
		for x := 0; x < w; x++ {
			solid[y][x] = t.Solid(x, y)
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var cell terrain.Cell = 0
			if solid[y][x] {
				cell |= terrain.Solid
			}
			// mark borders where this cell meets a solid/non-solid neighbor
			if y-1 >= 0 && solid[y-1][x] != solid[y][x] {
				cell |= terrain.BorderNorth
			}
			if y+1 < h && solid[y+1][x] != solid[y][x] {
				cell |= terrain.BorderSouth
			}
			if x-1 >= 0 && solid[y][x-1] != solid[y][x] {
				cell |= terrain.BorderWest
			}
			if x+1 < w && solid[y][x+1] != solid[y][x] {
				cell |= terrain.BorderEast
			}
			if err := t.Fill(x, y, cell); err != nil {
				return err
			}
		}
	}

	// Add some extra decorative walls / border flags inside open areas.
	// These are low-probability so they don't destroy maze connectivity.
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			// only consider non-solid floor cells
			if t.Solid(x, y) {
				continue
			}
			// small chance to add a decorative border on a floor cell
			if rand.Float32() < 0.02 {
				// pick one random border to add
				var pick terrain.Cell
				switch rand.Intn(4) {
				case 0:
					pick = terrain.BorderNorth
				case 1:
					pick = terrain.BorderSouth
				case 2:
					pick = terrain.BorderWest
				default:
					pick = terrain.BorderEast
				}
				// reconstruct current flags and add the pick
				var cell terrain.Cell = 0
				if t.Solid(x, y) {
					cell |= terrain.Solid
				}
				if t.HasFlag(x, y, terrain.BorderNorth) {
					cell |= terrain.BorderNorth
				}
				if t.HasFlag(x, y, terrain.BorderSouth) {
					cell |= terrain.BorderSouth
				}
				if t.HasFlag(x, y, terrain.BorderWest) {
					cell |= terrain.BorderWest
				}
				if t.HasFlag(x, y, terrain.BorderEast) {
					cell |= terrain.BorderEast
				}
				cell |= pick
				if err := t.Fill(x, y, cell); err != nil {
					return err
				}
			}
			// very small chance to place a thin solid wall (obstacle)
			if rand.Float32() < 0.005 {
				if err := t.Fill(x, y, terrain.Solid); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
