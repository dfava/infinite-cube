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
				next := cur.ApplyMove(mv)
				if _, seen := g.Nodes[next]; seen {
					// We still need to record the edge if it's new
					tr := model.Transition{From: cur, Move: mv, To: next}
					g.Edges[cur] = append(g.Edges[cur], tr)
					continue
				}
				if tryApply(top, cur, mv, v, g) {
					g.Nodes[next] = struct{}{}
					queue = append(queue, next)
				}
			}
		}

		// Check if we need to explore simultaneous moves.
		// For InfiniteCube8, we know 2-hinge moves are sufficient.
		// We avoid 3-hinge moves to keep the graph clean and tests fast.
		for i := 0; i < len(top.Hinges); i++ {
			for j := i + 1; j < len(top.Hinges); j++ {
				h1, h2 := top.Hinges[i], top.Hinges[j]
				nextPoses := []model.HingePose{model.Pose0, model.Pose180}

				for _, np1 := range nextPoses {
					for _, np2 := range nextPoses {
						if cur.Pose(h1.ID) == np1 && cur.Pose(h2.ID) == np2 {
							continue
						}

						mv := model.Move{Changes: []model.HingeChange{
							{Hinge: h1.ID, To: np1},
							{Hinge: h2.ID, To: np2},
						}}
						if cur.Pose(h1.ID) != np1 && cur.Pose(h2.ID) != np2 {
							next := cur.ApplyMove(mv)
							if _, seen := g.Nodes[next]; seen {
								tr := model.Transition{From: cur, Move: mv, To: next}
								g.Edges[cur] = append(g.Edges[cur], tr)
								continue
							}

							if tryApply(top, cur, mv, v, g) {
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
