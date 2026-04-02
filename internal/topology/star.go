package topology

import (
	"infinite-cube/internal/model"
)

// ExpandingStar returns a central cube with chains expanding from all six faces.
// This topology is interesting because it mimics biomimicry of something like
// a radiolarian or an echinoderm. The central hub is a "point of origin" for
// complex multi-appendage motion. It explores how coordination between
// several limbs can prevent entanglement or, conversely, create a cage.
func ExpandingStar() model.Topology {
	cubes := []model.CubeID{0}
	hinges := []model.Hinge{}

	// Hub cube 0.
	// We'll add 6 limbs, each 2 cubes long.
	// Limb 1: +X
	// Limb 2: -X
	// Limb 3: +Y
	// Limb 4: -Y
	// Limb 5: +Z
	// Limb 6: -Z

	// Limb 1 (+X direction)
	cubes = append(cubes, 1, 2)
	hinges = append(hinges, model.Hinge{
		ID: 0, A: 0, B: 1, AxisA: model.AxisZ, SignA: 1,
		AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
		AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
	})
	hinges = append(hinges, model.Hinge{
		ID: 1, A: 1, B: 2, AxisA: model.AxisY, SignA: 1,
		AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0.5},
		AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0.5},
	})

	// Limb 2 (-X direction)
	cubes = append(cubes, 3, 4)
	hinges = append(hinges, model.Hinge{
		ID: 2, A: 0, B: 3, AxisA: model.AxisZ, SignA: 1,
		AnchorA: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
		AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
	})
	hinges = append(hinges, model.Hinge{
		ID: 3, A: 3, B: 4, AxisA: model.AxisY, SignA: 1,
		AnchorA: model.Vec3{X: -0.5, Y: 0, Z: -0.5},
		AnchorB: model.Vec3{X: 0.5, Y: 0, Z: -0.5},
	})

	// Limb 3 (+Y direction)
	cubes = append(cubes, 5, 6)
	hinges = append(hinges, model.Hinge{
		ID: 4, A: 0, B: 5, AxisA: model.AxisX, SignA: 1,
		AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5},
		AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5},
	})
	hinges = append(hinges, model.Hinge{
		ID: 5, A: 5, B: 6, AxisA: model.AxisZ, SignA: 1,
		AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
		AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
	})

	// Limb 4 (-Y direction)
	cubes = append(cubes, 7, 8)
	hinges = append(hinges, model.Hinge{
		ID: 6, A: 0, B: 7, AxisA: model.AxisX, SignA: 1,
		AnchorA: model.Vec3{X: 0, Y: -0.5, Z: -0.5},
		AnchorB: model.Vec3{X: 0, Y: 0.5, Z: -0.5},
	})
	hinges = append(hinges, model.Hinge{
		ID: 7, A: 7, B: 8, AxisA: model.AxisZ, SignA: 1,
		AnchorA: model.Vec3{X: -0.5, Y: -0.5, Z: 0},
		AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
	})

	// Limb 5 (+Z direction)
	cubes = append(cubes, 9, 10)
	hinges = append(hinges, model.Hinge{
		ID: 8, A: 0, B: 9, AxisA: model.AxisY, SignA: 1,
		AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0.5},
		AnchorB: model.Vec3{X: 0.5, Y: 0, Z: -0.5},
	})
	hinges = append(hinges, model.Hinge{
		ID: 9, A: 9, B: 10, AxisA: model.AxisX, SignA: 1,
		AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5},
		AnchorB: model.Vec3{X: 0, Y: 0.5, Z: -0.5},
	})

	// Limb 6 (-Z direction)
	cubes = append(cubes, 11, 12)
	hinges = append(hinges, model.Hinge{
		ID: 10, A: 0, B: 11, AxisA: model.AxisY, SignA: 1,
		AnchorA: model.Vec3{X: -0.5, Y: 0, Z: -0.5},
		AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0.5},
	})
	hinges = append(hinges, model.Hinge{
		ID: 11, A: 11, B: 12, AxisA: model.AxisX, SignA: 1,
		AnchorA: model.Vec3{X: 0, Y: -0.5, Z: -0.5},
		AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5},
	})

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
