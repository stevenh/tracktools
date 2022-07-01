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
	lo, hi := muli[Int16_16, uint32](M, x, y)
	ret := Int16_16(hi<<M | lo>>N)
	ret += Int16_16((lo >> (N - 1)) & 1) // Round to nearest, instead of rounding down.
	return ret
}

// muli multiplies two integer values, returning the signed integer
// result as two signed values.
func muli[I ~int32 | ~int64, O uint32 | uint64](s int, u, v I) (lo, hi O) { // nolint: ireturn
	var mask I = 1<<s - 1

	u1 := O(u >> s)
	u0 := O(u & mask)
	v1 := O(v >> s)
	v0 := O(v & mask)

	w0 := u0 * v0
	t := u1*v0 + w0>>s
	w1 := t & O(mask)
	w2 := O(I(t) >> s)
	w1 += u0 * v1
	return O(u) * O(v), u1*v1 + w2 + O(I(w1)>>s)
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
	lo, hi := muli[Int32_32, uint64](M, x, y)
	ret := Int32_32(hi<<M | lo>>N)
	ret += Int32_32((lo >> (N - 1)) & 1) // Round to nearest, instead of rounding down.
	return ret
}
