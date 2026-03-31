package topology

import (
	"math"

	"infinite-cube/internal/model"
)

// ThreeCubeOpposed90 models 0-1-2 in a line with hinge barrels facing opposite directions.
// PoseA is disengaged (line), PoseB is engaged at 90 degrees.
func ThreeCubeOpposed90() model.Topology {
	cubes := []model.CubeID{0, 1, 2}
	hinges := []model.Hinge{
		{
			ID:      0,
			A:       0,
			B:       1,
			AxisA:   model.AxisZ,
			SignA:   1,
			AngleB:  math.Pi / 2,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0},
		},
		{
			ID:      1,
			A:       1,
			B:       2,
			AxisA:   model.AxisZ,
			SignA:   -1,
			AngleB:  math.Pi / 2,
			AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0},
			AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0},
		},
	}
	return model.Topology{Cubes: cubes, Hinges: hinges}
}
