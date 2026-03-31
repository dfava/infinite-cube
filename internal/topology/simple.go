package topology

import (
	"infinite-cube/internal/model"
)

func Simple() model.Topology {
	cubes := []model.CubeID{0, 1}
	hinges := []model.Hinge{
		{
			ID:      0,
			A:       0,
			B:       1,
			AxisA:   model.AxisX,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0},
		},
	}
	return model.Topology{cubes, hinges}
}
