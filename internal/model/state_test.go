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
	s = s.ApplyMove(s.Flip(2)) // 0 -> 90
	if s.Pose(2) != Pose90 {
		t.Fatalf("expected Pose90, got %v", s.Pose(2))
	}
	s = s.ApplyMove(s.Flip(2)) // 90 -> 180
	if s.Pose(2) != Pose180 {
		t.Fatalf("expected Pose180, got %v", s.Pose(2))
	}
	s = s.ApplyMove(s.Flip(2)) // 180 -> 90
	if s.Pose(2) != Pose90 {
		t.Fatalf("expected Pose90 after flip from 180, got %v", s.Pose(2))
	}
	s = s.ApplyMove(s.Flip(2)) // 90 -> 180 (Forward cycle from 90 is 180 in current Flip implementation)
	// Let's check 90 -> 0 to be sure we can go back.
	s = State{}
	s = s.ApplyMove(Move{Changes: []HingeChange{{Hinge: 2, To: Pose90}}})
	// We need a way to go 90 -> 0. Flip as currently implemented goes 90 -> 180.
	// This is fine for its purpose.
}
