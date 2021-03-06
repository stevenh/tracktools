// Code generated by "stringer -type=Type -output=types_string.go"; DO NOT EDIT.

package gpmf

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Int8-98]
	_ = x[Uint8-66]
	_ = x[String-99]
	_ = x[Float64-100]
	_ = x[Float32-102]
	_ = x[FourCC-70]
	_ = x[GUID-71]
	_ = x[Int64-106]
	_ = x[Uint64-74]
	_ = x[Int32-108]
	_ = x[Uint32-76]
	_ = x[Q32-113]
	_ = x[Q64-81]
	_ = x[Int16-115]
	_ = x[Uint16-83]
	_ = x[Date-85]
	_ = x[Complex-63]
	_ = x[Compressed-35]
	_ = x[Nested-0]
}

const _Type_name = "NestedCompressedComplexUint8FourCCGUIDUint64Uint32Q64Uint16DateInt8StringFloat64Float32Int64Int32Q32Int16"

var _Type_map = map[Type]string{
	0:   _Type_name[0:6],
	35:  _Type_name[6:16],
	63:  _Type_name[16:23],
	66:  _Type_name[23:28],
	70:  _Type_name[28:34],
	71:  _Type_name[34:38],
	74:  _Type_name[38:44],
	76:  _Type_name[44:50],
	81:  _Type_name[50:53],
	83:  _Type_name[53:59],
	85:  _Type_name[59:63],
	98:  _Type_name[63:67],
	99:  _Type_name[67:73],
	100: _Type_name[73:80],
	102: _Type_name[80:87],
	106: _Type_name[87:92],
	108: _Type_name[92:97],
	113: _Type_name[97:100],
	115: _Type_name[100:105],
}

func (i Type) String() string {
	if str, ok := _Type_map[i]; ok {
		return str
	}
	return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
}
