package model

// Vec3 represents a position in 3D space.
type Vec3 struct {
	X float64
	Y float64
	Z float64
}

// Quat represents an orientation quaternion.
type Quat struct {
	W float64
	X float64
	Y float64
	Z float64
}

// Pose represents one cube transform.
type Pose struct {
	P Vec3
	Q Quat
}
