package topology

import (
	"infinite-cube/internal/model"
)

// TwistedLoop returns a 7-cube loop with a "twist" in the hinge axes.
// This topology is interesting because it's an odd-numbered loop.
// With 6 cubes and alternating axes, you can form a symmetric hexagon.
// Adding a 7th cube and twisting the axes (X, Y, Z pattern doesn't line up perfectly
// on a 7-cube loop) creates a structure that might require a full "lap" to return
// to its original orientation, similar to a Möbius strip.
// It explores the concept of non-orientable-like behavior in a discrete system
// and the "infinity" of a path that takes longer than expected to close.
func TwistedLoop7() model.Topology {
	cubes := make([]model.CubeID, 7)
	for i := range 7 {
		cubes[i] = model.CubeID(i)
	}

	hinges := []model.Hinge{
		{ID: 0, A: 0, B: 1, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0.5}},
		{ID: 1, A: 1, B: 2, AxisA: model.AxisY, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0.5}},
		{ID: 2, A: 2, B: 3, AxisA: model.AxisZ, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0.5, Y: 0.5, Z: -0.5}},
		{ID: 3, A: 3, B: 4, AxisA: model.AxisX, AnchorA: model.Vec3{X: -0.5, Y: -0.5, Z: -0.5}, AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: -0.5}},
		{ID: 4, A: 4, B: 5, AxisA: model.AxisY, AnchorA: model.Vec3{X: -0.5, Y: -0.5, Z: -0.5}, AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: -0.5}},
		{ID: 5, A: 5, B: 6, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0.5}},
	}
	// Note: All SignA default to 0, which is handled as +1 if the code doesn't specify,
	// but I'll set them to 1 explicitly for clarity.
	for i := range hinges {
		hinges[i].SignA = 1
	}

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
