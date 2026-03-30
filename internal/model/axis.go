package model

import "fmt"

// Axis identifies a local hinge axis.
type Axis uint8

const (
	AxisX Axis = iota
	AxisY
	AxisZ
)

func AxisFromString(v string) (Axis, error) {
	switch v {
	case "X":
		return AxisX, nil
	case "Y":
		return AxisY, nil
	case "Z":
		return AxisZ, nil
	default:
		return 0, fmt.Errorf("axisA must be X, Y, or Z")
	}
}

func (a Axis) String() string {
	switch a {
	case AxisX:
		return "X"
	case AxisY:
		return "Y"
	case AxisZ:
		return "Z"
	default:
		return "?"
	}
}
