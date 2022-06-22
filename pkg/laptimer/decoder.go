package laptimer

import (
	"encoding/xml"
	"fmt"
	"io"

	"golang.org/x/text/encoding/ianaindex"
)

// Decoder reads Harry's LapTimer xml data files.
type Decoder struct {
	r io.Reader
}

// NewDecoder returns a fully initialised encoder which reads
// data from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode decodes data from the stream into v.
func (d *Decoder) Decode(v any) error {
	dec := xml.NewDecoder(d.r)
	dec.CharsetReader = func(charset string, r io.Reader) (io.Reader, error) {
		enc, err := ianaindex.IANA.Encoding(charset)
		if err != nil {
			return nil, fmt.Errorf("charset %s: %w", charset, err)
		}
		if enc == nil {
			// Assume it's compatible with (a subset of) UTF-8 encoding
			// Bug: https://github.com/golang/go/issues/19421
			return r, nil
		}
		return enc.NewDecoder().Reader(r), nil
	}
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
