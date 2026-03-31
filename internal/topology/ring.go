package topology

import (
	"infinite-cube/internal/model"
)

// RingLoop returns a closed loop of 6 cubes.
// This topology is interesting because 6 is the minimum number of cubes
// to form a loop that can fold in 3D (a 4-cube loop is rigid if restricted to 90 degree folds).
// It's like a simplified version of the Yoshimoto Cube, exploring the transition between
// a regular hexagon configuration and a folded, dense 3D structure.
// This has a "circular" infinity feeling where states repeat in cycles.
func RingLoop6() model.Topology {
	cubes := []model.CubeID{0, 1, 2, 3, 4, 5}
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
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
		},
		{
			ID:      2,
			A:       2,
			B:       3,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
		},
		{
			ID:      3,
			A:       3,
			B:       4,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
		},
		{
			ID:      4,
			A:       4,
			B:       5,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
		},
		{
			ID:      5,
			A:       5,
			B:       0,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
		},
	}

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
