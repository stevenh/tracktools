package gpmf

//go:generate stringer -type=Type -output=types_string.go

// Type represents an encoding type.
type Type byte

// Known encoding types.
const (
	// Int8 8-bit signed integer.
	Int8 Type = 'b'

	// Uint8 8-bit unsigned integer.
	Uint8 Type = 'B'

	// StringAscii single byte 'c' style ASCII character string.
	String Type = 'c'

	// Float64 64-bit double precison (IEEE 754).
	Float64 Type = 'd'

	// Float32 32 bit float (IEEE 754).
	Float32 Type = 'f'

	// FourCC 32 bit four character key.
	FourCC Type = 'F'

	// GUID 128 bit ID (like UUID).
	GUID Type = 'G'

	// Int64 64-bit signed number.
	Int64 Type = 'j'

	// Uint64 64-bit unsigned number.
	Uint64 Type = 'J'

	// Int32 32-bit signed integer.
	Int32 Type = 'l'

	// Uint32 32-bit unsigned integer.
	Uint32 Type = 'L'

	// Q32 32-bit Q Number.
	// Q number Q15.16 - 16-bit signed integer (A) with 16-bit fixed point (B)
	// for A.B values (range -32768.0 to 32767.99998).
	Q32 Type = 'q'

	// Q64 64-bit Q Number Q15.16.
	// Q number Q31.32 - 32-bit signed integer (A) with 32-bit fixed point (B)
	// for A.B value.
	Q64 Type = 'Q'

	// Int16 16-bit signed integer.
	Int16 Type = 's'

	// Uint16 16-bit unsigned integer.
	Uint16 Type = 'S'

	// Date UTC Date and Time (yymmddhhmmss.sss).
	Date Type = 'U'

	// Complex data structure is complex.
	// Base size in bytes, data is either opaque, or the stream
	// has a TYPE structure field for the Sample.
	Complex Type = '?'

	// Compressed huffman compressed TRM payloads.
	// 4-CC <type><size><rpt> <data ...> is compressed as 4-CC
	// '#'<new size/rpt> <type><size><rpt> <compressed data ...>.
	Compressed Type = '#'

	// Nested nested metadata
	Nested Type = 0
)
