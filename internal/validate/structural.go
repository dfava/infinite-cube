package validate

import (
	"fmt"
	"math"
	"sort"

	"infinite-cube/internal/kinematics"
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
		if math.IsNaN(h.Angle180) || math.IsInf(h.Angle180, 0) || math.IsNaN(h.Angle90) || math.IsInf(h.Angle90, 0) {
			issues = append(issues, fmt.Sprintf("hinge %d has non-finite Angle180 or Angle90", h.ID))
		}
		if h.Angle180 < 0 || h.Angle180 > math.Pi || h.Angle90 < 0 || h.Angle90 > math.Pi {
			issues = append(issues, fmt.Sprintf("hinge %d has invalid Angle180/Angle90 (expected 0..pi radians)", h.ID))
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

		// Edge-alignment checks
		if issue := checkHingeAlignment(h); issue != "" {
			issues = append(issues, issue)
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

	if len(top.Hinges) > 0 {
		var validMask uint32
		for _, h := range top.Hinges {
			if h.ID < 16 {
				validMask |= (0x3 << (2 * h.ID))
			}
		}
		if s.PoseBits&^validMask != 0 {
			issues = append(issues, "state has bits set for hinges that do not exist in topology")
		}
	} else if s.PoseBits != 0 {
		issues = append(issues, "state has bits set but topology has no hinges")
	}

	// Collision detection
	if len(issues) == 0 {
		solver := kinematics.NewDeterministicSolver()
		poses, err := solver.Poses(top, s)
		if err == nil {
			// Check for cube overlaps. Since each cube is a unit cube centered at Pose.P,
			// two cubes collide if the distance between their centers is less than 1.0.
			// However, adjacent cubes connected by a hinge share a face, so their distance is exactly 1.0.
			// We only flag if distance < 0.99 (allowing some floating point slack).
			ids := make([]model.CubeID, 0, len(poses))
			for id := range poses {
				ids = append(ids, id)
			}
			sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

			for i := 0; i < len(ids); i++ {
				for j := i + 1; j < len(ids); j++ {
					idA, idB := ids[i], ids[j]
					dist := poses[idA].P.Distance(poses[idB].P)
					if dist < 0.99 {
						issues = append(issues, fmt.Sprintf("collision detected between cube %d and cube %d (distance %.3f)", idA, idB, dist))
					}
				}
			}
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

func checkHingeAlignment(h model.Hinge) string {
	// A hinge is edge-aligned if the anchor has two coordinates at +/- 0.5.
	// The third coordinate corresponds to the axis along which the edge runs.
	// The hinge AxisA MUST match this edge direction.

	type coord struct {
		val float64
		idx int
	}
	coords := []coord{{h.AnchorA.X, 0}, {h.AnchorA.Y, 1}, {h.AnchorA.Z, 2}}

	fixed := make([]int, 0)
	var freeIdx int
	for _, c := range coords {
		if math.Abs(math.Abs(c.val)-0.5) < 1e-6 {
			fixed = append(fixed, c.idx)
		} else {
			freeIdx = c.idx
		}
	}

	if len(fixed) < 2 {
		return fmt.Sprintf("hinge %d is not aligned with an edge of cube A (at least two anchor coordinates must be +/- 0.5)", h.ID)
	}

	// For a hinge on an edge, the axis must be the one that is NOT fixed.
	// In some cases (corners), more than 2 coordinates could be +/- 0.5.
	// But normally, a hinge axis is one of X, Y, Z.
	var edgeDir model.Axis
	if len(fixed) == 3 {
		// Corner case: any axis passing through this corner might be considered.
		// However, for an infinite cube, it's usually one of the edges meeting at the corner.
		// We'll allow the axis if it's any of the three.
		// Wait, if it's a corner, we don't know which edge it is.
		// But the user said "hinge is placed on an edge".
		// Let's assume the hinge axis defines the edge.
		edgeDir = h.AxisA
	} else {
		edgeDir = model.Axis(freeIdx)
	}

	if h.AxisA != edgeDir {
		return fmt.Sprintf("hinge %d axis %s does not match edge direction %s of cube A", h.ID, h.AxisA, edgeDir)
	}

	// Shared edge check: in Pose0 (identity rotation), the anchors must refer to the same world point
	// if we assume they are adjacent along some axis.
	// Actually, the simpler check: the distance between AnchorA and AnchorB must be exactly 1.0
	// (or they must be on opposite faces of the cubes that are touching).
	// If cubes are unit and axis-aligned, they touch if their centers are 1 unit apart.
	// In that case, the contact point is AnchorA in A's frame and AnchorB in B's frame.
	// For them to be the "same" point when the cubes are adjacent, they must be "mirrored"
	// across the contact plane.
	// E.g. if contact is at X=0.5 for A and X=-0.5 for B, then AnchorA.X=0.5, AnchorB.X=-0.5,
	// and AnchorA.Y == AnchorB.Y, AnchorA.Z == AnchorB.Z.

	dx := math.Abs(h.AnchorA.X - h.AnchorB.X)
	dy := math.Abs(h.AnchorA.Y - h.AnchorB.Y)
	dz := math.Abs(h.AnchorA.Z - h.AnchorB.Z)

	// One of these should be 1.0, the others should be 0.0.
	sums := dx + dy + dz
	if math.Abs(sums-1.0) > 1e-6 || (dx > 1e-6 && dx < 0.99) || (dy > 1e-6 && dy < 0.99) || (dz > 1e-6 && dz < 0.99) {
		return fmt.Sprintf("hinge %d anchors do not define a shared edge between adjacent cubes", h.ID)
	}

	return ""
}
