package gpmf

import (
	"math"
)

const (
	faceTypeDef = "Lffff"
)

// Face represents camera pointing direction.
type Face struct {
	ID uint32
	X  float32
	Y  float32
	W  float32
	H  float32
}

func parseFace(e *Element) error {
	e.initMetadata()
	if e.Header.Count == 0 {
		// Nothing to do.
		return nil
	}

	if err := validateTypeDef(e, 20, faceTypeDef); err != nil {
		return err
	}

	count := int(e.Header.Count)
	d := make([]Face, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+20 {
		f := Face{
			ID: byteOrder.Uint32(e.raw[j:]),
			X:  math.Float32frombits(byteOrder.Uint32(e.raw[j+4:])),
			Y:  math.Float32frombits(byteOrder.Uint32(e.raw[j+8:])),
			W:  math.Float32frombits(byteOrder.Uint32(e.raw[j+12:])),
			H:  math.Float32frombits(byteOrder.Uint32(e.raw[j+16:])),
		}
		d[i] = f
	}

	e.Data = d

	return nil
}
