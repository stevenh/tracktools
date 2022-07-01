package gpmf

import (
	"errors"
	"fmt"
	"io"
)

// Reader is a gpmf reader.
type Reader struct {
	want       map[string]struct{}
	devices    int
	deviceID   string
	deviceName string
}

// NewReader returns a new Reader.
func NewReader() *Reader {
	return &Reader{}
}

// Read reads and returns kvl Elements from v.
func (re *Reader) Read(r io.Reader) ([]*Element, error) {
	e := NewElement(nil)

	// TODO(steve): remove this initialisation of e as
	// it's just for debugging.
	e.Header.Type = Nested
	for i, v := range []byte(KeyStream) {
		e.Header.Key[i] = v
	}

	if err := re.read(r, e); err != nil {
		return nil, err
	}

	return e.Nested, nil
}

func (re *Reader) read(r io.Reader, parent *Element) error {
	for {
		e := NewElement(parent)
		if err := e.ReadHeader(r); err != nil {
			if errors.Is(err, io.EOF) {
				// This is only place where EOF is expected.
				return nil
			}
			return fmt.Errorf("reader: read header: %w", err)
		}

		if e.Header.Nested() {
			// Nested elements.
			lr := io.LimitReader(r, e.Total)
			if err := re.read(lr, e); err != nil {
				return err
			}
		} else if err := e.ReadData(r); err != nil {
			return err
		}

		if err := e.DiscardPadding(r); err != nil {
			return err
		}

		if err := parent.Add(e); err != nil {
			return err
		}
	}
}
