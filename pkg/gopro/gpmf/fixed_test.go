package gpmf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var testCases = []struct {
	x      float64
	s16_16 string
	s32_32 string
	floor  int
	round  int
	ceil   int
}{{
	x:      0,
	s16_16: "0:00000",
	s32_32: "0:0000000000",
	floor:  0,
	round:  0,
	ceil:   0,
}, {
	x:      1,
	s16_16: "1:00000",
	s32_32: "1:0000000000",
	floor:  1,
	round:  1,
	ceil:   1,
}, {
	x:      1.25,
	s16_16: "1:16384",
	s32_32: "1:1073741824",
	floor:  1,
	round:  1,
	ceil:   2,
}, {
	x:      2.5,
	s16_16: "2:32768",
	s32_32: "2:2147483648",
	floor:  2,
	round:  3,
	ceil:   3,
}, {
	x:      63 / 64.0,
	s16_16: "0:64512",
	s32_32: "0:4227858432",
	floor:  0,
	round:  1,
	ceil:   1,
}, {
	x:      -0.5,
	s16_16: "-0:32768",
	s32_32: "-0:2147483648",
	floor:  -1,
	round:  +0,
	ceil:   +0,
}, {
	x:      -4.125,
	s16_16: "-4:08192",
	s32_32: "-4:0536870912",
	floor:  -5,
	round:  -4,
	ceil:   -4,
}, {
	x:      -7.75,
	s16_16: "-7:49152",
	s32_32: "-7:3221225472",
	floor:  -8,
	round:  -8,
	ceil:   -7,
}}

func TestInt16_16(t *testing.T) {
	const one = Int16_16(1 << 16)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.x), func(t *testing.T) {
			x := Int16_16(tc.x * (1 << 16))
			t.Run("String", func(t *testing.T) {
				require.Equal(t, tc.s16_16, x.String())
			})
			t.Run("Floor", func(t *testing.T) {
				require.Equal(t, tc.floor, x.Floor())
			})
			t.Run("Round", func(t *testing.T) {
				require.Equal(t, tc.round, x.Round())
			})
			t.Run("Ceil", func(t *testing.T) {
				require.Equal(t, tc.ceil, x.Ceil())
			})
			t.Run("Mul", func(t *testing.T) {
				require.Equal(t, x, x.Mul(one))
			})
			t.Run("mul", func(t *testing.T) {
				require.Equal(t, x.mul(one), x)
			})
		})
	}
}

func TestInt32_32(t *testing.T) {
	const one = Int32_32(1 << 32)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.x), func(t *testing.T) {
			x := Int32_32(tc.x * (1 << 32))
			t.Run("String", func(t *testing.T) {
				require.Equal(t, tc.s32_32, x.String())
			})
			t.Run("Floor", func(t *testing.T) {
				require.Equal(t, tc.floor, x.Floor())
			})
			t.Run("Round", func(t *testing.T) {
				require.Equal(t, tc.round, x.Round())
			})
			t.Run("Ceil", func(t *testing.T) {
				require.Equal(t, tc.ceil, x.Ceil())
			})
			t.Run("Mul", func(t *testing.T) {
				require.Equal(t, x, x.Mul(one))
			})
		})
	}
}
