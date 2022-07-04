package gpmf

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strings"
	"time"
)

const (
	// alignment is the alignment in bytes of gpfm data (32bit).
	alignment = 4

	// dataFormat for Date.
	dateFormat = "060102150405.000"
)

var (
	// standardUnitsFix fixes standard units by replacing ASCII codes
	// with the correct UTF-8 characters.
	standardUnitsFix = strings.NewReplacer(
		"\xB0", "°",
		"\xB2", "²",
		"\xB3", "³",
		"\xB5", "µ",
	)

	// byteOrder is the byte order for gpfm data.
	byteOrder = binary.BigEndian
)

// Element represents a klv element.
type Element struct {
	// Header is the header of the klv.
	Header Header

	// Total is the total data size including padding calculated by ReadHeader.
	Total int64

	// Metadata represents element metadata as applied by sticky items.
	Metadata map[string]interface{} `json:",omitempty"`

	// Nested represents any nested Key Length Values.
	Nested []*Element `json:",omitempty"`

	// Data is the data formatted and scaled if needed data.
	Data any `json:",omitempty"`

	// raw is raw data.
	raw []byte

	// level is the level of this element in stream.
	level int

	// elements is a map of elements by key, duplicates will be lost.
	elements map[string]*Element

	// scale is the last unused scale element.
	scale []float64

	// typeDef is the format for custom types.
	typeDef string

	// size is the total data size excluding padding calculated by ReadHeader.
	size int64

	// padding is the amount of padding calculated by ReaderHeader.
	padding int64

	// parent is the parent element.
	parent *Element
}

// NewElement returns a Element with it's internal structures initialised.
func NewElement(parent *Element) *Element {
	e := &Element{
		elements: make(map[string]*Element),
		Metadata: make(map[string]interface{}),
		parent:   parent,
	}

	if parent != nil {
		e.level = parent.level + 1
	}

	return e
}

// Add adds a nested element other to e.
func (e *Element) Add(other *Element) error {
	other.parent = e
	if err := other.format(e); err != nil {
		return err
	}

	e.Nested = append(e.Nested, other)

	return nil
}

// MarshalJSON implements json.Marshaler.
func (e *Element) MarshalJSON() ([]byte, error) {
	type Alias Element
	v := &struct {
		Level int
		*Alias
	}{
		Level: e.level,
		Alias: (*Alias)(e),
	}
	return json.Marshal(v)
}

// ReadHeader reads the header details from r and calculates its sizing.
func (e *Element) ReadHeader(r io.Reader) error {
	if err := binary.Read(r, byteOrder, &e.Header); err != nil {
		return fmt.Errorf("element: read header: %w", err)
	}

	if err := e.Header.validate(); err != nil {
		return err
	}

	e.size = int64(e.Header.Size) * int64(e.Header.Count)
	e.Total = int64(math.Ceil(float64(e.size)/alignment) * alignment)
	e.padding = e.Total - e.size

	return nil
}

// ReadData reads the total size from r storing it in e.raw.
func (e *Element) ReadData(r io.Reader) error {
	e.raw = make([]byte, e.size)
	if n, err := io.ReadFull(r, e.raw); err != nil {
		return fmt.Errorf("element: read data %d of %d: %w", n, e.size, err)
	}
	return nil
}

// DiscardPadding reads and discards the padding from r leaving it ready for the next Element.
func (e *Element) DiscardPadding(r io.Reader) error {
	if _, err := io.CopyN(ioutil.Discard, r, e.padding); err != nil {
		return fmt.Errorf("element: discard padding: %w", err)
	}

	return nil
}

// String implements Stringer.
func (e Element) String() string {
	return fmt.Sprintf("Key: %s Type: %s Size: %d Count: %d Total: %d Padding: %d",
		e.Header.FourCC(),
		e.Header.Type,
		e.Header.Size,
		e.Header.Count,
		e.Total,
		e.padding,
	)
}

// MetadataByKey returns the elements metadata for key.
func (e *Element) MetadataByKey(key string) (any, bool) {
	v, ok := e.Metadata[keyNames[key]]
	return v, ok
}

