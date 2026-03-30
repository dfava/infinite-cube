package model

import "fmt"

// State is an FSM node. Bit i stores pose of hinge i.
type State struct {
	PoseBits uint16
}

// Pose returns the current discrete pose for a hinge.
func (s State) Pose(h HingeID) HingePose {
	mask := uint16(1) << h
	if s.PoseBits&mask != 0 {
		return PoseB
	}
	return PoseA
}

// ApplyMove returns a new state with the move applied.
func (s State) ApplyMove(m Move) State {
	mask := uint16(1) << m.Hinge
	if m.To == PoseB {
		s.PoseBits |= mask
	} else {
		s.PoseBits &^= mask
	}
	return s
}

// Flip returns a new state where hinge h toggles pose.
func (s State) Flip(h HingeID) State {
	mask := uint16(1) << h
	s.PoseBits ^= mask
	return s
}

func (s State) String() string {
	return fmt.Sprintf("State{PoseBits: %016b}", s.PoseBits)
}
