package model

import (
	"math"
	"testing"
)

func TestVec3Operations(t *testing.T) {
	tests := []struct {
		name     string
		op       func() Vec3
		expected Vec3
	}{
		{
			name: "Add",
			op: func() Vec3 {
				return Vec3{X: 1, Y: 2, Z: 3}.Add(Vec3{X: 4, Y: 5, Z: 6})
			},
			expected: Vec3{X: 5, Y: 7, Z: 9},
		},
		{
			name: "Sub",
			op: func() Vec3 {
				return Vec3{X: 5, Y: 7, Z: 9}.Sub(Vec3{X: 1, Y: 2, Z: 3})
			},
			expected: Vec3{X: 4, Y: 5, Z: 6},
		},
		{
			name: "Scale",
			op: func() Vec3 {
				return Vec3{X: 1, Y: 2, Z: 3}.Scale(2)
			},
			expected: Vec3{X: 2, Y: 4, Z: 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.op()
			if !got.AlmostEqual(tt.expected, 1e-9) {
				t.Errorf("%s got %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestVec3Distance(t *testing.T) {
	tests := []struct {
		name     string
		v1, v2   Vec3
		expected float64
	}{
		{
			name:     "Same point",
			v1:       Vec3{X: 1, Y: 2, Z: 3},
			v2:       Vec3{X: 1, Y: 2, Z: 3},
			expected: 0,
		},
		{
			name:     "X axis",
			v1:       Vec3{X: 0, Y: 0, Z: 0},
			v2:       Vec3{X: 3, Y: 0, Z: 0},
			expected: 3,
		},
		{
			name:     "Diagonal",
			v1:       Vec3{X: 0, Y: 0, Z: 0},
			v2:       Vec3{X: 1, Y: 1, Z: 1},
			expected: math.Sqrt(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v1.Distance(tt.v2)
			if math.Abs(got-tt.expected) > 1e-9 {
				t.Errorf("%s got %f, want %f", tt.name, got, tt.expected)
			}
		})
	}
}

func TestPoseAlmostEqual(t *testing.T) {
	p1 := Pose{P: Vec3{X: 1, Y: 0, Z: 0}, Q: QuatIdentity()}
	p2 := Pose{P: Vec3{X: 1.000000001, Y: 0, Z: 0}, Q: QuatIdentity()}
	if !p1.AlmostEqual(p2, 1e-8) {
		t.Errorf("expected poses to be almost equal")
	}
	if p1.AlmostEqual(p2, 1e-10) {
		t.Errorf("expected poses NOT to be almost equal")
	}
}

func TestAxis(t *testing.T) {
	tests := []struct {
		a    Axis
		s    string
		unit Vec3
	}{
		{AxisX, "X", Vec3{X: 1}},
		{AxisY, "Y", Vec3{Y: 1}},
		{AxisZ, "Z", Vec3{Z: 1}},
	}
	for _, tt := range tests {
		if tt.a.String() != tt.s {
			t.Errorf("Axis %v String() = %v, want %v", tt.a, tt.a.String(), tt.s)
		}
		if !tt.a.UnitVector().AlmostEqual(tt.unit, 1e-9) {
			t.Errorf("Axis %v UnitVector() = %v, want %v", tt.a, tt.a.UnitVector(), tt.unit)
		}
		got, _ := AxisFromString(tt.s)
		if got != tt.a {
			t.Errorf("AxisFromString(%v) = %v, want %v", tt.s, got, tt.a)
		}
	}
	if _, err := AxisFromString("W"); err == nil {
		t.Errorf("expected error for invalid axis string")
	}
}
