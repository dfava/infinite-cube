package topology

import (
	"infinite-cube/internal/model"
)

// HTree returns a fractal-like H-tree topology of cubes of a given level.
// Levels = 1: A central bar (2 cubes)
// Levels = 2: Four perpendicular bars at the ends (total 6 cubes if no extra cubes)
// Higher levels: Branches have more cubes to make room for lower levels.
// This topology mimics recursive self-similar structures found in nature.
func HTree(levels int) model.Topology {
	if levels < 1 {
		return model.Topology{}
	}

	var cubes []model.CubeID
	var hinges []model.Hinge

	// Helper to add a chain of cubes
	addChain := func(startCube model.CubeID, numCubes int, isHorizontal bool, side int) model.CubeID {
		prev := startCube
		for range numCubes {
			next := model.CubeID(len(cubes))
			cubes = append(cubes, next)

			var anchorA, anchorB model.Vec3
			if isHorizontal {
				// Chain along X
				anchorA = model.Vec3{X: float64(side) * 0.5, Y: 0.5, Z: 0}
				anchorB = model.Vec3{X: -float64(side) * 0.5, Y: 0.5, Z: 0}
			} else {
				// Chain along Y
				anchorA = model.Vec3{X: 0.5, Y: float64(side) * 0.5, Z: 0}
				anchorB = model.Vec3{X: 0.5, Y: -float64(side) * 0.5, Z: 0}
			}

			hinges = append(hinges, model.Hinge{
				ID:      model.HingeID(len(hinges)),
				A:       prev,
				B:       next,
				AxisA:   model.AxisZ,
				SignA:   1,
				AnchorA: anchorA,
				AnchorB: anchorB,
			})
			prev = next
		}
		return prev
	}

	// Cube 0 and 1 form the base of the central bar
	cubes = append(cubes, 0, 1)

	hinges = append(hinges, model.Hinge{
		ID:      0,
		A:       0,
		B:       1,
		AxisA:   model.AxisZ,
		SignA:   1,
		AnchorA: model.Vec3{X: 0.5, Y: 0.5, Z: 0},
		AnchorB: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
	})

	// centralBarExtraCubes ensures Level 1 is long enough
	// For levels=3, we need 2^1 = 2 extra cubes on each side of central bar?
	// Actually, let's use 2^(levels-currentLevel-1) extra cubes.
	centralExtra := 0
	if levels > 1 {
		centralExtra = 1 << (levels - 2)
	}

	rootL := addChain(0, centralExtra, true, -1)
	rootR := addChain(1, centralExtra, true, 1)

	var addBranches func(parent model.CubeID, currentLevel int, isHorizontal bool, side int)
	addBranches = func(parent model.CubeID, currentLevel int, isHorizontal bool, side int) {
		if currentLevel >= levels {
			return
		}

		// At each end, we branch into TWO directions.
		// These branches themselves might need to be long.
		numExtra := 0
		if currentLevel < levels-1 {
			numExtra = 1 << (levels - currentLevel - 2)
		}

		// Child 1
		c1Start := model.CubeID(len(cubes))
		cubes = append(cubes, c1Start)

		var anchorP1, anchorC1 model.Vec3
		if isHorizontal {
			anchorP1 = model.Vec3{X: float64(side) * 0.5, Y: 0.5, Z: 0}
			anchorC1 = model.Vec3{X: float64(side) * 0.5, Y: -0.5, Z: 0}
		} else {
			anchorP1 = model.Vec3{X: 0.5, Y: float64(side) * 0.5, Z: 0}
			anchorC1 = model.Vec3{X: -0.5, Y: float64(side) * 0.5, Z: 0}
		}

		hinges = append(hinges, model.Hinge{
			ID:      model.HingeID(len(hinges)),
			A:       parent,
			B:       c1Start,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: anchorP1,
			AnchorB: anchorC1,
		})

		leaf1 := addChain(c1Start, numExtra, !isHorizontal, 1)

		// Child 2
		c2Start := model.CubeID(len(cubes))
		cubes = append(cubes, c2Start)

		var anchorP2, anchorC2 model.Vec3
		if isHorizontal {
			anchorP2 = model.Vec3{X: float64(side) * 0.5, Y: -0.5, Z: 0}
			anchorC2 = model.Vec3{X: float64(side) * 0.5, Y: 0.5, Z: 0}
		} else {
			anchorP2 = model.Vec3{X: -0.5, Y: float64(side) * 0.5, Z: 0}
			anchorC2 = model.Vec3{X: 0.5, Y: float64(side) * 0.5, Z: 0}
		}

		hinges = append(hinges, model.Hinge{
			ID:      model.HingeID(len(hinges)),
			A:       parent,
			B:       c2Start,
			AxisA:   model.AxisZ,
			SignA:   1,
			AnchorA: anchorP2,
			AnchorB: anchorC2,
		})

		leaf2 := addChain(c2Start, numExtra, !isHorizontal, -1)

		addBranches(leaf1, currentLevel+1, !isHorizontal, 1)
		addBranches(leaf2, currentLevel+1, !isHorizontal, -1)
	}

	addBranches(rootL, 1, true, -1)
	addBranches(rootR, 1, true, 1)

	return model.Topology{Cubes: cubes, Hinges: hinges}
}
