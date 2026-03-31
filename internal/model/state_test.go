package model

import (
	"testing"
)

func TestApplyMoveAndPose(t *testing.T) {
	s := State{}
	if got := s.Pose(2); got != Pose0 {
		t.Fatalf("expected Pose0, got %v", got)
	}

	s = s.ApplyMove(Move{Changes: []HingeChange{{Hinge: 2, To: Pose180}}})
	if got := s.Pose(2); got != Pose180 {
		t.Fatalf("expected Pose180, got %v", got)
	}

	s = s.ApplyMove(Move{Changes: []HingeChange{{Hinge: 2, To: Pose0}}})
	if got := s.Pose(2); got != Pose0 {
		t.Fatalf("expected Pose0 after reset, got %v", got)
	}
}

func TestFlip(t *testing.T) {
	s := State{}
	s = s.ApplyMove(s.Flip(2)) // to Pose90
	s = s.ApplyMove(s.Flip(2)) // to Pose180
	s = s.ApplyMove(s.Flip(2)) // to Pose0
	if s.PoseBits != 0 {
		t.Fatalf("expected initial state, got %v", s)
	}
}
