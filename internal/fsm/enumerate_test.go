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
