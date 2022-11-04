package gpmf

import (
	"encoding/json"
	"fmt"
	"unicode"
)

// Header represents a klv header.
// https://github.com/gopro/gpmf-parser#klv-design
type Header struct {
	// Key is a 7-bit ASCII Key (FourCC).
	// https://github.com/gopro/gpmf-parser#fourcc
	Key [4]byte

	// The traditional Length is made up of three fields.
	// https://github.com/gopro/gpmf-parser#length-type-size-repeat-structure

	// Type is the type of the data elements.
	// https://github.com/gopro/gpmf-parser#type
	Type Type

	// Size is the size of the Element data.
	// This can result it multiple items of Type.
	// https://github.com/gopro/gpmf-parser#structure-size
	Size byte

	// Count is the number of Elements of Size.
	// https://github.com/gopro/gpmf-parser#repeat
	Count uint16
}

// Nested returns true if its type is Nested,
// false otherwise.
func (h Header) Nested() bool {
	return h.Type == Nested
}

// Scalable returns true if its type is scalable,
// false otherwise.
func (h Header) Scalable() bool {
	switch h.Type { //nolint: exhaustive
	case Int8, Uint8,
		Int16, Uint16,
		Int32, Uint32,
		Int64, Uint64,
		Float64, Float32,
		Q32, Q64:
		return true
	}

	return false
}

// FourCC returns the Key as string.
func (h Header) FourCC() string {
	return string(h.Key[:])
}

// MarshalJSON implements json.Marshaler.
func (h *Header) MarshalJSON() ([]byte, error) {
	type Alias Header
	return json.Marshal(&struct {
		Key  string
		Type string
		*Alias
	}{
		Key:   string(h.Key[:]),
		Type:  h.Type.String(),
		Alias: (*Alias)(h),
	})
}

// String implements stringer.
func (h Header) String() string {
	return fmt.Sprintf("key: %s type: %s size: %d count: %d",
		h.FourCC(),
		h.Type,
		h.Size,
		h.Count,
	)
}

// validate validates h.
func (h Header) validate() error {
	for i, c := range h.Key[:] {
		if c > unicode.MaxASCII {
			return fmt.Errorf("header: invalid key[%d] = 0x%02x", i, c)
		}
	}

	if h.Type == Compressed {
		return fmt.Errorf("header: key %q: compressed data not supported", h.FourCC())
	}

	return nil
}
