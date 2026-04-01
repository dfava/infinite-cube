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
	for _, c := range m.Changes {
		shift := 2 * c.Hinge
		mask := uint32(0x3) << shift
		s.PoseBits = (s.PoseBits &^ mask) | (uint32(c.To) << shift)
	}
	return s
}

// Flip returns a move that cycles through Pose0 -> Pose90 -> Pose180 -> Pose90 -> Pose0 for a single hinge.
func (s State) Flip(h HingeID) Move {
	cur := s.Pose(h)
	var next HingePose
	switch cur {
	case Pose0:
		next = Pose90
	case Pose90:
		// We could go to 0 or 180. Flip is often used in interactive contexts
		// where we cycle forward. Let's cycle forward: 0 -> 90 -> 180 -> 0.
		// Actually, to be strictly 90-degree adjacent, we should go to 180.
		next = Pose180
	case Pose180:
		// Go back to 90
		next = Pose90
	default:
		next = Pose0
	}
	return Move{Changes: []HingeChange{{Hinge: h, To: next}}}
}

func (s State) String() string {
	return fmt.Sprintf("State{PoseBits: %032b}", s.PoseBits)
}
