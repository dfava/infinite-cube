package main

import (
	"fmt"

	"infinite-cube/internal/fsm"
	"infinite-cube/internal/model"
	"infinite-cube/internal/topology"
	"infinite-cube/internal/validate"
)

func main() {
	top := topology.InfiniteCube8()
	start := model.State{}
	validator := validate.PermissiveValidator{}

	graph := fsm.Enumerate(top, start, validator)

	edges := 0
	for _, out := range graph.Edges {
		edges += len(out)
	}

	fmt.Printf("infinite-cube FSM scaffold\n")
	fmt.Printf("cubes=%d hinges=%d reachable_states=%d transitions=%d\n", len(top.Cubes), len(top.Hinges), len(graph.Nodes), edges)
}
