package topology

import (
	"infinite-cube/internal/model"
	"testing"
)

func TestHTreeParametric(t *testing.T) {
	tests := []struct {
		levels     int
		wantCubes  int
		wantHinges int
	}{
		{levels: 1, wantCubes: 2, wantHinges: 1},
		{levels: 2, wantCubes: 8, wantHinges: 7},
		{levels: 3, wantCubes: 22, wantHinges: 21},
	}

	for _, tt := range tests {
		top := HTree(tt.levels)
		if len(top.Cubes) != tt.wantCubes {
			t.Errorf("HTree(%d) got %d cubes, want %d", tt.levels, len(top.Cubes), tt.wantCubes)
		}
		if len(top.Hinges) != tt.wantHinges {
			t.Errorf("HTree(%d) got %d hinges, want %d", tt.levels, len(top.Hinges), tt.wantHinges)
		}

		// Verify State support for this many hinges
		s := model.State{}
		for i := 0; i < len(top.Hinges); i++ {
			hID := model.HingeID(i)
			// Try to set and get Pose180 for each hinge to test bit shifts
			move := model.Move{Changes: []model.HingeChange{{Hinge: hID, To: model.Pose180}}}
			s = s.ApplyMove(move)
			if s.Pose(hID) != model.Pose180 {
				t.Errorf("HTree(%d) hinge %d: failed to set/get Pose180", tt.levels, i)
			}
		}

		// Verify cube IDs are sequential and match the count
		for i := 0; i < len(top.Cubes); i++ {
			found := false
			for _, id := range top.Cubes {
				if id == model.CubeID(i) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("HTree(%d) missing cube ID %d", tt.levels, i)
			}
		}

		// Verify hinge IDs are sequential
		for i := 0; i < len(top.Hinges); i++ {
			found := false
			for _, h := range top.Hinges {
				if h.ID == model.HingeID(i) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("HTree(%d) missing hinge ID %d", tt.levels, i)
			}
		}
	}
}
