package kinematics

import (
	"fmt"
	"math"
	"sort"

	"infinite-cube/internal/model"
)

// DeterministicSolver computes poses by propagating hinge transforms across the topology graph.
// It is deterministic but intentionally simple; it is a kinematic scaffold, not a full rigid-body solver.
type DeterministicSolver struct {
	// ComponentSpacing separates disconnected components.
	ComponentSpacing float64
}

func NewDeterministicSolver() DeterministicSolver {
	return DeterministicSolver{
		ComponentSpacing: 3,
	}
}

func (s DeterministicSolver) Poses(top model.Topology, state model.State) (map[model.CubeID]model.Pose, error) {
	if len(top.Cubes) == 0 {
		return nil, fmt.Errorf("empty topology")
	}
	if s.ComponentSpacing <= 0 {
		return nil, fmt.Errorf("invalid ComponentSpacing %v", s.ComponentSpacing)
	}

	cubeSet := make(map[model.CubeID]struct{}, len(top.Cubes))
	for _, c := range top.Cubes {
		if _, exists := cubeSet[c]; exists {
			return nil, fmt.Errorf("duplicate cube ID %d", c)
		}
		cubeSet[c] = struct{}{}
	}

	type incident struct {
		hinge model.Hinge
		fromA bool
	}
	adj := make(map[model.CubeID][]incident, len(top.Cubes))
	for _, c := range top.Cubes {
		adj[c] = nil
	}
	for _, h := range top.Hinges {
		if _, ok := cubeSet[h.A]; !ok {
			return nil, fmt.Errorf("hinge %d references unknown cube A=%d", h.ID, h.A)
		}
		if _, ok := cubeSet[h.B]; !ok {
			return nil, fmt.Errorf("hinge %d references unknown cube B=%d", h.ID, h.B)
		}
		adj[h.A] = append(adj[h.A], incident{hinge: h, fromA: true})
		adj[h.B] = append(adj[h.B], incident{hinge: h, fromA: false})
	}

	for c := range adj {
		sort.Slice(adj[c], func(i, j int) bool {
			return adj[c][i].hinge.ID < adj[c][j].hinge.ID
		})
	}

	cubes := append([]model.CubeID(nil), top.Cubes...)
	sort.Slice(cubes, func(i, j int) bool { return cubes[i] < cubes[j] })

	poses := make(map[model.CubeID]model.Pose, len(top.Cubes))
	visited := make(map[model.CubeID]bool, len(top.Cubes))
	component := 0

	for _, root := range cubes {
		if visited[root] {
			continue
		}
		rootPose := model.Pose{
			P: model.Vec3{X: float64(component) * s.ComponentSpacing},
			Q: model.QuatIdentity(),
		}
		poses[root] = rootPose
		visited[root] = true

		queue := []model.CubeID{root}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]

			curPose := poses[cur]
			for _, in := range adj[cur] {
				var next model.CubeID
				var nextPose model.Pose
				if in.fromA {
					next = in.hinge.B
					nextPose = propagateAtoB(curPose, in.hinge, state)
				} else {
					next = in.hinge.A
					nextPose = propagateBtoA(curPose, in.hinge, state)
				}

				if !visited[next] {
					poses[next] = nextPose
					visited[next] = true
					queue = append(queue, next)
					continue
				}

				if !poses[next].AlmostEqual(nextPose, 1e-6) {
					return nil, fmt.Errorf("inconsistent kinematic cycle at cube %d via hinge %d", next, in.hinge.ID)
				}
			}
		}
		component++
	}

	return poses, nil
}

func propagateAtoB(a model.Pose, h model.Hinge, s model.State) model.Pose {
	qRel := hingeRelativeRotation(h, s)
	qB := a.Q.Mul(qRel).Normalize()
	worldAnchor := a.P.Add(a.Q.Rotate(h.AnchorA))
	return model.Pose{
		P: worldAnchor.Sub(qB.Rotate(h.AnchorB)),
		Q: qB,
	}
}

func propagateBtoA(b model.Pose, h model.Hinge, s model.State) model.Pose {
	qRel := hingeRelativeRotation(h, s)
	qA := b.Q.Mul(qRel.Conj()).Normalize()
	worldAnchor := b.P.Add(b.Q.Rotate(h.AnchorB))
	return model.Pose{
		P: worldAnchor.Sub(qA.Rotate(h.AnchorA)),
		Q: qA,
	}
}

func hingeRelativeRotation(h model.Hinge, s model.State) model.Quat {
	pose := s.Pose(h.ID)
	var angle float64
	switch pose {
	case model.Pose0:
		return model.QuatIdentity()
	case model.Pose90:
		angle = h.Angle90
		if angle == 0 {
			angle = math.Pi / 2
		}
	case model.Pose180:
		angle = h.Angle180
		if angle == 0 {
			angle = math.Pi
		}
	}
	sign := 1.0
	if h.SignA < 0 {
		sign = -1
	}
	axis := h.AxisA.UnitVector()
	return model.QuatFromAxisAngle(axis, sign*angle)
}
