package fsm

import (
	"fmt"
	"testing"

	"infinite-cube/internal/model"
	"infinite-cube/internal/topology"
	"infinite-cube/internal/validate"
)

func TestInfiniteCube8StateCount(t *testing.T) {
	top := topology.InfiniteCube8()
	start := model.State{}
	v := validate.StructuralValidator{}

	g := Enumerate(top, start, v)

	fmt.Printf("InfiniteCube8 with StructuralValidator found %d reachable nodes\n", len(g.Nodes))

	for s := range g.Nodes {
		fmt.Printf("State: %08x\n", s.PoseBits)
		for _, tr := range g.Edges[s] {
			fmt.Printf("  -> %08x (move: %v)\n", tr.To.PoseBits, tr.Move)
		}
	}

	if len(g.Nodes) != 10 {
		t.Errorf("expected 10 nodes for InfiniteCube8, got %d", len(g.Nodes))
	}
}

func TestInfiniteCube8EnumerateMultiple(t *testing.T) {
	top := topology.InfiniteCube8()
	v := validate.StructuralValidator{}

	allStates := make(map[model.State]bool)

	// Start from all 3^8 possible states
	for i := uint32(0); i < 6561; i++ {
		bits := uint32(0)
		temp := i
		for h := uint32(0); h < 8; h++ {
			p := model.HingePose(temp % 3)
			temp /= 3
			bits |= (uint32(p) << (2 * h))
		}

		s := model.State{PoseBits: bits}
		if v.ValidState(top, s) {
			allStates[s] = true
		}
	}

	fmt.Printf("Found %d valid '3-pose' states total\n", len(allStates))

	reachableFromStart := Enumerate(top, model.State{}, v)
	fmt.Printf("Reachable from start: %d\n", len(reachableFromStart.Nodes))

	components := 0
	visitedGlobal := make(map[model.State]bool)

	// Start from 0 to ensure we find the 10-node component first if possible
	startState := model.State{PoseBits: 0}
	if allStates[startState] {
		components++
		g := Enumerate(top, startState, v)
		fmt.Printf("Component %d (starting %08x) has %d nodes\n", components, startState.PoseBits, len(g.Nodes))
		for node := range g.Nodes {
			visitedGlobal[node] = true
		}
	}

	for s := range allStates {
		if visitedGlobal[s] {
			continue
		}
		components++
		g := Enumerate(top, s, v)
		fmt.Printf("Component %d (starting %08x) has %d nodes\n", components, s.PoseBits, len(g.Nodes))
		if len(g.Nodes) <= 10 {
			fmt.Printf("  -> Component nodes:\n")
			for n := range g.Nodes {
				fmt.Printf("     State: %08x\n", n.PoseBits)
			}
		}
		for node := range g.Nodes {
			visitedGlobal[node] = true
		}
	}
	fmt.Printf("Total components found among valid states: %d\n", components)
}
