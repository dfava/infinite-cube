package topology

import "infinite-cube/internal/model"

// ThreeCubeLine creates cubes 0-1-2 connected in a straight line in PoseA.
func ThreeCubeLine() model.Topology {
	cubes := []model.CubeID{0, 1, 2}
	hinges := []model.Hinge{
		{
			ID:      0,
			A:       0,
			B:       1,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0},
		},
		{
			ID:      1,
			A:       1,
			B:       2,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0},
		},
	}
	return model.Topology{Cubes: cubes, Hinges: hinges}
}

// ThreeCubeL creates cubes 0-1-2 where the second hinge exits cube 1 orthogonally.
func ThreeCubeL() model.Topology {
	cubes := []model.CubeID{0, 1, 2}
	hinges := []model.Hinge{
		{
			ID:      0,
			A:       0,
			B:       1,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0},
		},
		{
			ID:      1,
			A:       1,
			B:       2,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0},
		},
	}
	return model.Topology{Cubes: cubes, Hinges: hinges}
}
