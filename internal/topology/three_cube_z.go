package topology

import "infinite-cube/internal/model"

func ThreeCubeZ() model.Topology {
	cubes := []model.CubeID{0, 1, 2}
	hinges := []model.Hinge{
		{
			ID:      0,
			A:       0,
			B:       1,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
		},
		{
			ID:      1,
			A:       1,
			B:       2,
			AxisA:   model.AxisY,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0.5},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0.5},
		},
	}
	return model.Topology{Cubes: cubes, Hinges: hinges}
}
