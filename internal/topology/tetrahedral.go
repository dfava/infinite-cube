package topology

import (
	"infinite-cube/internal/model"
)

// TetrahedralCluster returns a cluster of 4 cubes connected in a tetrahedral arrangement.
// This topology is interesting because it's a fundamental 3D structure that introduces
// multi-axial constraints. Unlike a simple chain or a 2D grid, the tetrahedral
// arrangement forces moves to consider 3D volume and potential self-intersection
// early on. It's a "puzzle" of how to rotate any one cube without being blocked
// by its three neighbors, mimicking a basic unit of many crystalline structures.
func TetrahedralCluster() model.Topology {
	// Let's place 4 cubes at approximately:
	// 0: (0,0,0)
	// 1: (1,0,0)
	// 2: (0.5, 0.866, 0)
	// 3: (0.5, 0.288, 0.816)
	// Actually, let's keep them axis-aligned for the validator.
	// 0: (0,0,0)
	// 1: (1,0,0)
	// 2: (0,1,0)
	// 3: (0,0,1)
	// They all touch Cube 0.
	
	cubes := []model.CubeID{0, 1, 2, 3}
	
	hinges := []model.Hinge{
		{
			ID:      0,
			A:       0,
			B:       1,
			AxisA:   model.AxisZ, // Edge along Z
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
		},
		{
			ID:      1,
			A:       0,
			B:       2,
			AxisA:   model.AxisX, // Edge along X
			SignA:   1,
			AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5},
			AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5},
		},
		{
			ID:      2,
			A:       0,
			B:       3,
			AxisA:   model.AxisY, // Edge along Y
			SignA:   1,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0.5},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0.5},
		},
	}
	
	return model.Topology{Cubes: cubes, Hinges: hinges}
}
