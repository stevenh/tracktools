package gpmf

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
