package fsm

import (
	"infinite-cube/internal/model"
	"infinite-cube/internal/validate"
)

// Enumerate explores the reachable FSM state space.
// maxSimultaneous specifies the maximum number of hinges to check simultaneously.
func Enumerate(top model.Topology, start model.State, v validate.Validator, maxSimultaneous int) *Graph {
	g := NewGraph()
	if !v.ValidState(top, start) {
		return g
	}

	components := findComponents(top)

	queue := []model.State{start}
	g.Nodes[start] = struct{}{}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, comp := range components {
			// Explore all subsets of the component from size 1 up to maxSimultaneous
			for k := 1; k <= maxSimultaneous && k <= len(comp); k++ {
				subsets := combinations(comp, k)
				for _, subset := range subsets {
					trySubsets(top, cur, subset, v, g, &queue)
				}
			}
		}
	}

	return g
}

func trySubsets(top model.Topology, cur model.State, subset []model.HingeID, v validate.Validator, g *Graph, queue *[]model.State) {
	// For the given subset of hinges, try all possible valid adjacent pose combinations.
	// Since we only move 90 degrees at a time, each hinge has 1 or 2 possible next poses.
	var hinges []model.Hinge
	for _, id := range subset {
		hinges = append(hinges, findHinge(top, id))
	}

	var options [][]model.HingePose
	for _, h := range hinges {
		curP := cur.Pose(h.ID)
		var nextPoses []model.HingePose
		switch curP {
		case model.Pose0:
			nextPoses = []model.HingePose{model.Pose90}
		case model.Pose90:
			nextPoses = []model.HingePose{model.Pose0, model.Pose180}
		case model.Pose180:
			nextPoses = []model.HingePose{model.Pose90}
		}
		options = append(options, nextPoses)
	}

	// Generate Cartesian product of all nextPoses options
	for _, combo := range cartesianProduct(options) {
		mv := model.Move{Changes: make([]model.HingeChange, len(subset))}
		for i, id := range subset {
			mv.Changes[i] = model.HingeChange{Hinge: id, To: combo[i]}
		}

		next := cur.ApplyMove(mv)
		if _, seen := g.Nodes[next]; seen {
			if len(subset) > 1 {
				if hasValidSubset(top, cur, subset, combo, v) {
					continue
				}
			}
			// Record the edge if it's new
			tr := model.Transition{From: cur, Move: mv, To: next}
			g.Edges[cur] = append(g.Edges[cur], tr)
			continue
		}

		if len(subset) > 1 {
			// A simultaneous move is only recorded if it is "atomic".
			// We define atomic as: no proper subset of these changes is a valid move from the current state.
			if hasValidSubset(top, cur, subset, combo, v) {
				continue
			}
		}

		if tryApply(top, cur, mv, v, g) {
			g.Nodes[next] = struct{}{}
			*queue = append(*queue, next)
		}
	}
}

func hasValidSubset(top model.Topology, cur model.State, subset []model.HingeID, combo []model.HingePose, v validate.Validator) bool {
	// Check all non-empty proper subsets of the proposed simultaneous move.
	// If any subset is a valid move from the same starting state, then the full move is not atomic.
	n := len(subset)
	// We use a bitmask from 1 to 2^n - 2 to explore all non-empty proper subsets.
	for i := 1; i < (1<<n)-1; i++ {
		var subChanges []model.HingeChange
		for j := 0; j < n; j++ {
			if (i>>j)&1 == 1 {
				subChanges = append(subChanges, model.HingeChange{Hinge: subset[j], To: combo[j]})
			}
		}
		subMv := model.Move{Changes: subChanges}
		subNext := cur.ApplyMove(subMv)
		if v.ValidState(top, subNext) && v.ValidTransition(top, cur, subMv, subNext) {
			return true
		}
	}
	return false
}

func findHinge(top model.Topology, id model.HingeID) model.Hinge {
	for _, h := range top.Hinges {
		if h.ID == id {
			return h
		}
	}
	panic("hinge not found")
}

func combinations(ids []model.HingeID, k int) [][]model.HingeID {
	var result [][]model.HingeID
	var combine func(start int, current []model.HingeID)
	combine = func(start int, current []model.HingeID) {
		if len(current) == k {
			tmp := make([]model.HingeID, k)
			copy(tmp, current)
			result = append(result, tmp)
			return
		}
		for i := start; i < len(ids); i++ {
			combine(i+1, append(current, ids[i]))
		}
	}
	combine(0, nil)
	return result
}

func cartesianProduct(options [][]model.HingePose) [][]model.HingePose {
	if len(options) == 0 {
		return [][]model.HingePose{{}}
	}
	var result [][]model.HingePose
	tails := cartesianProduct(options[1:])
	for _, head := range options[0] {
		for _, tail := range tails {
			res := make([]model.HingePose, len(tail)+1)
			res[0] = head
			copy(res[1:], tail)
			result = append(result, res)
		}
	}
	return result
}

func findComponents(top model.Topology) [][]model.HingeID {
	cubeToHinges := make(map[model.CubeID][]model.HingeID)
	for _, h := range top.Hinges {
		cubeToHinges[h.A] = append(cubeToHinges[h.A], h.ID)
		cubeToHinges[h.B] = append(cubeToHinges[h.B], h.ID)
	}

	visited := make(map[model.HingeID]bool)
	var components [][]model.HingeID

	for _, h := range top.Hinges {
		if visited[h.ID] {
			continue
		}
		var component []model.HingeID
		queue := []model.HingeID{h.ID}
		visited[h.ID] = true

		for len(queue) > 0 {
			currID := queue[0]
			queue = queue[1:]
			component = append(component, currID)

			currHinge := findHinge(top, currID)
			cubes := []model.CubeID{currHinge.A, currHinge.B}

			for _, cube := range cubes {
				for _, neighborID := range cubeToHinges[cube] {
					if !visited[neighborID] {
						visited[neighborID] = true
						queue = append(queue, neighborID)
					}
				}
			}
		}
		components = append(components, component)
	}
	return components
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
