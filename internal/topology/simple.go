package topology

import (
	"infinite-cube/internal/model"
)

func Simple() model.Topology {
	cubes := []model.CubeID{0, 1}
	hinges := []model.Hinge{
		{ID: 0, A: 0, B: 1, AxisA: model.AxisX, SignA: 1},
	}
	return model.Topology{cubes, hinges}
}
