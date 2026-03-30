# infinite-cube

A Go scaffold for modeling an infinite cube as a finite state machine with deterministic kinematics and legality checks.

## Layout

- `cmd/infinite-cube`: basic CLI entrypoint
- `cmd/topology-viz`: local web visualizer for topology editing + validator feedback
- `internal/model`: core data types (`State`, `Move`, `Topology`, etc.)
- `internal/topology`: constructors for specific toy topologies
- `internal/kinematics`: pose computation interfaces and stubs
- `internal/validate`: legality interfaces and implementations
- `internal/fsm`: graph generation/enumeration

## Run

```bash
go run ./cmd/infinite-cube
```

```bash
go run ./cmd/topology-viz
# open http://localhost:8080
```

## Topology Visualizer

The visualizer lets you:

- Load `simple` or `infinite8` topology presets
- Edit hinge fields `ID`, `A`, `B`, `AxisA`, `SignA`
- Toggle hinge pose bits (`PoseB`) for a candidate state
- Run validator diagnostics from `validate.AnalyzeState`

`AxisA` and `SignA` are shown directly on hinge labels:

- `AxisA` picks the local rotation axis on cube `A`
- `SignA` sets positive direction convention (`+1` or `-1`)

## Next steps

1. Encode the exact physical hinge layout in `internal/topology`.
2. Implement deterministic pose generation in `internal/kinematics`.
3. Replace permissive legality with collision-aware checks in `internal/validate`.
4. Export state graphs as JSON/DOT from a CLI command.
