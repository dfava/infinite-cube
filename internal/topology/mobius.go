package topology

import (
	"infinite-cube/internal/model"
)

// MobiusStrip returns a strip of 10 cubes that, when manipulated, can form a Mobius loop.
// This topology is interesting because it's a 3D realization of a non-orientable surface.
// While it starts as a linear chain to avoid immediate kinematic conflicts, its hinge
// configuration (alternating axes and specific anchors) is designed to allow it to
// twist and close onto itself in a non-trivial way. It mimics the properties of a
// Mobius strip using rigid components, exploring how local constraints lead to
// global topological puzzles.
func MobiusStrip() model.Topology {
	// We use enough cubes to allow for a 180-degree twist and closure.
	numCubes := 10
	cubes := make([]model.CubeID, numCubes)
	for i := range numCubes {
		cubes[i] = model.CubeID(i)
	}

	hinges := make([]model.Hinge, numCubes-1)

	for i := 0; i < numCubes-1; i++ {
		// Standard connection: cubes adjacent along X, but hinge axis rotates
		// every few cubes to facilitate the "twist".
		axis := model.AxisZ
		if i >= 4 && i <= 6 {
			axis = model.AxisY
		}

		var anchorA, anchorB model.Vec3
		if axis == model.AxisZ {
			anchorA = model.Vec3{X: 0.5, Y: 0.5, Z: 0}
			anchorB = model.Vec3{X: -0.5, Y: 0.5, Z: 0}
		} else {
			anchorA = model.Vec3{X: 0.5, Y: 0, Z: 0.5}
			anchorB = model.Vec3{X: -0.5, Y: 0, Z: 0.5}
		}

		hinges[i] = model.Hinge{
			ID:      model.HingeID(i),
			A:       model.CubeID(i),
			B:       model.CubeID(i + 1),
			AxisA:   axis,
			SignA:   1,
			AnchorA: anchorA,
			AnchorB: anchorB,
		}
	}

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
