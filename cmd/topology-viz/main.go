package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"infinite-cube/internal/kinematics"
	"infinite-cube/internal/model"
	"infinite-cube/internal/topology"
	"infinite-cube/internal/validate"
)

//go:embed web/*
var webFS embed.FS

type topologyJSON struct {
	Cubes  []int       `json:"cubes"`
	Hinges []hingeJSON `json:"hinges"`
}

type hingeJSON struct {
	ID    int    `json:"id"`
	A     int    `json:"a"`
	B     int    `json:"b"`
	AxisA string `json:"axisA"`
	SignA int8   `json:"signA"`
}

type validateRequest struct {
	Topology topologyJSON `json:"topology"`
	PoseBits uint16       `json:"poseBits"`
}

type validateResponse struct {
	Valid       bool     `json:"valid"`
	Issues      []string `json:"issues"`
	HingeCount  int      `json:"hingeCount"`
	CubeCount   int      `json:"cubeCount"`
	PoseBits    uint16   `json:"poseBits"`
	PoseBitsBin string   `json:"poseBitsBin"`
	PresetUsed  string   `json:"presetUsed,omitempty"`
}

type vec3JSON struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type quatJSON struct {
	W float64 `json:"w"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type cubePoseJSON struct {
	Cube int      `json:"cube"`
	P    vec3JSON `json:"p"`
	Q    quatJSON `json:"q"`
}

type posesResponse struct {
	PoseBits uint16         `json:"poseBits"`
	Poses    []cubePoseJSON `json:"poses"`
}

type apiError struct {
	Error string `json:"error"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/topology", handleTopology)
	mux.HandleFunc("/api/validate", handleValidate)
	mux.HandleFunc("/api/poses", handlePoses)
	mux.Handle("/", http.FileServer(http.FS(webFS)))

	addr := ":8080"
	log.Printf("topology-viz listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func handleTopology(w http.ResponseWriter, r *http.Request) {
	preset := r.URL.Query().Get("preset")
	if preset == "" {
		preset = "infinite8"
	}

	var top model.Topology
	switch preset {
	case "simple":
		top = topology.Simple()
	case "infinite8":
		top = topology.InfiniteCube8()
	default:
		http.Error(w, "unknown preset", http.StatusBadRequest)
		return
	}

	writeJSON(w, toTopologyJSON(top))
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req validateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	top, err := fromTopologyJSON(req.Topology)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := model.State{PoseBits: req.PoseBits}
	report := validate.AnalyzeState(top, s)
	resp := validateResponse{
		Valid:       len(report.Issues) == 0,
		Issues:      report.Issues,
		HingeCount:  len(top.Hinges),
		CubeCount:   len(top.Cubes),
		PoseBits:    req.PoseBits,
		PoseBitsBin: strconv.FormatUint(uint64(req.PoseBits), 2),
	}
	writeJSON(w, resp)
}

func handlePoses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONWithStatus(w, http.StatusMethodNotAllowed, apiError{Error: "method not allowed"})
		return
	}

	var req validateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONWithStatus(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}

	top, err := fromTopologyJSON(req.Topology)
	if err != nil {
		writeJSONWithStatus(w, http.StatusBadRequest, apiError{Error: err.Error()})
		return
	}

	state := model.State{PoseBits: req.PoseBits}
	if report := validate.AnalyzeState(top, state); len(report.Issues) != 0 {
		writeJSONWithStatus(w, http.StatusBadRequest, apiError{Error: report.Issues[0]})
		return
	}

	solver := kinematics.NewDeterministicSolver()
	poses, err := solver.Poses(top, state)
	if err != nil {
		writeJSONWithStatus(w, http.StatusBadRequest, apiError{Error: err.Error()})
		return
	}

	resp := posesResponse{
		PoseBits: req.PoseBits,
		Poses:    make([]cubePoseJSON, 0, len(poses)),
	}
	for _, c := range top.Cubes {
		p := poses[c]
		resp.Poses = append(resp.Poses, cubePoseJSON{
			Cube: int(c),
			P: vec3JSON{
				X: p.P.X,
				Y: p.P.Y,
				Z: p.P.Z,
			},
			Q: quatJSON{
				W: p.Q.W,
				X: p.Q.X,
				Y: p.Q.Y,
				Z: p.Q.Z,
			},
		})
	}
	writeJSON(w, resp)
}

func toTopologyJSON(top model.Topology) topologyJSON {
	cubes := make([]int, 0, len(top.Cubes))
	for _, c := range top.Cubes {
		cubes = append(cubes, int(c))
	}
	hinges := make([]hingeJSON, 0, len(top.Hinges))
	for _, h := range top.Hinges {
		hinges = append(hinges, hingeJSON{
			ID:    int(h.ID),
			A:     int(h.A),
			B:     int(h.B),
			AxisA: h.AxisA.String(),
			SignA: h.SignA,
		})
	}
	return topologyJSON{Cubes: cubes, Hinges: hinges}
}

func fromTopologyJSON(tj topologyJSON) (model.Topology, error) {
	cubes := make([]model.CubeID, 0, len(tj.Cubes))
	for _, c := range tj.Cubes {
		if c < 0 || c > 255 {
			return model.Topology{}, fmt.Errorf("cube ID %d out of range", c)
		}
		cubes = append(cubes, model.CubeID(c))
	}
	hinges := make([]model.Hinge, 0, len(tj.Hinges))
	for _, h := range tj.Hinges {
		if h.ID < 0 || h.ID > 255 {
			return model.Topology{}, fmt.Errorf("hinge ID %d out of range", h.ID)
		}
		if h.A < 0 || h.A > 255 || h.B < 0 || h.B > 255 {
			return model.Topology{}, fmt.Errorf("hinge %d has cube endpoint out of range", h.ID)
		}
		axis, err := model.AxisFromString(h.AxisA)
		if err != nil {
			return model.Topology{}, fmt.Errorf("hinge %d: %w", h.ID, err)
		}
		hinges = append(hinges, model.Hinge{
			ID:    model.HingeID(h.ID),
			A:     model.CubeID(h.A),
			B:     model.CubeID(h.B),
			AxisA: axis,
			SignA: h.SignA,
		})
	}
	return model.Topology{Cubes: cubes, Hinges: hinges}, nil
}

func writeJSON(w http.ResponseWriter, v any) {
	writeJSONWithStatus(w, http.StatusOK, v)
}

func writeJSONWithStatus(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
	}
}
