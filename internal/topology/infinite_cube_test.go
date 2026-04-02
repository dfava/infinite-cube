package topology

import (
	"fmt"
	"infinite-cube/internal/fsm"
	"infinite-cube/internal/model"
	"infinite-cube/internal/validate"
	"testing"

	"github.com/stretchr/testify/require"
)

// This test helps prevent regressions in fsm.Enumerate.
func TestInfiniteCube(t *testing.T) {
	cube := InfiniteCube8()
	report := validate.AnalyzeTopology(cube)
	for _, issue := range report.Issues {
		fmt.Println(issue)
	}
	require.Equal(t, 0, len(report.Issues), "Expected no issues, but found %d", len(report.Issues))

	graph := fsm.Enumerate(cube, model.State{}, &validate.StructuralValidator{}, 2)
	require.Equal(t, 24, graph.NumNodes())
	require.Equal(t, 64, graph.NumEdges())
}
