package topology

import (
	"infinite-cube/internal/model"
)

// FigureEight returns a double-loop topology of 10 cubes connected in two 6-cube rings sharing a common side.
// This topology is interesting because it introduces "crossing" paths and more complex cycles.
// It feels like a more advanced version of the RingLoop, where a move in the shared central bar
// affects the constraints and state of both loops simultaneously.
// It explores the idea of interconnected systems where changes in one domain propagate to another.
func FigureEight() model.Topology {
	// Let's envision two rings of 6 cubes sharing two cubes (and one hinge).
	// Ring 1: 0, 1, 2, 3, 4, 5
	// Ring 2: 4, 5, 6, 7, 8, 9
	// Shared cubes: 4, 5
	// Shared hinge: between 4 and 5
	cubes := make([]model.CubeID, 10)
	for i := range 10 {
		cubes[i] = model.CubeID(i)
	}

	hinges := []model.Hinge{
		// Ring 1 (Simple linear chain)
		{ID: 0, A: 0, B: 1, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0.5}},
		{ID: 1, A: 1, B: 2, AxisA: model.AxisY, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0.5}},
		{ID: 2, A: 2, B: 3, AxisA: model.AxisZ, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0.5, Y: 0.5, Z: -0.5}},
		{ID: 3, A: 3, B: 4, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0.5}},
		{ID: 4, A: 4, B: 5, AxisA: model.AxisY, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0.5}},

		// Ring 2 branches off
		{ID: 5, A: 5, B: 6, AxisA: model.AxisZ, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0.5, Y: 0.5, Z: -0.5}},
		{ID: 6, A: 6, B: 7, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0.5}},
		{ID: 7, A: 7, B: 8, AxisA: model.AxisY, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0.5}},
		{ID: 8, A: 8, B: 9, AxisA: model.AxisZ, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0.5, Y: 0.5, Z: -0.5}},
	}

	for i := range hinges {
		hinges[i].SignA = 1
	}

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
