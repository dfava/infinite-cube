package main

import (
	"bytes"
	"encoding/json"
	"infinite-cube/internal/topology"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIParametric(t *testing.T) {
	tests := []struct {
		preset     string
		param      string
		wantCubes  int
		wantHinges int
	}{
		{"h-tree", "1", 2, 1},
		{"h-tree", "2", 8, 7},
		{"book", "2", 4, 4},
		{"snake", "3", 3, 2},
		{"tetrahedral-cluster", "", 4, 3},
		{"expanding-star", "", 13, 12},
		{"mobius-strip", "", 10, 9},
		{"", "", 2, 1}, // Default two-cube-hinge
	}

	for _, tt := range tests {
		url := "/api/topology"
		if tt.preset != "" {
			url += "?preset=" + tt.preset
			if tt.param != "" {
				url += "&param=" + tt.param
			}
		}
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		handleTopology(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GET %s returned status %d", url, w.Code)
			continue
		}

		var resp topologyJSON
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Errorf("GET %s failed to decode JSON: %v", url, err)
			continue
		}

		if len(resp.Cubes) != tt.wantCubes {
			t.Errorf("GET %s got %d cubes, want %d", url, len(resp.Cubes), tt.wantCubes)
		}
		if len(resp.Hinges) != tt.wantHinges {
			t.Errorf("GET %s got %d hinges, want %d", url, len(resp.Hinges), tt.wantHinges)
		}
	}
}

func TestEnumerateParametricCache(t *testing.T) {
	// 1. Enumerate h-tree level 1
	top1 := topology.HTree(1)
	req1 := validateRequest{
		Topology:   toTopologyJSON(top1),
		PresetName: "h-tree",
		Param:      1,
	}
	body1, _ := json.Marshal(req1)
	w1 := httptest.NewRecorder()
	handleEnumerate(w1, httptest.NewRequest("POST", "/api/enumerate", bytes.NewReader(body1)))

	var resp1 enumerateResponse
	if err := json.NewDecoder(w1.Body).Decode(&resp1); err != nil {
		t.Fatalf("failed to decode resp1: %v", err)
	}

	// 2. Enumerate h-tree level 2
	top2 := topology.HTree(2)
	req2 := validateRequest{
		Topology:   toTopologyJSON(top2),
		PresetName: "h-tree",
		Param:      2,
	}
	body2, _ := json.Marshal(req2)
	w2 := httptest.NewRecorder()
	handleEnumerate(w2, httptest.NewRequest("POST", "/api/enumerate", bytes.NewReader(body2)))

	var resp2 enumerateResponse
	if err := json.NewDecoder(w2.Body).Decode(&resp2); err != nil {
		t.Fatalf("failed to decode resp2: %v", err)
	}

	if len(resp1.States) == len(resp2.States) {
		t.Errorf("Expected different number of states for h-tree levels 1 and 2, but both got %d", len(resp1.States))
	}
}
