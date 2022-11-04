package gpmf

import (
	"fmt"
)

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

// floatType returns a slice of T created using f.
func floatType[T ~[]E, E any](e *Element, size int, f func([]float64) E) error {
	vals, err := floatSlice(e.Data)
	if err != nil {
		return err
	}

	var elem E
	if len(vals)%size != 0 {
		return fmt.Errorf("parse %T: invalid number of elements %d (not a multiple of %d)", elem, len(vals), size)
	}

	d := make(T, len(vals)/size)
	for i := range d {
		d[i] = f(vals[0:size])
		vals = vals[size:]
	}

	e.Data = d

	return nil
}

// floatSlice converts data to a []float64.
func floatSlice(data any) ([]float64, error) { //nolint: gocyclo,cyclop
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
