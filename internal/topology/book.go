package topology

import (
	"infinite-cube/internal/model"
)

func Book(n int) model.Topology {
	if n < 1 {
		n = 1
	}
	cubes := make([]model.CubeID, n*2)
	for i := range n * 2 {
		cubes[i] = model.CubeID(i)
	}

	// Evens connect to odds. Evens connect to events. Odds connect to odds
	// 0---h---1
	// |       |
	// 2---h---3
	// |       |
	// 4---h---	5
	hinges := make([]model.Hinge, n+2*(n-1))
	for iteration := range n {
		hinge := iteration * 3
		cube := iteration * 2
		hinges[hinge] = model.Hinge{
			ID:      model.HingeID(hinge),
			A:       model.CubeID(cube),
			B:       model.CubeID(cube + 1),
			AxisA:   model.AxisY,
			SignA:   1,
			AnchorA: model.Vec3{0.5, 0, -0.5},
			AnchorB: model.Vec3{-0.5, 0, -0.5},
		}
		if iteration == n-1 {
			break
		}
		hinges[hinge+1] = model.Hinge{
			ID:      model.HingeID(hinge + 1),
			A:       model.CubeID(cube),
			B:       model.CubeID(cube + 2),
			AxisA:   model.AxisX,
			SignA:   -1,
			AnchorA: model.Vec3{0, -0.5, 0.5},
			AnchorB: model.Vec3{0, 0.5, 0.5},
		}
		hinges[hinge+2] = model.Hinge{
			ID:      model.HingeID(hinge + 2),
			A:       model.CubeID(cube + 1),
			B:       model.CubeID(cube + 3),
			AxisA:   model.AxisX,
			SignA:   -1,
			AnchorA: model.Vec3{0, -0.5, 0.5},
			AnchorB: model.Vec3{0, 0.5, 0.5},
		}
	}
	return model.Topology{Cubes: cubes, Hinges: hinges}
}
