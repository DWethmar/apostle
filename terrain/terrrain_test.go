package terrain_test

import (
	"testing"

	"github.com/dwethmar/apostle/point"
	"github.com/dwethmar/apostle/terrain"
)

func TestTraversable_NoObstacles(t *testing.T) {
	tr := terrain.New()
	if !tr.Traversable(point.New(10, 10), terrain.North) {
		t.Errorf("expected traversable north")
	}
	if !tr.Traversable(point.New(10, 10), terrain.South) {
		t.Errorf("expected traversable south")
	}
	if !tr.Traversable(point.New(10, 10), terrain.East) {
		t.Errorf("expected traversable east")
	}
	if !tr.Traversable(point.New(10, 10), terrain.West) {
		t.Errorf("expected traversable west")
	}
}

func TestTraversable_SolidBlocksMovement(t *testing.T) {
	tr := terrain.New()
	if err := tr.Fill(10, 9, terrain.Solid); err != nil { // solid block north of (10,10)
		t.Fatalf("failed to fill solid cell: %v", err)
	}
	if tr.Traversable(point.New(10, 10), terrain.North) {
		t.Errorf("expected blocked north due to solid cell")
	}
	if !tr.Traversable(point.New(10, 10), terrain.South) {
		t.Errorf("expected traversable south")
	}
}

func TestTraversable_BorderBlocksMovement(t *testing.T) {
	tr := terrain.New()

	// Border on current cell
	if err := tr.Fill(10, 10, terrain.BorderNorth); err != nil {
		t.Fatalf("failed to fill border cell %d %d: %v", 10, 10, err)
	}
	if tr.Traversable(point.New(10, 10), terrain.North) {
		t.Errorf("expected blocked north due to border on current cell")
	}

	// Border on target cell
	if err := tr.Fill(10, 10, 0); err != nil {
		t.Fatalf("failed to fill empty cell: %v", err)
	}
	if err := tr.Fill(10, 9, terrain.BorderSouth); err != nil {
		t.Fatalf("failed to fill border cell %d %d: %v", 10, 9, err)
	}
	if tr.Traversable(point.New(10, 10), terrain.North) {
		t.Errorf("expected blocked north due to border on target cell")
	}
}

func TestTraversable_BordersDontBlockOtherDirections(t *testing.T) {
	tr := terrain.New()
	_ = tr.Fill(10, 10, terrain.BorderNorth)

	if !tr.Traversable(point.New(10, 10), terrain.South) {
		t.Errorf("expected traversable south")
	}
	if !tr.Traversable(point.New(10, 10), terrain.East) {
		t.Errorf("expected traversable east")
	}
	if !tr.Traversable(point.New(10, 10), terrain.West) {
		t.Errorf("expected traversable west")
	}
}

func TestTraversable_OutOfBounds(t *testing.T) {
	tr := terrain.New()

	if tr.Traversable(point.New(0, 0), terrain.West) {
		t.Errorf("expected blocked west out of bounds")
	}
	if tr.Traversable(point.New(0, 0), terrain.North) {
		t.Errorf("expected blocked north out of bounds")
	}
	if tr.Traversable(point.New(19, 19), terrain.South) {
		t.Errorf("expected blocked south out of bounds")
	}
	if tr.Traversable(point.New(19, 19), terrain.East) {
		t.Errorf("expected blocked east out of bounds")
	}
}

func TestBorders_ReturnsSetBorders(t *testing.T) {
	tr := terrain.New()
	_ = tr.Fill(5, 5, terrain.BorderNorth|terrain.BorderWest)

	borders := tr.Walls(5, 5)
	if len(borders) != 2 {
		t.Errorf("expected 2 borders, got %d", len(borders))
	}

	foundNorth := false
	foundWest := false
	for _, b := range borders {
		if b == terrain.BorderNorth {
			foundNorth = true
		}
		if b == terrain.BorderWest {
			foundWest = true
		}
	}

	if !foundNorth || !foundWest {
		t.Errorf("expected both north and west borders")
	}
}

func TestSolid_ChecksCorrectly(t *testing.T) {
	tr := terrain.New()
	if err := tr.Fill(3, 3, terrain.Solid); err != nil {
		t.Fatalf("failed to fill solid cell: %v", err)
	}

	if !tr.Solid(3, 3) {
		t.Errorf("expected solid cell at (3,3)")
	}
	if tr.Solid(2, 2) {
		t.Errorf("expected non-solid cell at (2,2)")
	}
}

func Test_Walk(t *testing.T) {
	tr := terrain.New()

	for x := range 20 {
		for y := range 20 {
			if err := tr.Fill(x, y, terrain.Solid|terrain.BorderNorth); err != nil {
				t.Fatalf("failed to fill cell (%d, %d): %v", x, y, err)
			}
		}
	}

	for s := range tr.Walk() {
		// Check that all cells are solid and have a border
		if s.Cell&terrain.Solid == 0 {
			t.Errorf("expected solid cell at (%d, %d)", s.X, s.Y)
		}
		if s.Cell&terrain.BorderNorth == 0 {
			t.Errorf("expected border north at (%d, %d)", s.X, s.Y)
		}
	}
}
