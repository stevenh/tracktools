package gpmf

import (
	"fmt"
)

type number interface {
	int8 | uint8 |
		int16 | uint16 |
		int32 | uint32 |
		float64 | float32 |
		int64 | uint64 |
		Int16_16 | Int32_32
}

// Scale represents a scale definition.
type Scale []float64

// scale applies scaling to data and returns it as a slice of float64.
func scale[N number](data []N, scale Scale) []float64 {
	n := len(scale)
	r := make([]float64, len(data))
	for i, v := range data {
		r[i] = float64(v) / scale[i%n]
	}
	return r
}

// parseScale parses scale data.
func parseScale(e *Element) error {
	d, err := floatSlice(e.Data)
	if err != nil {
		return err
	}

	s := Scale(d)
	e.Data = s
	e.parent.scale = s

	return nil
}

// toFloatSlice converts data to a float64 slice.
func toFloatSlice[N number](data []N) []float64 {
	r := make(Scale, len(data))
	for i, v := range data {
		r[i] = float64(v)
	}
	return r
}

// toFloatSliceSingle converts data to a float64 slice.
func toFloatSliceSingle[N number](data N) []float64 {
	return Scale([]float64{float64(data)})
}

// floatSlice converts data to a []float64.
func floatSlice(data any) ([]float64, error) { // nolint: gocyclo,cyclop
	switch v := data.(type) {
	case []int8:
		return toFloatSlice(v), nil
	case int8:
		return toFloatSliceSingle(v), nil
	case []uint8:
		return toFloatSlice(v), nil
	case uint8:
		return toFloatSliceSingle(v), nil
	case []int16:
		return toFloatSlice(v), nil
	case int16:
		return toFloatSliceSingle(v), nil
	case []uint16:
		return toFloatSlice(v), nil
	case uint16:
		return toFloatSliceSingle(v), nil
	case []int32:
		return toFloatSlice(v), nil
	case int32:
		return toFloatSliceSingle(v), nil
	case []uint32:
		return toFloatSlice(v), nil
	case uint32:
		return toFloatSliceSingle(v), nil
	case []int64:
		return toFloatSlice(v), nil
	case int64:
		return toFloatSliceSingle(v), nil
	case []uint64:
		return toFloatSlice(v), nil
	case uint64:
		return toFloatSliceSingle(v), nil
	case []float32:
		return toFloatSlice(v), nil
	case float32:
		return toFloatSliceSingle(v), nil
	case []float64:
		return toFloatSlice(v), nil
	case float64:
		return toFloatSliceSingle(v), nil
	// TODO(steve) is this correct for Q types?
	case []Int16_16:
		return toFloatSlice(v), nil
	case Int16_16:
		return toFloatSliceSingle(v), nil
	case []Int32_32:
		return toFloatSlice(v), nil
	case Int32_32:
		return toFloatSliceSingle(v), nil
	default:
		return nil, fmt.Errorf("element: to scale invalid element type %T", v)
	}
}
