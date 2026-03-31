package model

// CubeID identifies one of the rigid mini-cubes.
type CubeID uint8

// HingeID identifies one hinge in the fixed topology.
type HingeID uint8

// HingePose is the discrete hinge pose used by the FSM.
type HingePose uint8

const (
	Pose0   HingePose = iota // 0 degrees
	Pose90                   // 90 degrees
	Pose180                  // 180 degrees
)

// Hinge describes one connection in the immutable toy topology.
type Hinge struct {
	ID    HingeID
	A     CubeID
	B     CubeID
	AxisA Axis
	SignA int8 // +1 or -1 for orientation convention
	// Angle90/Angle180 are hinge rotations (radians) used when state pose is Pose90/Pose180.
	// Zero defaults to pi/2 for Angle90 and pi for Angle180.
	Angle90  float64
	Angle180 float64
	// AnchorA/AnchorB are hinge attachment points in each cube's local frame.
	AnchorA Vec3
	AnchorB Vec3
}

// Topology contains fixed cube and hinge connectivity.
type Topology struct {
	Cubes  []CubeID
	Hinges []Hinge
}

// Move represents one hinge move.
type Move struct {
	Hinge HingeID
	To    HingePose
}

// Transition is a directed edge in the FSM graph.
type Transition struct {
	From State
	Move Move
	To   State
}
