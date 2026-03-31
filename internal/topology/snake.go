package topology

import (
	"infinite-cube/internal/model"
)

// SnakeChain returns a linear chain of cubes connected by hinges with alternating axes.
// This topology is interesting because it mimics a biological structure like a vertebral column
// or a polypeptide chain. The alternating axes (X then Y then Z) create a complex range of
// motion, allowing the "snake" to curl into various compact shapes, almost like protein folding.
// It explores the "puzzle" aspect of trying to reach a specific 3D configuration from a straight line.
func SnakeChain(numCubes int) model.Topology {
	if numCubes < 2 {
		numCubes = 2
	}
	cubes := make([]model.CubeID, numCubes)
	for i := range numCubes {
		cubes[i] = model.CubeID(i)
	}

	hinges := make([]model.Hinge, numCubes-1)
	axes := []model.Axis{model.AxisX, model.AxisY, model.AxisZ}

	for i := range numCubes - 1 {
		axis := axes[i%3]
		var anchorA, anchorB model.Vec3

		// We connect cubes along the Z-axis, but the hinge axis alternates.
		// If we imagine cubes of size 1x1x1 centered at origin:
		// Cube i's "forward" face is at Z=0.5
		// Cube i+1's "backward" face is at Z=-0.5
		// The hinge axis sits ON the face.

		switch axis {
		case model.AxisX:
			// Hinge along X axis on the Z-face
			anchorA = model.Vec3{X: 0, Y: 0.5, Z: 0.5}
			anchorB = model.Vec3{X: 0, Y: 0.5, Z: -0.5}
		case model.AxisY:
			// Hinge along Y axis on the Z-face
			anchorA = model.Vec3{X: 0.5, Y: 0, Z: 0.5}
			anchorB = model.Vec3{X: 0.5, Y: 0, Z: -0.5}
		case model.AxisZ:
			// Hinge along Z axis on the X-face (to change direction of the chain)
			anchorA = model.Vec3{X: 0.5, Y: 0.5, Z: 0}
			anchorB = model.Vec3{X: -0.5, Y: 0.5, Z: 0}
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
