package gpmf

import (
	"time"
)

// offseter is implemented by type which can update their offsets.
type offseter interface {
	offsets(start, end time.Duration)
}

// number represents a number we can scale.
type number interface {
	int8 | uint8 |
		int16 | uint16 |
		int32 | uint32 |
		float64 | float32 |
		int64 | uint64 |
		Int16_16 | Int32_32
}
