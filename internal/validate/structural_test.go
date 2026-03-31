package validate

import (
	"testing"

	"infinite-cube/internal/model"
	"infinite-cube/internal/topology"
)

func TestAnalyzeStateValid(t *testing.T) {
	top := topology.TwoCubeHinge()
	s := model.State{}
	report := AnalyzeState(top, s)
	if len(report.Issues) != 0 {
		t.Fatalf("expected no issues, got %v", report.Issues)
	}
}

func TestAnalyzeStateFlagsOutOfRangeBits(t *testing.T) {
	top := topology.TwoCubeHinge() // Hinge ID 0
	// 2 bits per hinge. Hinge 0 uses bits 0,1.
	// Hinge 1 would use bits 2,3.
	// Setting bits 2 or 3 should trigger error.
	s := model.State{PoseBits: 0b100}
	report := AnalyzeState(top, s)
	if len(report.Issues) == 0 {
		t.Fatalf("expected issues for out-of-range state bit")
	}
}

func TestStructuralValidatorTransition(t *testing.T) {
	top := topology.TwoCubeHinge()
	v := StructuralValidator{}
	from := model.State{}
	mv := model.Move{Changes: []model.HingeChange{{Hinge: 0, To: model.Pose180}}}
	to := from.ApplyMove(mv)
	if !v.ValidTransition(top, from, mv, to) {
		t.Fatalf("expected transition to be valid")
	}
}

func TestAnalyzeTopology(t *testing.T) {
	tests := []struct {
		name     string
		top      model.Topology
		hasIssue bool
	}{
		{
			name:     "Valid Two Cube",
			top:      topology.TwoCubeHinge(),
			hasIssue: false,
		},
		{
			name: "Duplicate Cube ID",
			top: model.Topology{
				Cubes: []model.CubeID{0, 0},
			},
			hasIssue: true,
		},
		{
			name: "Self Hinge",
			top: model.Topology{
				Cubes: []model.CubeID{0},
				Hinges: []model.Hinge{
					{ID: 0, A: 0, B: 0, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5}},
				},
			},
			hasIssue: true,
		},
		{
			name: "Unknown Cube B",
			top: model.Topology{
				Cubes: []model.CubeID{0},
				Hinges: []model.Hinge{
					{ID: 0, A: 0, B: 1, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5}},
				},
			},
			hasIssue: true,
		},
		{
			name: "Duplicate Hinge ID",
			top: model.Topology{
				Cubes: []model.CubeID{0, 1},
				Hinges: []model.Hinge{
					{ID: 0, A: 0, B: 1, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5}},
					{ID: 0, A: 0, B: 1, AxisA: model.AxisY, AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0.5}, AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0.5}},
				},
			},
			hasIssue: true,
		},
		{
			name: "Multiple Hinges Same Pair",
			top: model.Topology{
				Cubes: []model.CubeID{0, 1},
				Hinges: []model.Hinge{
					{ID: 0, A: 0, B: 1, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5}},
					{ID: 1, A: 0, B: 1, AxisA: model.AxisY, AnchorA: model.Vec3{X: 0.5, Y: 0, Z: 0.5}, AnchorB: model.Vec3{X: -0.5, Y: 0, Z: 0.5}},
				},
			},
			hasIssue: true,
		},
		{
			name: "Non-contiguous Hinge ID",
			top: model.Topology{
				Cubes: []model.CubeID{0, 1, 2},
				Hinges: []model.Hinge{
					{ID: 0, A: 0, B: 1, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5}},
					{ID: 2, A: 1, B: 2, AxisA: model.AxisX, AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5}, AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5}},
				},
			},
			hasIssue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := AnalyzeTopology(tt.top)
			if (len(report.Issues) > 0) != tt.hasIssue {
				t.Errorf("expected issues: %v, got: %v", tt.hasIssue, report.Issues)
			}
		})
	}
}

func TestHingeAlignment(t *testing.T) {
	tests := []struct {
		name     string
		hinge    model.Hinge
		hasIssue bool
	}{
		{
			name: "Valid Edge Alignment",
			hinge: model.Hinge{
				ID:      0,
				AxisA:   model.AxisX,
				AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5},
				AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5},
			},
			hasIssue: false,
		},
		{
			name: "Invalid Anchor (not on edge)",
			hinge: model.Hinge{
				ID:      0,
				AxisA:   model.AxisX,
				AnchorA: model.Vec3{X: 0, Y: 0.4, Z: 0.5},
				AnchorB: model.Vec3{X: 0, Y: -0.4, Z: 0.5},
			},
			hasIssue: true,
		},
		{
			name: "Axis Mismatch",
			hinge: model.Hinge{
				ID:      0,
				AxisA:   model.AxisY, // Should be AxisX for AnchorA with fixed Y,Z
				AnchorA: model.Vec3{X: 0, Y: 0.5, Z: 0.5},
				AnchorB: model.Vec3{X: 0, Y: -0.5, Z: 0.5},
			},
			hasIssue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue := checkHingeAlignment(tt.hinge)
			if (issue != "") != tt.hasIssue {
				t.Errorf("expected issue: %v, got: %v", tt.hasIssue, issue)
			}
		})
	}
}
