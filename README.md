# infinite-cube

A Go scaffold for modeling an infinite cube as a finite state machine with deterministic kinematics and legality checks.

## Layout

- `cmd/infinite-cube`: CLI entrypoint
- `internal/model`: core data types (`State`, `Move`, `Topology`, etc.)
- `internal/topology`: constructors for specific toy topologies
- `internal/kinematics`: pose computation interfaces and stubs
- `internal/validate`: state/transition legality interfaces and implementations
- `internal/fsm`: graph generation/enumeration

## Next steps

1. Encode the exact hinge layout for your physical cube in `internal/topology`.
2. Implement deterministic pose generation in `internal/kinematics`.
3. Implement collision checks in `internal/validate`.
4. Expand the CLI to export the state graph as JSON/DOT.
