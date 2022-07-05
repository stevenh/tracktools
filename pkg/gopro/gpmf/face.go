package gpmf

import (
	"math"
)

const (
	faceSizeHero6 = 20
	faceSizeHero7 = 92
	faceSizeHero8 = 28
)

var (
	faceTypeDefs = map[string]byte{
		"Lffff":                   faceSizeHero6,
		"Lffffffffffffffffffffff": faceSizeHero7,
		"Lffffff":                 faceSizeHero8,
	}
)

// Face represents face detection.
type Face struct {
	ID uint32
	X  float32
	Y  float32
	W  float32
	H  float32
	// Only available on Hero 7+
	Smile float32
	// TODO(steve): complete the definition.
	// https://github.com/gopro/gpmf-parser/issues/163
}

func parseFace(e *Element) error {
	e.initMetadata()
	if e.Header.Count == 0 {
		// Nothing to do.
		return nil
	}

	if err := validateTypeDef(e, faceTypeDefs); err != nil {
		return err
	}

	count := int(e.Header.Count)
	size := int(e.Header.Size)
	d := make([]Face, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		f := Face{
			ID: byteOrder.Uint32(e.raw[j:]),
			X:  math.Float32frombits(byteOrder.Uint32(e.raw[j+4:])),
			Y:  math.Float32frombits(byteOrder.Uint32(e.raw[j+8:])),
			W:  math.Float32frombits(byteOrder.Uint32(e.raw[j+12:])),
			H:  math.Float32frombits(byteOrder.Uint32(e.raw[j+16:])),
		}
		switch size {
		case faceSizeHero7:
			f.Smile = math.Float32frombits(byteOrder.Uint32(e.raw[j+92:]))
		case faceSizeHero8:
			f.Smile = math.Float32frombits(byteOrder.Uint32(e.raw[j+24:]))
		}
		d[i] = f
	}

	e.Data = d

	return nil
}
