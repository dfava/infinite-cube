# Infinite Cube Simulator

A Go-based simulator for modeling an "Infinite Cube" (a mechanical puzzle of 8 hinged cubes) and other modular cube topologies. It features deterministic kinematics, finite state machine (FSM) exploration, and a web-based visualizer.

## Capabilities

The simulator provides a robust framework for exploring the configuration space of hinged cubes:

- **Topology Modeling**: Define arbitrary connections between cubes using hinges with specific anchor points, rotation axes, and angle limits.
- **Deterministic Kinematics**: Automatically computes the 3D position and orientation (pose) of every cube based on the hinge states (0°, 90°, or 180°). It handles kinematic propagation and detects inconsistent cycles.
- **FSM Enumeration**: Explores the entire reachable state space of a topology. It supports both single-hinge flips and simultaneous multi-hinge moves (essential for complex topologies where single moves might be physically blocked).
- **Structural Validation**: Checks for:
  - **Hinge Alignment**: Ensures hinges are correctly placed on cube edges.
  - **Connectivity**: Validates that the topology forms a valid kinematic tree or graph.
  - **Collision Detection**: Identifies physical overlaps between cubes in any given state.
- **Web Visualization**: A real-time interactive dashboard to edit topologies, toggle hinge states, and see immediate validation feedback.

## Project Structure

- `cmd/infinite-cube`: CLI entrypoint for analysis and simulation.
- `cmd/topology-viz`: Local web server for interactive topology editing and visualization.
- `internal/model`: Core data types (`State`, `Pose`, `Topology`, `Vec3`, `Quat`).
- `internal/topology`: Library of predefined topologies (from simple 2-cube pairs to the full 8-cube Infinite Cube).
- `internal/kinematics`: Implementation of the deterministic pose solver.
- `internal/validate`: Structural and collision-aware legality checks.
- `internal/fsm`: Graph generation and state space enumeration logic.

## Getting Started

### Prerequisites
- Go 1.23 or later.

### Running the Visualizer
The web visualizer is the best way to explore existing topologies:

```bash
go run ./cmd/topology-viz
```
Then open [http://localhost:8080](http://localhost:8080) in your browser.

### Running the CLI
```bash
go run ./cmd/infinite-cube
```

## Existing Topologies

The `internal/topology` package includes several presets:

1. **Two Cube**: Minimal two-cube pair connected by a single hinge.
2. **Three Cube Line**: Three cubes connected in a straight line.
3. **Three Cube Z**: Three cubes with hinges on different axes, creating a "Z" shape when folded.
4. **Three Cube Opposed**: Three cubes with hinges on opposite faces/edges.
5. **Infinite Cube (8-cube)**: The full 8-cube layout representing the classic "Infinite Cube" puzzle.

## Technical Details

### Kinematic Solver
The solver uses a Breadth-First Search (BFS) starting from a root cube (usually Cube 0) to propagate transforms across the hinge graph. Each hinge defines a relative rotation and translation between two cubes. The solver uses quaternions for robust rotation math and ensures that any cycles in the topology are kinematically consistent.

### FSM Enumeration
The explorer identifies all valid configurations by treating the puzzle as a Finite State Machine. A transition between states is considered valid if the resulting configuration has no cube collisions and follows the kinematic constraints of the hinges. It can detect "locked" states or complex paths that require multiple hinges to move in unison.

### Collision Detection
Collision detection is performed by calculating the distance between cube centers in 3D space. Since all cubes are unit cubes (1x1x1), a distance of less than 1.0 (with a small epsilon for floating point math) indicates a collision.

## Contributing

1. Encode new physical hinge layouts in `internal/topology`.
2. Extend `internal/validate` with more sophisticated physics or constraint checks.
3. Improve the FSM exploration heuristics in `internal/fsm`.
