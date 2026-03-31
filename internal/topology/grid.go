package topology

import (
	"infinite-cube/internal/model"
)

// Grid2x2 returns a 2x2 grid of cubes (4 cubes total) connected in a square.
// This is interesting because it introduces redundant paths and cycles in a 2D-like structure.
// While each individual hinge has a simple motion, the collective constraints of the
// square loop mean that not all combinations of hinge poses are physically possible.
// It explores the "puzzle" of how a flat sheet-like structure can be deformed
// when constrained by a closed loop.
func Grid2x2() model.Topology {
	// Cubes:
	// (0,0) - 0, (1,0) - 1
	// (0,1) - 2, (1,1) - 3
	cubes := []model.CubeID{0, 1, 2, 3}

	hinges := []model.Hinge{
		// Horizontal hinges
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
			A:       2,
			B:       3,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
		},
		// Vertical hinges
		{
			ID:      2,
			A:       0,
			B:       2,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
		},
		{
			ID:      3,
			A:       1,
			B:       3,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
		},
	}

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
