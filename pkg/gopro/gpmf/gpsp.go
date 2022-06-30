package gpmf

import (
	"fmt"
)

func parseGPSDoP(e *Element) error {
	v, ok := e.Data.(uint16)
	if !ok {
		return fmt.Errorf("gps dop: unexpected data type %T (expected uint16)", e.Data)
	}

	e.Data = float64(v) / 100

	return parseMetadata(e)
}
