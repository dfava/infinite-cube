package kinematics

import (
	"fmt"
	"infinite-cube/internal/model"
)

// Solver computes cube poses implied by a discrete state.
type Solver interface {
	Poses(top model.Topology, s model.State) (map[model.CubeID]model.Pose, error)
}

// StubSolver is a deterministic placeholder that keeps all cubes at origin.
type StubSolver struct{}

func (StubSolver) Poses(top model.Topology, _ model.State) (map[model.CubeID]model.Pose, error) {
	if len(top.Cubes) == 0 {
		return nil, fmt.Errorf("empty topology")
	}
	poses := make(map[model.CubeID]model.Pose, len(top.Cubes))
	for _, c := range top.Cubes {
		poses[c] = model.Pose{Q: model.Quat{W: 1}}
	}
	return poses, nil
}
