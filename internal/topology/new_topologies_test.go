package topology

import (
	"infinite-cube/internal/kinematics"
	"infinite-cube/internal/model"
	"infinite-cube/internal/validate"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTopologiesValid(t *testing.T) {
	tests := []struct {
		name string
		top  func() model.Topology
	}{
		{"TetrahedralCluster", TetrahedralCluster},
		{"ExpandingStar", ExpandingStar},
		{"MobiusStrip", MobiusStrip},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			top := tt.top()

			// 1. Structural Validation
			report := validate.AnalyzeTopology(top)
			for _, issue := range report.Issues {
				t.Errorf("Structural issue in %s: %s", tt.name, issue)
			}
			require.Empty(t, report.Issues, "Expected no structural issues in %s topology", tt.name)

			// 2. Kinematic Consistency (Initial State)
			solver := kinematics.NewDeterministicSolver()
			_, err := solver.Poses(top, model.State{})
			require.NoError(t, err, "Expected kinematic consistency in %s at Pose0", tt.name)
		})
	}
}
