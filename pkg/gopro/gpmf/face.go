package gpmf

import (
	"math"
)

// Face type sizes.
const (
	faceSizeHero6  = 20
	faceSizeHero7  = 92
	faceSizeHero8  = 28
	faceSizeHero10 = 14
)

// Face type definitions.
const (
	faceDefHero6  = "Lffff"
	faceDefHero7  = "Lffffffffffffffffffffff"
	faceDefHero8  = "Lffffff"
	faceDefHero10 = "BBSSSSSBB"
)

var (
	// faceTypeDefs is a map of all known face type definitions.
	faceTypeDefs = map[string]byte{
		faceDefHero6:  faceSizeHero6,
		faceDefHero7:  faceSizeHero7,
		faceDefHero8:  faceSizeHero8,
		faceDefHero10: faceSizeHero10,
	}
)

// face is an interface that represents all face types.
type face interface {
	Face6 | Face7 | Face8 | Face10
}

// Face6 represents face detection for Hero 6.
type Face6 struct {
	// ID is the unique ID of the face.
	ID uint32

	// X is the starting coordinate on the horizontal axis of the face bounding box.
	X float32

	// Y is the starting coordinate on the vertical axis of the face bounding box.
	Y float32

	// Width is the width of the face bounding box.
	Width float32

	// Height is the hight of the face bounding box.
	Height float32
}

// Face7 represents a detected face for Hero 7.
type Face7 struct {
	Face6

	// Smile is the percentage confidence that the face is smiling.
	Smile float32
}

// Face8 represents a detected face for Hero 7.
type Face8 struct {
	Face7

	// Confidence is the percentage confidence that the face contains a face.
	Confidence float32
}

// Face10 represents a detected face for Hero 10+.
type Face10 struct {
	// Version is the version of the face definition.
	Version byte

	// Confidence is the percentage confidence that the face contains a face.
	Confidence byte

	// ID is the unique ID of the face.
	ID uint16

	// X is the starting coordinate on the horizontal axis of the face bounding box.
	X uint16

	// Y is the starting coordinate on the vertical axis of the face bounding box.
	Y uint16

	// Width is the width of the face bounding box.
	Width uint16

	// Height is the hight of the face bounding box.
	Height uint16

	// Smile is the percentage confidence that the face is smiling.
	Smile byte

	// Blink is the percentage confidence that the face is blinking.
	Blink byte
}

// parseFace parses a face detecting which format to be used.
func parseFace(e *Element) error {
	e.initMetadata()
	if e.Header.Count == 0 {
		// Nothing to do.
		return nil
	}

	def, err := validateTypeDef(e, faceTypeDefs)
	if err != nil {
		return err
	}

	switch def {
	case faceDefHero6:
		parseFaces(e, parseFace6)
	case faceDefHero7:
		parseFaces(e, parseFace7)
	case faceDefHero8:
		parseFaces(e, parseFace8)
	case faceDefHero10:
		parseFaces(e, parseFace10)
	}

	return nil
}

// parseFaces parses a set of faces using fn.
func parseFaces[T face](e *Element, fn func(*T, []byte)) {
	count := int(e.Header.Count)
	size := int(e.Header.Size)
	d := make([]T, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		var f T
		fn(&f, e.raw[j:])
		d[i] = f
	}

	e.Data = d
}

// parseFace6 parses a Hero 6 face.
func parseFace6(f *Face6, raw []byte) {
	f.ID = byteOrder.Uint32(raw)
	f.X = math.Float32frombits(byteOrder.Uint32(raw[4:]))
	f.Y = math.Float32frombits(byteOrder.Uint32(raw[8:]))
	f.Width = math.Float32frombits(byteOrder.Uint32(raw[12:]))
	f.Height = math.Float32frombits(byteOrder.Uint32(raw[16:]))
}

// parseFace7 parses a Hero 7 face.
func parseFace7(f *Face7, raw []byte) {
	parseFace6(&f.Face6, raw)
	f.Smile = math.Float32frombits(byteOrder.Uint32(raw[92:]))
}

// parseFace8 parses a Hero 8+ face.
func parseFace8(f *Face8, raw []byte) {
	parseFace6(&f.Face6, raw)
	f.Confidence = math.Float32frombits(byteOrder.Uint32(raw[20:]))
	f.Smile = math.Float32frombits(byteOrder.Uint32(raw[24:]))
}

// parseFace10 parses a Hero 10+ face.
func parseFace10(f *Face10, raw []byte) {
	f.Version = raw[0]
	f.Confidence = raw[1]
	f.ID = byteOrder.Uint16(raw[3:])
	f.X = byteOrder.Uint16(raw[5:])
	f.Y = byteOrder.Uint16(raw[7:])
	f.Width = byteOrder.Uint16(raw[9:])
	f.Height = byteOrder.Uint16(raw[11:])
	f.Smile = raw[13]
	f.Blink = raw[14]
}
