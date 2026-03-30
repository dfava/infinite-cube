package validate

import "infinite-cube/internal/model"

// Validator decides if states and transitions are legal.
type Validator interface {
	ValidState(top model.Topology, s model.State) bool
	ValidTransition(top model.Topology, from model.State, mv model.Move, to model.State) bool
}

// PermissiveValidator allows everything. Replace with collision-aware checks.
type PermissiveValidator struct{}

func (PermissiveValidator) ValidState(_ model.Topology, _ model.State) bool {
	return true
}

func (PermissiveValidator) ValidTransition(_ model.Topology, _ model.State, _ model.Move, _ model.State) bool {
	return true
}
