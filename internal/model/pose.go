package model

// Pose represents one cube transform.
type Pose struct {
	P Vec3
	Q Quat
}

func (p Pose) AlmostEqual(other Pose, eps float64) bool {
	if p.P.Distance(other.P) > eps {
		return false
	}
	return p.Q.AlmostEqual(other.Q, eps)
}
