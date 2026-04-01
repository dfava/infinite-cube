package topology

import (
	"infinite-cube/internal/model"
)

// BranchingStar returns a star-shaped topology with a central cube and 6 arms.
// This topology is interesting because it reflects radial symmetry found in nature,
// like a starfish or certain blossoms. It creates a "decision" hub at the center
// where a single move can change the orientation of multiple appendages.
// It explores a different kind of "infinity" where the number of possible states
// grows exponentially from a common origin, rather than a linear or circular sequence.
func BranchingStar() model.Topology {
	// 0 is central cube, 1-6 are the arms
	cubes := []model.CubeID{0, 1, 2, 3, 4, 5, 6}
	hinges := []model.Hinge{
		{
			ID:      0,
			A:       0,
			B:       1,
			AxisA:   model.AxisX,
			SignA:   1,
			AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5},
			AnchorB: model.Vec3{X: 0, Y: 0.5, Z: 0.5},
		},
		{
			ID:      1,
			A:       0,
			B:       2,
			AxisA:   model.AxisX,
			SignA:   1,
			AnchorA: model.Vec3{X: 0, Y: -0.5, Z: -0.5},
			AnchorB: model.Vec3{X: 0, Y: -0.5, Z: -0.5},
		},
		{
			ID:      2,
			A:       0,
			B:       3,
			AxisA:   model.AxisY,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: -0.5},
			AnchorB: model.Vec3{X: 0.5, Y: 0, Z: -0.5},
		},
		{
			ID:      3,
			A:       0,
			B:       4,
			AxisA:   model.AxisY,
			SignA:   1,
			AnchorA: model.Vec3{X: -0.5, Y: 0, Z: 0.5},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0.5},
		},
		{
			ID:      4,
			A:       0,
			B:       5,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
		},
		{
			ID:      5,
			A:       0,
			B:       6,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
		},
	}

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
