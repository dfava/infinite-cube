package topology

import (
	"infinite-cube/internal/model"
)

// InfiniteCube8 returns a placeholder 8-cube/8-hinge layout.
// Replace hinge definitions with your physically exact layout.
func InfiniteCube8() model.Topology {
	cubes := []model.CubeID{0, 1, 2, 3, 4, 5, 6, 7}
	// We are starting with 8 cubes and NO hinges.
	// As we add hinges together, we'll fill this list.
	hinges := []model.Hinge{
		{
			ID:      0,
			A:       0,
			B:       1,
			AxisA:   model.AxisX,
			SignA:   1,
			AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5},
			AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5},
		},
		{
			ID:      1,
			A:       0,
			B:       2,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
		},
		{
			ID:      2,
			A:       1,
			B:       3,
			AxisA:   model.AxisZ,
			SignA:   -1,
			AnchorA: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
		},
		// Right column
		{
			ID:      3,
			A:       3,
			B:       4,
			AxisA:   model.AxisY,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: -0.5},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: -0.5},
		},
		{
			ID:      4,
			A:       2,
			B:       5,
			AxisA:   model.AxisY,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: -0.5},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: -0.5},
		},
		{
			ID:      5,
			A:       6,
			B:       4,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
		},
		{
			ID:      6,
			A:       7,
			B:       5,
			AxisA:   model.AxisZ,
			SignA:   -1,
			AnchorA: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
		},
		{
			ID:      7,
			A:       7,
			B:       6,
			AxisA:   model.AxisX,
			SignA:   1,
			AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5},
			AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5},
		},
	}

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
