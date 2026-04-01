package topology

import "infinite-cube/internal/model"

// TwoCubeHinge is a minimal two-cube topology with one hinge attachment.
func TwoCubeHinge() model.Topology {
	cubes := []model.CubeID{0, 1}
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
	}
	return model.Topology{Cubes: cubes, Hinges: hinges}
}

// This is an interesting configuration in the sense that
// it is used to demonstrate a possible issue with collision detection.
// The hinge is placed in a way that the cubes start side by side at Pose0.
// At Pose180, the cubes are again side by side.
// However, to get from Pose0 to Pose180, the hinge must go through Pose90,
// and at Pose90, the cubes collide.
// So, if the collision detection algorithm is not implemented correctly,
// the system will fail to exclude the Topology at Pose180.
func TwoCubeHingeThrough() model.Topology {
	cubes := []model.CubeID{0, 1}
	hinges := []model.Hinge{
		{
			ID:      0,
			A:       0,
			B:       1,
			AxisA:   model.AxisZ,
			SignA:   -1,
			AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
		},
	}
	return model.Topology{Cubes: cubes, Hinges: hinges}
}
