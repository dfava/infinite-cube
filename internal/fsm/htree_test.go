package fsm

import (
	"testing"

	"infinite-cube/internal/model"
	"infinite-cube/internal/topology"
	"infinite-cube/internal/validate"
)

func TestHTreeTransitions(t *testing.T) {
	top := topology.HTree(2)
	start := model.State{}
	v := validate.PermissiveValidator{}

	g := Enumerate(top, start, v, 2)

	// Check if any transition from start involves both hinge 1 and 4.
	found14 := false
	for _, tr := range g.Edges[start] {
		if len(tr.Move.Changes) == 2 {
			h1 := tr.Move.Changes[0].Hinge
			h2 := tr.Move.Changes[1].Hinge
			if (h1 == 1 && h2 == 4) || (h1 == 4 && h2 == 1) {
				found14 = true
				break
			}
		}
	}

	if found14 {
		t.Errorf("found unwanted simultaneous transition for hinges 1 and 4 in HTree")
	}
}
