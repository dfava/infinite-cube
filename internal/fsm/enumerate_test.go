package fsm

import (
	"testing"

	"infinite-cube/internal/model"
	"infinite-cube/internal/topology"
	"infinite-cube/internal/validate"
)

func TestIntuition(t *testing.T) {
	top := topology.TwoCubeHinge()
	start := model.State{}
	validator := validate.PermissiveValidator{}
	g := Enumerate(top, start, validator, 2)
	if len(g.Nodes) == 0 {
		t.Fatalf("expected at least one reachable node")
	}
}

func TestEnumerateHypercubeCountWithPermissiveValidator(t *testing.T) {
	top := topology.InfiniteCube8()
	start := model.State{}
	v := validate.PermissiveValidator{}

	g := Enumerate(top, start, v, 2)

	// wantNodes := wantNodes * len(top.Hinges)
	// Each hinge can be in 3 poses: 0, 90, 180.
	// With PermissiveValidator, all combinations are reachable since we can go 0 <-> 90 <-> 180.
	wantNodes := 1
	for range len(top.Hinges) {
		wantNodes *= 3
	}
	if len(g.Nodes) != wantNodes {
		t.Fatalf("expected %d nodes, got %d", wantNodes, len(g.Nodes))
	}

	// Each node has moves to adjacent poses.
	// If a hinge is at 0 or 180, it has 1 adjacent pose (90).
	// If a hinge is at 90, it has 2 adjacent poses (0 and 180).
	// Total single edges = Sum over all states of (number of hinges at 90 * 2 + number of hinges at 0/180 * 1)

	// For one hinge:
	// States: 0, 90, 180
	// Edges: 0->90, 180->90, 90->0, 90->180 (Total 4)

	// For N hinges:
	// Number of single-move edges = N * 4 * 3^(N-1)
	n := len(top.Hinges)
	wantSingleEdges := n * 4
	for i := 0; i < n-1; i++ {
		wantSingleEdges *= 3
	}

	// For simultaneous moves of 2 hinges:
	// A pair of hinges (h1, h2) has moves to (adj1 x adj2).
	// If h1 is at 0, adj1={90} (size 1)
	// If h1 is at 90, adj1={0, 180} (size 2)
	// If h1 is at 180, adj1={90} (size 1)
	// Total edges for one pair:
	// (h1,h2) state: (0,0)->(90,90) [1 edge]
	// (0,90)->(90,0), (90,180) [2 edges]
	// (0,180)->(90,90) [1 edge]
	// (90,0)->(0,90), (180,90) [2 edges]
	// (90,90)->(0,0), (0,180), (180,0), (180,180) [4 edges]
	// (90,180)->(0,90), (180,90) [2 edges]
	// (180,0)->(90,90) [1 edge]
	// (180,90)->(90,0), (90,180) [2 edges]
	// (180,180)->(90,90) [1 edge]
	// Total edges for one pair across all 9 states = 1+2+1+2+4+2+1+2+1 = 16

	// With the new heuristic, independent simultaneous moves are only recorded if no
	// proper subset is valid. With PermissiveValidator, all single-hinge moves are valid.
	// Therefore, NO simultaneous moves will be recorded by Enumerate because they are
	// all considered "compound moves".
	wantEdges := wantSingleEdges
	gotEdges := 0
	for _, out := range g.Edges {
		gotEdges += len(out)
	}
	if gotEdges != wantEdges {
		t.Fatalf("expected %d edges (single-hinge moves only), got %d", wantEdges, gotEdges)
	}
}

type BlockedValidator struct {
	validate.PermissiveValidator
}

func (BlockedValidator) ValidState(top model.Topology, s model.State) bool {
	// Only Pose0 and Pose180 for both hinges simultaneously is valid.
	// Individual Pose180 for either hinge is "blocked" (e.g. by collision).
	h0 := s.Pose(0)
	h1 := s.Pose(1)
	if h0 != h1 {
		return false
	}
	return true
}

func TestEnumerateSimultaneousMoves(t *testing.T) {
	// Two cubes with two hinges (not physically possible but good for testing logic)
	// Or better: Use a custom validator that requires simultaneous moves.
	top := model.Topology{
		Cubes: []model.CubeID{0, 1},
		Hinges: []model.Hinge{
			{ID: 0, A: 0, B: 1},
			{ID: 1, A: 0, B: 1},
		},
	}
	start := model.State{}
	v := BlockedValidator{}

	g := Enumerate(top, start, v, 2)

	// In this test, single moves from Pose0,0 lead to Pose180,0 or Pose0,180.
	// Both are invalid according to BlockedValidator.
	// But a simultaneous move to Pose180,180 is valid.

	found := false
	for node := range g.Nodes {
		if node.Pose(0) == model.Pose180 && node.Pose(1) == model.Pose180 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("simultaneous move to Pose180,180 not found")
	}
}

func TestEnumerateInvalidStart(t *testing.T) {
	top := topology.TwoCubeHinge()
	v := &validate.StructuralValidator{}
	// Invalid state: out of range bits
	start := model.State{PoseBits: 0xFFFFFFFF}
	g := Enumerate(top, start, v, 2)
	if len(g.Nodes) != 0 {
		t.Fatalf("expected 0 nodes for invalid start state")
	}
}

func TestTwoCubeHingeThroughReachability(t *testing.T) {
	// TwoCubeHingeThrough() is designed so that Pose0 and Pose180 are collision-free,
	// but Pose90 (the mandatory intermediate state) has a collision.
	// Therefore, Pose180 should be unreachable from Pose0.
	top := topology.TwoCubeHingeThrough()
	start := model.State{} // Pose0 for all hinges
	v := &validate.StructuralValidator{}

	g := Enumerate(top, start, v, 2)

	// Only the initial Pose0 state should be reachable.
	if len(g.Nodes) != 1 {
		t.Errorf("expected 1 reachable node, got %d", len(g.Nodes))
	}

	for node := range g.Nodes {
		if node.Pose(0) != model.Pose0 {
			t.Errorf("expected only Pose0 to be reachable, but found %v", node.Pose(0))
		}
	}
}
