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
	g := Enumerate(top, start, validator)
	if len(g.Nodes) == 0 {
		t.Fatalf("expected at least one reachable node")
	}
}

func TestEnumerateHypercubeCountWithPermissiveValidator(t *testing.T) {
	top := topology.InfiniteCube8()
	start := model.State{}
	v := validate.PermissiveValidator{}

	g := Enumerate(top, start, v)

	wantNodes := 1 << len(top.Hinges)
	if len(g.Nodes) != wantNodes {
		t.Fatalf("expected %d nodes, got %d", wantNodes, len(g.Nodes))
	}

	wantEdges := wantNodes * len(top.Hinges)
	gotEdges := 0
	for _, out := range g.Edges {
		gotEdges += len(out)
	}
	if gotEdges != wantEdges {
		t.Fatalf("expected %d edges, got %d", wantEdges, gotEdges)
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

	g := Enumerate(top, start, v)

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
	v := validate.StructuralValidator{}
	// Invalid state: out of range bits
	start := model.State{PoseBits: 0xFFFFFFFF}
	g := Enumerate(top, start, v)
	if len(g.Nodes) != 0 {
		t.Fatalf("expected 0 nodes for invalid start state")
	}
}
