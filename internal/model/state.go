package model

import "fmt"

// State is an FSM node. Each hinge i uses 2 bits in PoseBits.
type State struct {
	PoseBits uint32
}

// Pose returns the current discrete pose for a hinge.
func (s State) Pose(h HingeID) HingePose {
	shift := 2 * h
	return HingePose((s.PoseBits >> shift) & 0x3)
}

// ApplyMove returns a new state with the move applied.
func (s State) ApplyMove(m Move) State {
	shift := 2 * m.Hinge
	mask := uint32(0x3) << shift
	s.PoseBits = (s.PoseBits &^ mask) | (uint32(m.To) << shift)
	return s
}

// Flip cycles through PoseA -> PoseB -> PoseC -> PoseA.
func (s State) Flip(h HingeID) State {
	cur := s.Pose(h)
	var next HingePose
	switch cur {
	case PoseA:
		next = PoseB
	case PoseB:
		next = PoseC
	default:
		next = PoseA
	}
	return s.ApplyMove(Move{Hinge: h, To: next})
}

func (s State) String() string {
	return fmt.Sprintf("State{PoseBits: %032b}", s.PoseBits)
}
