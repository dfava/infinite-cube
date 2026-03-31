package kinematics

import (
	"testing"

	"infinite-cube/internal/model"
	"infinite-cube/internal/topology"
)

func TestDeterministicSolverSimplePose0(t *testing.T) {
	top := topology.TwoCubeHinge()
	solver := NewDeterministicSolver()

	poses, err := solver.Poses(top, model.State{})
	if err != nil {
		t.Fatalf("Poses returned error: %v", err)
	}

	p0, ok := poses[0]
	if !ok {
		t.Fatalf("missing pose for cube 0")
	}
	p1, ok := poses[1]
	if !ok {
		t.Fatalf("missing pose for cube 1")
	}

	if p0.P != (model.Vec3{}) {
		t.Fatalf("expected cube 0 at origin, got %+v", p0.P)
	}
	if p0.Q != (model.Quat{W: 1}) {
		t.Fatalf("expected cube 0 identity orientation, got %+v", p0.Q)
	}

	if p1.P != (model.Vec3{X: 1}) {
		t.Fatalf("expected cube 1 at +X offset from anchors, got %+v", p1.P)
	}
	if p1.Q != (model.Quat{W: 1}) {
		t.Fatalf("expected Pose0 to keep identity orientation, got %+v", p1.Q)
	}
}

func TestDeterministicSolverSimplePose180FlipsOrientation(t *testing.T) {
	top := topology.TwoCubeHinge()
	solver := NewDeterministicSolver()
	state := model.State{}.ApplyMove(model.Move{Changes: []model.HingeChange{{Hinge: 0, To: model.Pose180}}})

	poses, err := solver.Poses(top, state)
	if err != nil {
		t.Fatalf("Poses returned error: %v", err)
	}

	p1 := poses[1]
	// Now TwoCubeHinge uses AxisZ.
	// 180deg around +Z => quaternion approximately (0,0,0,1) up to sign.
	if !p1.Q.AlmostEqual(model.Quat{Z: 1}, 1e-6) {
		t.Fatalf("expected cube 1 orientation to rotate around Z for Pose180, got %+v", p1.Q)
	}
}

func TestDeterministicSolverDisconnectedComponents(t *testing.T) {
	top := model.Topology{Cubes: []model.CubeID{0, 1, 2}}
	solver := NewDeterministicSolver()

	poses, err := solver.Poses(top, model.State{})
	if err != nil {
		t.Fatalf("Poses returned error: %v", err)
	}

	if poses[0].P.X != 0 || poses[1].P.X != solver.ComponentSpacing || poses[2].P.X != 2*solver.ComponentSpacing {
		t.Fatalf("unexpected disconnected component placement: c0=%+v c1=%+v c2=%+v", poses[0].P, poses[1].P, poses[2].P)
	}
}

func TestDeterministicSolverRejectsUnknownCubeReference(t *testing.T) {
	top := model.Topology{
		Cubes: []model.CubeID{0},
		Hinges: []model.Hinge{
			{ID: 0, A: 0, B: 9, AxisA: model.AxisX, SignA: 1},
		},
	}
	solver := NewDeterministicSolver()

	if _, err := solver.Poses(top, model.State{}); err == nil {
		t.Fatalf("expected error for unknown cube reference")
	}
}

func TestThreeCubeOpposed90HasLineLLinePattern(t *testing.T) {
	solver := NewDeterministicSolver()
	top := topology.ThreeCubeOpposed90()

	classify := func(bits uint32) (bool, error) {
		poses, err := solver.Poses(top, model.State{PoseBits: bits})
		if err != nil {
			return false, err
		}
		return collinear3(poses[0].P, poses[1].P, poses[2].P), nil
	}

	// 00, 01, 10, 11 in binary-pose layout (1 bit per hinge)
	// Now mapped to 2 bits per hinge.
	// Pose0=0, Pose180=2.
	// State 00: H0=0, H1=0 => bits 0
	// State 01: H0=2, H1=0 => bits 2
	// State 10: H0=0, H1=2 => bits 8 (2 << (2*1))
	// State 11: H0=2, H1=2 => bits 10 (2 | (2 << 2))

	c00, err := classify(0)
	if err != nil {
		t.Fatalf("state 00 error: %v", err)
	}
	c01, err := classify(2)
	if err != nil {
		t.Fatalf("state 01 error: %v", err)
	}
	c10, err := classify(8)
	if err != nil {
		t.Fatalf("state 10 error: %v", err)
	}
	c11, err := classify(10)
	if err != nil {
		t.Fatalf("state 11 error: %v", err)
	}

	if !c00 || c01 || c10 || !c11 {
		t.Fatalf("unexpected shape pattern; want line/L/L/line, got 00=%v 01=%v 10=%v 11=%v", c00, c01, c10, c11)
	}
}

func collinear3(a, b, c model.Vec3) bool {
	v1 := b.Sub(a)
	v2 := c.Sub(b)
	cross := model.Vec3{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: v1.Z*v2.X - v1.X*v2.Z,
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
	return cross.Distance(model.Vec3{}) <= 1e-6
}

func TestDeterministicSolverInconsistentCycle(t *testing.T) {
	// Cube 0 --H0-- Cube 1
	//   \           /
	//    H1       H2
	//     \       /
	//      Cube 2
	// We'll set up anchors so that it's consistent in Pose0 but inconsistent if we change one hinge.
	top := model.Topology{
		Cubes: []model.CubeID{0, 1, 2},
		Hinges: []model.Hinge{
			{
				ID:      0,
				A:       0,
				B:       1,
				AxisA:   model.AxisZ,
				AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0},
				AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0},
			},
			{
				ID:      1,
				A:       0,
				B:       2,
				AxisA:   model.AxisZ,
				AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0},
				AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0},
			},
			{
				ID:      2,
				A:       1,
				B:       2,
				AxisA:   model.AxisZ,
				AnchorA: model.Vec3{X: -0.5, Y: 0.5, Z: 0},
				AnchorB: model.Vec3{X: 0.5, Y: -0.5, Z: 0},
			},
		},
	}
	solver := NewDeterministicSolver()

	// Pose0 should be consistent (a triangle of cubes)
	_, err := solver.Poses(top, model.State{})
	if err != nil {
		t.Fatalf("expected consistent Pose0, got %v", err)
	}

	// Now make one hinge Pose180. This should break the cycle consistency.
	state := model.State{}.ApplyMove(model.Move{Changes: []model.HingeChange{{Hinge: 0, To: model.Pose180}}})
	_, err = solver.Poses(top, state)
	if err == nil {
		t.Fatalf("expected inconsistent kinematic cycle error")
	}
}

func TestDeterministicSolverEmptyTopology(t *testing.T) {
	solver := NewDeterministicSolver()
	if _, err := solver.Poses(model.Topology{}, model.State{}); err == nil {
		t.Fatalf("expected error for empty topology")
	}
}

func TestDeterministicSolverDuplicateCubeID(t *testing.T) {
	solver := NewDeterministicSolver()
	top := model.Topology{Cubes: []model.CubeID{0, 0}}
	if _, err := solver.Poses(top, model.State{}); err == nil {
		t.Fatalf("expected error for duplicate cube ID")
	}
}
