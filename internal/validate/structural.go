package validate

import (
	"fmt"
	"math"
	"sort"

	"infinite-cube/internal/model"
)

// StructuralValidator enforces basic topology/state consistency checks.
type StructuralValidator struct{}

func (StructuralValidator) ValidState(top model.Topology, s model.State) bool {
	return len(AnalyzeState(top, s).Issues) == 0
}

func (StructuralValidator) ValidTransition(top model.Topology, from model.State, mv model.Move, to model.State) bool {
	if !(StructuralValidator{}).ValidState(top, from) || !(StructuralValidator{}).ValidState(top, to) {
		return false
	}
	if !hasHinge(top, mv.Hinge) {
		return false
	}
	expected := from.ApplyMove(mv)
	if expected != to {
		return false
	}
	return to.Pose(mv.Hinge) == mv.To
}

// DiagnosticReport provides human-readable validation feedback.
type DiagnosticReport struct {
	Issues []string
}

// AnalyzeTopology checks static topology consistency.
func AnalyzeTopology(top model.Topology) DiagnosticReport {
	issues := make([]string, 0)

	if len(top.Cubes) == 0 {
		issues = append(issues, "topology has no cubes")
	}
	if len(top.Hinges) == 0 {
		issues = append(issues, "topology has no hinges")
	}
	if len(top.Hinges) > 16 {
		issues = append(issues, "topology has more than 16 hinges; State.PoseBits currently supports up to 16")
	}

	cubeSet := make(map[model.CubeID]struct{}, len(top.Cubes))
	for _, c := range top.Cubes {
		if _, exists := cubeSet[c]; exists {
			issues = append(issues, fmt.Sprintf("cube ID %d appears more than once", c))
			continue
		}
		cubeSet[c] = struct{}{}
	}

	hingeIDs := make(map[model.HingeID]struct{}, len(top.Hinges))
	pairs := make(map[[2]model.CubeID]model.HingeID, len(top.Hinges))
	for _, h := range top.Hinges {
		if _, exists := hingeIDs[h.ID]; exists {
			issues = append(issues, fmt.Sprintf("hinge ID %d appears more than once", h.ID))
		} else {
			hingeIDs[h.ID] = struct{}{}
		}
		if h.ID >= 16 {
			issues = append(issues, fmt.Sprintf("hinge ID %d exceeds PoseBits capacity (max 15)", h.ID))
		}
		if _, ok := cubeSet[h.A]; !ok {
			issues = append(issues, fmt.Sprintf("hinge %d references unknown cube A=%d", h.ID, h.A))
		}
		if _, ok := cubeSet[h.B]; !ok {
			issues = append(issues, fmt.Sprintf("hinge %d references unknown cube B=%d", h.ID, h.B))
		}
		if h.A == h.B {
			issues = append(issues, fmt.Sprintf("hinge %d has same endpoint cube (%d)", h.ID, h.A))
		}
		if h.AxisA > model.AxisZ {
			issues = append(issues, fmt.Sprintf("hinge %d has invalid AxisA value %d", h.ID, h.AxisA))
		}
		if h.SignA != 1 && h.SignA != -1 {
			issues = append(issues, fmt.Sprintf("hinge %d has invalid SignA value %d (expected +1 or -1)", h.ID, h.SignA))
		}
		if !vecFinite(h.AnchorA) {
			issues = append(issues, fmt.Sprintf("hinge %d has non-finite AnchorA", h.ID))
		}
		if !vecFinite(h.AnchorB) {
			issues = append(issues, fmt.Sprintf("hinge %d has non-finite AnchorB", h.ID))
		}

		p := canonicalPair(h.A, h.B)
		if prev, exists := pairs[p]; exists {
			issues = append(issues, fmt.Sprintf("hinge %d duplicates cube pair (%d,%d) already used by hinge %d", h.ID, p[0], p[1], prev))
		} else {
			pairs[p] = h.ID
		}
	}

	if len(top.Hinges) > 0 {
		ids := make([]int, 0, len(hingeIDs))
		for id := range hingeIDs {
			ids = append(ids, int(id))
		}
		sort.Ints(ids)
		for i := range ids {
			if ids[i] != i {
				issues = append(issues, "hinge IDs should be contiguous starting at 0 for predictable state indexing")
				break
			}
		}
	}

	return DiagnosticReport{Issues: issues}
}

// AnalyzeState checks topology and state-bit consistency.
func AnalyzeState(top model.Topology, s model.State) DiagnosticReport {
	report := AnalyzeTopology(top)
	issues := append([]string{}, report.Issues...)

	if len(top.Hinges) < 16 {
		mask := uint16(0xFFFF) << len(top.Hinges)
		if s.PoseBits&mask != 0 {
			issues = append(issues, "state has bits set for hinges that do not exist in topology")
		}
	}

	return DiagnosticReport{Issues: issues}
}

func hasHinge(top model.Topology, id model.HingeID) bool {
	for _, h := range top.Hinges {
		if h.ID == id {
			return true
		}
	}
	return false
}

func canonicalPair(a, b model.CubeID) [2]model.CubeID {
	if a <= b {
		return [2]model.CubeID{a, b}
	}
	return [2]model.CubeID{b, a}
}

func vecFinite(v model.Vec3) bool {
	return !math.IsNaN(v.X) && !math.IsInf(v.X, 0) &&
		!math.IsNaN(v.Y) && !math.IsInf(v.Y, 0) &&
		!math.IsNaN(v.Z) && !math.IsInf(v.Z, 0)
}
