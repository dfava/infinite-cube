package topology

import (
	"infinite-cube/internal/model"
)

// HTree returns a fractal-like H-tree topology of cubes.
// This topology is interesting because it mimics recursive self-similar structures
// often found in nature, like the branching of lungs or vascular systems.
// It creates a structure with a high surface-area-to-volume ratio as it expands,
// and its movements can feel very coordinated or chaotic depending on which
// level of the recursion is being manipulated.
func HTree() model.Topology {
	// Let's build a simple 2-level H-tree
	// Level 0: Central bar (2 cubes)
	// Level 1: Four perpendicular bars at the ends (4 more cubes)
	// Total: 6 cubes (0, 1) central, (2, 3) off 0, (4, 5) off 1
	cubes := []model.CubeID{0, 1, 2, 3, 4, 5}

	hinges := []model.Hinge{
		// Central bar
		{
			ID:      0,
			A:       0,
			B:       1,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
		},
		// Branches off cube 0
		{
			ID:      1,
			A:       0,
			B:       2,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
		},
		{
			ID:      2,
			A:       0,
			B:       3,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
		},
		// Branches off cube 1
		{
			ID:      3,
			A:       1,
			B:       4,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
		},
		{
			ID:      4,
			A:       1,
			B:       5,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
			AnchorB: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
		},
	}

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
