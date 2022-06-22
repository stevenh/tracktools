package convert

import (
	"math"

	"github.com/stevenh/tracktools/pkg/laptimer"
)

func round2dp(v float64) laptimer.Float2dp {
	return laptimer.Float2dp(math.Round(v*100) / 100)
}

func round1dp(v float64) laptimer.Float1dp {
	return laptimer.Float1dp(math.Round(v*10) / 10)
}

func round0dp(v float64) laptimer.Float0dp {
	return laptimer.Float0dp(math.Round(v))
}

func roundint(v float64) int {
	return int(math.Round(v))
}
