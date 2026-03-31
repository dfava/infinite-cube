package fsm

import (
	"infinite-cube/internal/model"
	"infinite-cube/internal/validate"
)

// Enumerate explores the reachable FSM state space.
// It primarily tries single-hinge flips. If the user topology implies that
// some hinges must move together, it might be necessary to try combinations.
func Enumerate(top model.Topology, start model.State, v validate.Validator) *Graph {
	g := NewGraph()
	if !v.ValidState(top, start) {
		return g
	}

	queue := []model.State{start}
	g.Nodes[start] = struct{}{}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		// Single-hinge moves
		for _, h := range top.Hinges {
			nextPoses := []model.HingePose{model.Pose0, model.Pose180}
			for _, np := range nextPoses {
				if cur.Pose(h.ID) == np {
					continue
				}
				mv := model.Move{Changes: []model.HingeChange{{Hinge: h.ID, To: np}}}
				if tryApply(top, cur, mv, v, g) {
					next := cur.ApplyMove(mv)
					if _, seen := g.Nodes[next]; !seen {
						g.Nodes[next] = struct{}{}
						queue = append(queue, next)
					}
				}
			}
		}

		// Check if we need to explore simultaneous moves.
		// A simple heuristic: if no single-hinge moves were possible from this state,
		// or if we're in a specific topology that we know requires it.
		// For now, let's only try multi-hinge moves if NO single-hinge moves were found from this node.
		// However, that might be too restrictive.
		// Let's stick to the issue description's requirement: support them.
		// To keep the graph clean, we'll only add multi-hinge transitions if they lead to NEW states
		// OR if they are specifically required (but we don't have a "required" flag yet).

		// Actually, the issue is that it's NOT possible to make a move one hinge at a time.
		// So tryApply for single hinges will return false.
		// Let's just always try pairs if single moves didn't work, OR just always try pairs but that's expensive.

		// Let's refine: Always try single moves. Then try pairs only for those that didn't work as singles?
		// No, let's just always try pairs but maybe only for InfiniteCube? No, keep it general.

		// To fix the test regression, we should only try pairs if they are "necessary" or if we want full exploration.
		// For now, let's try pairs only if they lead to a valid state that wasn't reachable via a SINGLE valid move.

		for i := 0; i < len(top.Hinges); i++ {
			for j := i + 1; j < len(top.Hinges); j++ {
				h1, h2 := top.Hinges[i], top.Hinges[j]
				nextPoses := []model.HingePose{model.Pose0, model.Pose180}

				for _, np1 := range nextPoses {
					for _, np2 := range nextPoses {
						if cur.Pose(h1.ID) == np1 && cur.Pose(h2.ID) == np2 {
							continue
						}

						// If this transition (or its components) is already possible via single moves,
						// we might skip it to keep the graph simple, but a simultaneous move is distinct.
						// The issue says "it is not possible to make a move one hinge at a time".
						// This implies that tryApply(h1) and tryApply(h2) would both be false.

						mv1 := model.Move{Changes: []model.HingeChange{{Hinge: h1.ID, To: np1}}}
						mv2 := model.Move{Changes: []model.HingeChange{{Hinge: h2.ID, To: np2}}}

						v1 := v.ValidState(top, cur.ApplyMove(mv1)) && v.ValidTransition(top, cur, mv1, cur.ApplyMove(mv1))
						v2 := v.ValidState(top, cur.ApplyMove(mv2)) && v.ValidTransition(top, cur, mv2, cur.ApplyMove(mv2))

						if v1 && v2 {
							// Both can move independently. While they COULD move together,
							// it's not strictly necessary to explore it as a single transition
							// unless we want to model all possible simultaneous actions.
							continue
						}

						mv := model.Move{Changes: []model.HingeChange{
							{Hinge: h1.ID, To: np1},
							{Hinge: h2.ID, To: np2},
						}}
						if tryApply(top, cur, mv, v, g) {
							next := cur.ApplyMove(mv)
							if _, seen := g.Nodes[next]; !seen {
								g.Nodes[next] = struct{}{}
								queue = append(queue, next)
							}
						}
					}
				}
			}
		}
	}

	return g
}

func tryApply(top model.Topology, cur model.State, mv model.Move, v validate.Validator, g *Graph) bool {
	next := cur.ApplyMove(mv)
	if !v.ValidState(top, next) || !v.ValidTransition(top, cur, mv, next) {
		return false
	}

	tr := model.Transition{From: cur, Move: mv, To: next}
	g.Edges[cur] = append(g.Edges[cur], tr)
	return true
}
