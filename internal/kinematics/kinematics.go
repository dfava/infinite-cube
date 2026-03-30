package kinematics

import "infinite-cube/internal/model"

// Solver computes cube poses implied by a discrete state.
type Solver interface {
	Poses(top model.Topology, s model.State) (map[model.CubeID]model.Pose, error)
}
