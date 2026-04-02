package topology

import (
	"infinite-cube/internal/fsm"
	"infinite-cube/internal/model"
	"infinite-cube/internal/validate"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTwoCubeHinge(t *testing.T) {
	cube := TwoCubeHinge()
	report := validate.AnalyzeTopology(cube)
	require.Empty(t, report.Issues, "Expected no issues in TwoCubeHinge topology")

	// The initial state for fsm.Enumerate is all hinges at Pose0 by default.
	graph := fsm.Enumerate(cube, model.State{}, &validate.StructuralValidator{}, 1)

	// 1 hinge, 3 possible poses: Pose0, Pose90, Pose180.
	// Since there are no collisions or constraints, all 3 should be reachable.
	require.Equal(t, 3, graph.NumNodes(), "Expected 3 nodes in the FSM graph")

	// Transitions:
	// Pose0 -> Pose90
	// Pose90 -> Pose0
	// Pose90 -> Pose180
	// Pose180 -> Pose90
	// Total: 4 edges.
	require.Equal(t, 4, graph.NumEdges(), "Expected 4 edges in the FSM graph")
}