// initMetadata sets the metadata on e from its parents.
func (e *Element) initMetadata() {
	e.Metadata = e.parent.Metadata
	for v := e.parent; v.parent != nil; v = v.parent {
		for k, v := range v.Metadata {
			if _, ok := e.Metadata[k]; !ok {
				e.Metadata[k] = v
			}
		}
	}
}

// format returns the element data formatted according
// to its Header information.
func (e *Element) format(parent *Element) error {
	if err := e.formatBasic(parent); err != nil {
		return err
	}

	if parent.scale != nil {
		// Apply scaling.
		s, err := floatSlice(e.Data)
		if err != nil {
			return err
		}

		e.Data = scale(s, parent.scale)
		parent.scale = nil
	}

	if f := keyParsers[e.Header.FourCC()]; f != nil {
		if err := f(e); err != nil {
			return err
		}
	}

	return nil
}

// formatBasic stores the.Data version according
// to its Header information.
func (e *Element) formatBasic(parent *Element) error { // nolint: cyclop
	// Ensure raw data is valid ReadData will have ensured this.
	// TODO(steve): remove?
	if e.Header.Type != Nested && int64(len(e.raw)) != e.Total-e.padding {
		return fmt.Errorf("element: %s: unexpected raw len %d != %d", e, len(e.raw), e.Total-e.padding)
	}

	switch e.Header.Type {
	case Int8:
		return e.formatInt8s()
	case Uint8:
		return e.formatUint8s()
	case String, FourCC, GUID:
		return e.formatStrings()
	case Int16:
		return e.formatInt16s(parent)
	case Uint16:
		return e.formatUint16s()
	case Float32:
		return e.formatFloat32s()
	case Int32:
		return e.formatInt32s()
	case Uint32:
		return e.formatUint32s()
	case Q32:
		return e.formatInt16_16s()
	case Float64:
		return e.formatFloat64s()
	case Int64:
		return e.formatInt64s()
	case Uint64:
		return e.formatUint64s()
	case Q64:
		return e.formatInt32_32s()
	case Date:
		return e.formatDates()
	case Complex:
		// TODO(steve): support fully.
		return nil
	case Compressed:
		return fmt.Errorf("element: type %s not supported", e.Header.Type)
	case Nested:
		// Nested doesn't have raw data.
		return nil
	default:
		return fmt.Errorf("element: type %s unknown", e.Header.Type)
	}
}

// toString returns the buf converted to a string
// with replacements done and any null termination
// removed.
func (e Element) toString(buf []byte) string {
	// Remove null termination.
	ret := strings.TrimRight(string(buf), "\x00")

	if e.Header.FourCC() == KeyStandardUnits {
		return standardUnitsFix.Replace(ret)
	}

	return ret
}

func (e *Element) formatInt8s() error {
	if e.Header.Count == 1 {
		e.Data = int8(e.raw[0])
		return nil
	}

	d := make([]int8, e.Header.Count)
	for i, v := range e.raw {
		d[i] = int8(v)
	}

	e.Data = d

	return nil
}

func (e *Element) formatUint8s() error {
	if e.Header.Count == 1 {
		e.Data = uint8(e.raw[0]) // nolint: unconvert
		return nil
	}

	e.Data = []uint8(e.raw) // nolint: unconvert

	return nil
}

func (e *Element) formatStrings() error {
	if e.Header.Size == 1 || e.Header.Count == 1 {
		// String of chars.
		e.Data = e.toString(e.raw)
		return nil
	}

	// Multiple strings.
	d := make([]string, e.Header.Count)
	size := int(e.Header.Size)
	for i, j := 0, 0; i < int(e.Header.Count); i, j = i+1, j+size {
		d[i] = e.toString(e.raw[j : j+size])
	}

	e.Data = d

	return nil
}

func (e *Element) formatInt16s(parent *Element) error {
	size := 2
	count := int(e.size) / size
	if count == 1 {
		e.Data = int16(byteOrder.Uint16(e.raw))
		return nil
	}

	d := make([]int16, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		d[i] = int16(byteOrder.Uint16(e.raw[j:]))
	}

	e.Data = d

	return nil
}

