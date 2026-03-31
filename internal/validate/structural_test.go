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
	top := topology.TwoCubeHinge()
	s := model.State{PoseBits: 0b10}
	report := AnalyzeState(top, s)
	if len(report.Issues) == 0 {
		t.Fatalf("expected issues for out-of-range state bit")
	}
}

func TestStructuralValidatorTransition(t *testing.T) {
	top := topology.TwoCubeHinge()
	v := StructuralValidator{}
	from := model.State{}
	mv := model.Move{Hinge: 0, To: model.PoseB}
	to := from.ApplyMove(mv)
	if !v.ValidTransition(top, from, mv, to) {
		t.Fatalf("expected transition to be valid")
	}
}
