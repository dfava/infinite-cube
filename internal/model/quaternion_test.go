package model

import (
	"math"
	"testing"
)

func TestQuatOperations(t *testing.T) {
	tests := []struct {
		name     string
		op       func() Quat
		expected Quat
	}{
		{
			name: "Identity",
			op: func() Quat {
				return QuatIdentity()
			},
			expected: Quat{W: 1, X: 0, Y: 0, Z: 0},
		},
		{
			name: "FromAxisAngle X 90",
			op: func() Quat {
				return QuatFromAxisAngle(Vec3{X: 1, Y: 0, Z: 0}, math.Pi/2)
			},
			expected: Quat{W: math.Sqrt(0.5), X: math.Sqrt(0.5), Y: 0, Z: 0},
		},
		{
			name: "Normalize",
			op: func() Quat {
				return Quat{W: 2, X: 0, Y: 0, Z: 0}.Normalize()
			},
			expected: Quat{W: 1, X: 0, Y: 0, Z: 0},
		},
		{
			name: "Conjugate",
			op: func() Quat {
				return Quat{W: 1, X: 2, Y: 3, Z: 4}.Conj()
			},
			expected: Quat{W: 1, X: -2, Y: -3, Z: -4},
		},
		{
			name: "Multiply Identities",
			op: func() Quat {
				return QuatIdentity().Mul(QuatIdentity())
			},
			expected: QuatIdentity(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.op()
			if !got.AlmostEqual(tt.expected, 1e-7) {
				t.Errorf("%s got %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestQuatRotate(t *testing.T) {
	tests := []struct {
		name     string
		q        Quat
		v        Vec3
		expected Vec3
	}{
		{
			name:     "Identity rotate",
			q:        QuatIdentity(),
			v:        Vec3{X: 1, Y: 2, Z: 3},
			expected: Vec3{X: 1, Y: 2, Z: 3},
		},
		{
			name:     "Rotate X 90 around Z axis",
			q:        QuatFromAxisAngle(Vec3{X: 0, Y: 0, Z: 1}, math.Pi/2),
			v:        Vec3{X: 1, Y: 0, Z: 0},
			expected: Vec3{X: 0, Y: 1, Z: 0},
		},
		{
			name:     "Rotate X 180 around Y axis",
			q:        QuatFromAxisAngle(Vec3{X: 0, Y: 1, Z: 0}, math.Pi),
			v:        Vec3{X: 1, Y: 0, Z: 0},
			expected: Vec3{X: -1, Y: 0, Z: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.q.Rotate(tt.v)
			if !got.AlmostEqual(tt.expected, 1e-7) {
				t.Errorf("%s got %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestQuatAlmostEqual(t *testing.T) {
	q1 := Quat{W: 1, X: 0, Y: 0, Z: 0}
	q2 := Quat{W: -1, X: 0, Y: 0, Z: 0}
	if !q1.AlmostEqual(q2, 1e-9) {
		t.Errorf("expected q1 and -q1 to be almost equal")
	}
}
