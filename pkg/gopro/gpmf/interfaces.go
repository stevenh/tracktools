package gpmf

import (
	"time"
)

// offseter is implemented by type which can update their offsets.
type offseter interface {
	offsets(start, end time.Duration)
}
