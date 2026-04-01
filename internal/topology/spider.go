package topology

import (
	"infinite-cube/internal/model"
)

// Spider returns a topology with a 6-cube central ring and 3 branching "legs".
// This topology is interesting because it combines a closed loop (the "body")
// with branching structures (the "legs"). This hybrid approach explores how
// constraints in a central core affect the freedom of movement in appendages.
// It mimics biological structures like arthropods, where a rigid/constrained thorax
// provides the base for articulated limbs.
func Spider() model.Topology {
	// Ring: 0, 1, 2, 3, 4, 5
	// Legs: 6 (off 0), 7 (off 2), 8 (off 4)
	cubes := make([]model.CubeID, 9)
	for i := range 9 {
		cubes[i] = model.CubeID(i)
	}

	hinges := []model.Hinge{
		// Central "Ring" (Open chain in Pose0)
		{ID: 0, A: 0, B: 1, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0, Y: 0.5, Z: 0.5}},
		{ID: 1, A: 1, B: 2, AxisA: model.AxisY, AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0.5}, AnchorB: model.Vec3{X: 0.5, Y: 0, Z: 0.5}},
		{ID: 2, A: 2, B: 3, AxisA: model.AxisZ, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0}, AnchorB: model.Vec3{X: 0.5, Y: 0.5, Z: 0}},
		{ID: 3, A: 3, B: 4, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0, Y: 0.5, Z: 0.5}},
		{ID: 4, A: 4, B: 5, AxisA: model.AxisY, AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0.5}, AnchorB: model.Vec3{X: 0.5, Y: 0, Z: 0.5}},

		// Legs (connected to alternate cubes in the ring)
		{ID: 5, A: 0, B: 6, AxisA: model.AxisY, AnchorA: model.Vec3{X: -0.5, Y: 0, Z: -0.5}, AnchorB: model.Vec3{X: -0.5, Y: 0, Z: -0.5}},
		{ID: 6, A: 2, B: 7, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0, Y: -0.5, Z: -0.5}, AnchorB: model.Vec3{X: 0, Y: -0.5, Z: -0.5}},
		{ID: 7, A: 4, B: 8, AxisA: model.AxisZ, AnchorA: model.Vec3{X: -0.5, Y: -0.5, Z: 0}, AnchorB: model.Vec3{X: -0.5, Y: -0.5, Z: 0}},
	}

	for i := range hinges {
		hinges[i].SignA = 1
	}

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
