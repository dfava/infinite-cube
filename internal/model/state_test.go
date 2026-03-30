package model

import (
	"testing"
)

func TestApplyMoveAndPose(t *testing.T) {
	s := State{}
	if got := s.Pose(2); got != PoseA {
		t.Fatalf("expected PoseA, got %v", got)
	}

	s = s.ApplyMove(Move{Hinge: 2, To: PoseB})
	if got := s.Pose(2); got != PoseB {
		t.Fatalf("expected PoseB, got %v", got)
	}

	s = s.ApplyMove(Move{Hinge: 2, To: PoseA})
	if got := s.Pose(2); got != PoseA {
		t.Fatalf("expected PoseA after reset, got %v", got)
	}
}

func TestFlip(t *testing.T) {
	s := State{}
	s = s.Flip(2)
	s = s.Flip(2)
	if s.PoseBits != 0 {
		t.Fatalf("expected initial state, got %v", s)
	}
}
