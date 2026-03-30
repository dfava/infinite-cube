package topology

import (
	"infinite-cube/internal/model"
)

// InfiniteCube8 returns a placeholder 8-cube/8-hinge layout.
// Replace hinge definitions with your physically exact layout.
func InfiniteCube8() model.Topology {
	cubes := []model.CubeID{0, 1, 2, 3, 4, 5, 6, 7}
	hinges := []model.Hinge{
		{ID: 0, A: 0, B: 1, AxisA: model.AxisX, SignA: 1},
		{ID: 1, A: 1, B: 3, AxisA: model.AxisY, SignA: 1},
		{ID: 2, A: 3, B: 2, AxisA: model.AxisX, SignA: -1},
		{ID: 3, A: 2, B: 0, AxisA: model.AxisY, SignA: -1},
		{ID: 4, A: 4, B: 5, AxisA: model.AxisX, SignA: 1},
		{ID: 5, A: 5, B: 7, AxisA: model.AxisY, SignA: 1},
		{ID: 6, A: 7, B: 6, AxisA: model.AxisX, SignA: -1},
		{ID: 7, A: 6, B: 4, AxisA: model.AxisY, SignA: -1},
	}
	return model.Topology{Cubes: cubes, Hinges: hinges}
}
