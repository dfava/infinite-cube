package fsm

import "infinite-cube/internal/model"

// Graph stores reachable states and outgoing transitions.
type Graph struct {
	Nodes map[model.State]struct{}
	Edges map[model.State][]model.Transition
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[model.State]struct{}),
		Edges: make(map[model.State][]model.Transition),
	}
}

func (g *Graph) NumNodes() int {
	return len(g.Nodes)
}

func (g *Graph) NumEdges() int {
	n := 0
	for _, transitions := range g.Edges {
		n += len(transitions)
	}
	return n
}
