package model

import "math"

// Quat represents an orientation quaternion.
type Quat struct {
	W float64
	X float64
	Y float64
	Z float64
}

func QuatIdentity() Quat {
	return Quat{W: 1}
}

func QuatFromAxisAngle(axis Vec3, angle float64) Quat {
	n := math.Sqrt(axis.X*axis.X + axis.Y*axis.Y + axis.Z*axis.Z)
	if n == 0 {
		return QuatIdentity()
	}
	half := angle / 2
	s := math.Sin(half) / n
	return (Quat{
		W: math.Cos(half),
		X: axis.X * s,
		Y: axis.Y * s,
		Z: axis.Z * s,
	}).Normalize()
}

func (q Quat) Normalize() Quat {
	n := math.Sqrt(q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z)
	if n == 0 {
		return QuatIdentity()
	}
	return Quat{W: q.W / n, X: q.X / n, Y: q.Y / n, Z: q.Z / n}
}

func (q Quat) Conj() Quat {
	return Quat{W: q.W, X: -q.X, Y: -q.Y, Z: -q.Z}
}

func (q Quat) Mul(other Quat) Quat {
	return Quat{
		W: q.W*other.W - q.X*other.X - q.Y*other.Y - q.Z*other.Z,
		X: q.W*other.X + q.X*other.W + q.Y*other.Z - q.Z*other.Y,
		Y: q.W*other.Y - q.X*other.Z + q.Y*other.W + q.Z*other.X,
		Z: q.W*other.Z + q.X*other.Y - q.Y*other.X + q.Z*other.W,
	}
}

func (q Quat) Rotate(v Vec3) Vec3 {
	q = q.Normalize()
	p := Quat{X: v.X, Y: v.Y, Z: v.Z}
	rotated := q.Mul(p).Mul(q.Conj())
	return Vec3{X: rotated.X, Y: rotated.Y, Z: rotated.Z}
}

func (q Quat) AlmostEqual(other Quat, eps float64) bool {
	q = q.Normalize()
	other = other.Normalize()
	// q and -q represent the same orientation.
	da := math.Abs(q.W-other.W) + math.Abs(q.X-other.X) + math.Abs(q.Y-other.Y) + math.Abs(q.Z-other.Z)
	db := math.Abs(q.W+other.W) + math.Abs(q.X+other.X) + math.Abs(q.Y+other.Y) + math.Abs(q.Z+other.Z)
	if da < db {
		return da <= eps
	}
	return db <= eps
}
