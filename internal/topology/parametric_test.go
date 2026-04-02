package topology

import (
	"testing"
)

func TestBookParametric(t *testing.T) {
	tests := []struct {
		pages      int
		wantCubes  int
		wantHinges int
	}{
		{1, 2, 1},
		{2, 4, 4},
		{3, 6, 7},
	}
	for _, tt := range tests {
		top := Book(tt.pages)
		if len(top.Cubes) != tt.wantCubes {
			t.Errorf("Book(%d) cubes = %d, want %d", tt.pages, len(top.Cubes), tt.wantCubes)
		}
		if len(top.Hinges) != tt.wantHinges {
			t.Errorf("Book(%d) hinges = %d, want %d", tt.pages, len(top.Hinges), tt.wantHinges)
		}
	}
}

func TestSnakeParametric(t *testing.T) {
	tests := []struct {
		cubes      int
		wantCubes  int
		wantHinges int
	}{
		{2, 2, 1},
		{4, 4, 3},
		{6, 6, 5},
	}
	for _, tt := range tests {
		top := SnakeChain(tt.cubes)
		if len(top.Cubes) != tt.wantCubes {
			t.Errorf("SnakeChain(%d) cubes = %d, want %d", tt.cubes, len(top.Cubes), tt.wantCubes)
		}
		if len(top.Hinges) != tt.wantHinges {
			t.Errorf("SnakeChain(%d) hinges = %d, want %d", tt.cubes, len(top.Hinges), tt.wantHinges)
		}
	}
}
