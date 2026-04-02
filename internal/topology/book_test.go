package topology

import (
	"fmt"
	"infinite-cube/internal/fsm"
	"infinite-cube/internal/model"
	"infinite-cube/internal/validate"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBook(t *testing.T) {
	n := 3
	book := Book(n)
	report := validate.AnalyzeTopology(book)
	for _, issue := range report.Issues {
		fmt.Println(issue)
	}
	if len(report.Issues) != 0 {
		t.Fatalf("Expected no issues, but found %d", len(report.Issues))
	}

	changes := make([]model.HingeChange, n)
	for i := range changes {
		changes[i] = model.HingeChange{
			Hinge: model.HingeID(i * 3),
			To:    model.Pose90,
		}
	}
	move := model.Move{Changes: changes}
	to := model.State{PoseBits: 0b1000001000001}
	isValidTransition := validate.StructuralValidator{}.ValidTransition(book, model.State{}, move, to)
	require.True(t, isValidTransition)

	graph := fsm.Enumerate(book, model.State{}, validate.StructuralValidator{})
	transitions := graph.Edges[model.State{}]
	expected := []model.Transition{
		{
			From: model.State{},
			Move: model.Move{
				Changes: []model.HingeChange{
					{Hinge: 1, To: model.Pose90},
					{Hinge: 2, To: model.Pose90},
				},
			},
			To: model.State{PoseBits: 0b10100},
		},
		{
			From: model.State{},
			Move: model.Move{
				Changes: []model.HingeChange{
					{Hinge: 4, To: model.Pose90},
					{Hinge: 5, To: model.Pose90},
				},
			},
			To: model.State{PoseBits: 0b10100000000},
		},
		{ // The FSM struggles to identify this transition
			From: model.State{},
			Move: model.Move{
				Changes: []model.HingeChange{
					{Hinge: 0, To: model.Pose90},
					{Hinge: 3, To: model.Pose90},
					{Hinge: 6, To: model.Pose90},
				},
			},
			To: model.State{PoseBits: 0b1000001000001},
		},
	}
	require.Equal(t, expected, transitions)
}
