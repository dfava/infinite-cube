package fsm

import (
	"infinite-cube/internal/model"
	"infinite-cube/internal/validate"
)

// Enumerate explores the reachable FSM state space by single-hinge flips.
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

		for _, h := range top.Hinges {
			nextPose := model.Pose0
			if cur.Pose(h.ID) == model.Pose0 {
				nextPose = model.Pose180
			}
			mv := model.Move{Hinge: h.ID, To: nextPose}
			next := cur.ApplyMove(mv)

			if !v.ValidState(top, next) || !v.ValidTransition(top, cur, mv, next) {
				continue
			}

			tr := model.Transition{From: cur, Move: mv, To: next}
			g.Edges[cur] = append(g.Edges[cur], tr)

			if _, seen := g.Nodes[next]; seen {
				continue
			}
			g.Nodes[next] = struct{}{}
			queue = append(queue, next)
		}
	}

	return g
}
