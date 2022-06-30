package gpmf

import (
	"fmt"
)

// Int16_16 is a signed 16.16 fixed-point number.
//
// The integer part ranges from -32768 to 32768, inclusive. The
// fractional part has 16 bits of precision.
//
// For example, the number one-and-a-quarter is Int16_16(1<<16 + 1<<14).
type Int16_16 int32

// String returns a human-readable representation of a 16.16 fixed-point number.
//
// For example, the number one-and-a-quarter becomes "1:16384".
func (x Int16_16) String() string {
	const shift, mask = 16, 1<<16 - 1
	if x >= 0 {
		return fmt.Sprintf("%d:%05d", int32(x>>shift), int32(x&mask))
	}
	x = -x
	if x >= 0 {
		return fmt.Sprintf("-%d:%05d", int32(x>>shift), int32(x&mask))
	}
	return "-32768:00" // The minimum value is -(1<<15).
}

// Floor returns the greatest integer value less than or equal to x.
//
// Its return type is int, not Int16_16.
func (x Int16_16) Floor() int { return int((x + 0x00) >> 16) }

// Round returns the nearest integer value to x. Ties are rounded up.
//
// Its return type is int, not Int16_16.
func (x Int16_16) Round() int { return int((x + 0x8000) >> 16) }

// Ceil returns the least integer value greater than or equal to x.
//
// Its return type is int, not Int16_16.
func (x Int16_16) Ceil() int { return int((x + 0xffff) >> 16) }

// Mul returns x*y in 26.6 fixed-point arithmetic.
func (x Int16_16) Mul(y Int16_16) Int16_16 {
	return Int16_16((int64(x)*int64(y) + 1<<15) >> 16)
}

// mul (with a lower case 'm') is an alternative implementation of Int16_16.Mul
// (with an upper case 'M'). It has the same structure as the Int32_32.Mul
// implementation, but Int16_16.mul is easier to test since Go has built-in
// 64-bit integers.
func (x Int16_16) mul(y Int16_16) Int16_16 {
	const M, N = 16, 16
	lo, hi := muli32(int32(x), int32(y))
	ret := Int16_16(hi<<M | lo>>N)
	ret += Int16_16((lo >> (N - 1)) & 1) // Round to nearest, instead of rounding down.
	return ret
}

// muli32 multiplies two int32 values, returning the 64-bit signed integer
// result as two uint32 values.
//
// muli32 isn't used directly by this package, but it has the same structure as
// muli64, and muli32 is easier to test since Go has built-in 64-bit integers.
func muli32(u, v int32) (lo, hi uint32) {
	const (
		s    = 16
		mask = 1<<s - 1
	)

	u1 := uint32(u >> s)
	u0 := uint32(u & mask)
	v1 := uint32(v >> s)
	v0 := uint32(v & mask)

	w0 := u0 * v0
	t := u1*v0 + w0>>s
	w1 := t & mask
	w2 := uint32(int32(t) >> s)
	w1 += u0 * v1
	return uint32(u) * uint32(v), u1*v1 + w2 + uint32(int32(w1)>>s)
}

// Int32_32 is a signed 32.32 fixed-point number.
//
// The integer part ranges from -2147483647 to 2147483647,
// inclusive. The fractional part has 32 bits of precision.
//
// For example, the number one-and-a-quarter is Int32_32(1<<32 + 1<<30).
type Int32_32 int64

// String returns a human-readable representation of a 32.32 fixed-point
// number.
//
// For example, the number one-and-a-quarter becomes "1:1073741824".
func (x Int32_32) String() string {
	const shift, mask = 32, 1<<32 - 1
	if x >= 0 {
		return fmt.Sprintf("%d:%010d", int64(x>>shift), int64(x&mask))
	}
	x = -x
	if x >= 0 {
		return fmt.Sprintf("-%d:%010d", int64(x>>shift), int64(x&mask))
	}
	return "-2147483647:0000" // The minimum value is -(1<<31).
}

// Floor returns the greatest integer value less than or equal to x.
//
// Its return type is int, not Int32_32.
func (x Int32_32) Floor() int { return int((x + 0x000) >> 32) }

// Round returns the nearest integer value to x. Ties are rounded up.
//
// Its return type is int, not Int32_32.
func (x Int32_32) Round() int { return int((x + 0x80000000) >> 32) }

// Ceil returns the least integer value greater than or equal to x.
//
// Its return type is int, not Int32_32.
func (x Int32_32) Ceil() int { return int((x + 0xFFFFFFFF) >> 32) }

// Mul returns x*y in 32.32 fixed-point arithmetic.
func (x Int32_32) Mul(y Int32_32) Int32_32 {
	const M, N = 32, 32
	lo, hi := muli64(int64(x), int64(y))
	ret := Int32_32(hi<<M | lo>>N)
	ret += Int32_32((lo >> (N - 1)) & 1) // Round to nearest, instead of rounding down.
	return ret
}

// muli64 multiplies two int64 values, returning the 128-bit signed integer
// result as two uint64 values.
//
// This implementation is similar to $GOROOT/src/runtime/softfloat64.go's mullu
// function, which is in turn adapted from Hacker's Delight.
func muli64(u, v int64) (lo, hi uint64) {
	const (
		s    = 32
		mask = 1<<s - 1
	)

	u1 := uint64(u >> s)
	u0 := uint64(u & mask)
	v1 := uint64(v >> s)
	v0 := uint64(v & mask)

	w0 := u0 * v0
	t := u1*v0 + w0>>s
	w1 := t & mask
	w2 := uint64(int64(t) >> s)
	w1 += u0 * v1
	return uint64(u) * uint64(v), u1*v1 + w2 + uint64(int64(w1)>>s)
}