func (e *Element) formatUint16s() error {
	size := 2
	count := int(e.size) / size
	if count == 1 {
		e.Data = byteOrder.Uint16(e.raw)
		return nil
	}

	d := make([]uint16, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		d[i] = byteOrder.Uint16(e.raw[j:])
	}

	e.Data = d

	return nil
}

func (e *Element) formatFloat32s() error {
	size := 4
	count := int(e.size) / size
	if count == 1 {
		e.Data = math.Float32frombits(byteOrder.Uint32(e.raw))
		return nil
	}

	d := make([]float32, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		d[i] = math.Float32frombits(byteOrder.Uint32(e.raw[j:]))
	}

	e.Data = d

	return nil
}

func (e *Element) formatInt32s() error {
	size := 4
	count := int(e.size) / size
	if count == 1 {
		e.Data = int32(byteOrder.Uint32(e.raw))
		return nil
	}

	d := make([]int32, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		d[i] = int32(byteOrder.Uint32(e.raw[j:]))
	}

	e.Data = d

	return nil
}

func (e *Element) formatUint32s() error {
	size := 4
	count := int(e.size) / size
	if count == 1 {
		e.Data = byteOrder.Uint32(e.raw)
		return nil
	}

	d := make([]uint32, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		d[i] = byteOrder.Uint32(e.raw[j:])
	}

	e.Data = d

	return nil
}

func (e *Element) formatInt16_16s() error {
	size := 8
	count := int(e.size) / size
	if count == 1 {
		e.Data = Int16_16(byteOrder.Uint32(e.raw))
		return nil
	}

	d := make([]Int16_16, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		d[i] = Int16_16(byteOrder.Uint32(e.raw[j:]))
	}

	e.Data = d

	return nil
}

func (e *Element) formatFloat64s() error {
	size := 8
	count := int(e.size) / size
	if count == 1 {
		e.Data = math.Float64frombits(byteOrder.Uint64(e.raw))
		return nil
	}

	d := make([]float64, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		d[i] = math.Float64frombits(byteOrder.Uint64(e.raw[j:]))
	}

	e.Data = d

	return nil
}

func (e *Element) formatInt64s() error {
	size := 8
	count := int(e.size) / size
	if count == 1 {
		e.Data = int64(byteOrder.Uint64(e.raw))
		return nil
	}

	d := make([]int64, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		d[i] = int64(byteOrder.Uint64(e.raw[j:]))
	}

	e.Data = d

	return nil
}

func (e *Element) formatUint64s() error {
	size := 8
	count := int(e.size) / size
	if count == 1 {
		e.Data = byteOrder.Uint64(e.raw)
		return nil
	}

	d := make([]uint64, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		d[i] = byteOrder.Uint64(e.raw[j:])
	}

	e.Data = d

	return nil
}

func (e *Element) formatInt32_32s() error {
	size := 8
	count := int(e.size) / size
	if count == 1 {
		e.Data = Int32_32(byteOrder.Uint64(e.raw))
		return nil
	}

	d := make([]Int32_32, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		d[i] = Int32_32(byteOrder.Uint64(e.raw[j:]))
	}

	e.Data = d

	return nil
}

func (e *Element) formatDates() error {
	size := 16
	count := int(e.size) / size
	if count == 1 {
		date := string(e.raw)
		t, err := time.Parse(dateFormat, date)
		if err != nil {
			return fmt.Errorf("element: parse date %q: %w", date, err)
		}

		e.Data = t

		return nil
	}

	d := make([]time.Time, count)
	for i, j := 0, 0; i < count; i, j = i+1, j+size {
		date := string(e.raw[j : j+size])
		t, err := time.Parse(dateFormat, date)
		if err != nil {
			return fmt.Errorf("element: parse date %q: %w", date, err)
		}

		d[i] = t
	}

	e.Data = d

	return nil
}

func (e *Element) friendlyName() string {
	return friendlyName(e.Header.FourCC())
}
