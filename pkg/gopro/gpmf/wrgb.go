package gpmf

// WhiteBalanceRGB represents a white balance with RGB channels.
type WhiteBalanceRGB struct {
	Red   float64
	Green float64
	Blue  float64
}

func parseWhiteBalanceRGB(e *Element) error {
	e.metadata()
	return floatType(e, 3, func(vals []float64) WhiteBalanceRGB {
		return WhiteBalanceRGB{
			Red:   vals[0],
			Green: vals[1],
			Blue:  vals[2],
		}
	})
}
