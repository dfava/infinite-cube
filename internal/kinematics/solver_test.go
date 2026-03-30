package kinematics

import (
	"testing"

	"infinite-cube/internal/model"
	"infinite-cube/internal/topology"
)

func TestDeterministicSolverSimplePoseA(t *testing.T) {
	top := topology.Simple()
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

	// Simple topology hinge: axis X, sign +1, so offset uses +Y in cube A frame.
	if p1.P != (model.Vec3{Y: 1}) {
		t.Fatalf("expected cube 1 at +Y offset, got %+v", p1.P)
	}
	if p1.Q != (model.Quat{W: 1}) {
		t.Fatalf("expected PoseA to keep identity orientation, got %+v", p1.Q)
	}
}

func TestDeterministicSolverSimplePoseBFlipsOrientation(t *testing.T) {
	top := topology.Simple()
	solver := NewDeterministicSolver()
	state := model.State{}.ApplyMove(model.Move{Hinge: 0, To: model.PoseB})

	poses, err := solver.Poses(top, state)
	if err != nil {
		t.Fatalf("Poses returned error: %v", err)
	}

	p1 := poses[1]
	// 180deg around +X => quaternion approximately (0,1,0,0) up to sign.
	if !p1.Q.AlmostEqual(model.Quat{X: 1}, 1e-6) {
		t.Fatalf("expected cube 1 orientation to rotate around X for PoseB, got %+v", p1.Q)
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
