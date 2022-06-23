package laptimer

import (
	"bufio"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// EncoderOpt represents an Encoder option.
type EncoderOpt func(*Encoder) error

// Compress enables gzip compression.
func Compress() EncoderOpt {
	return func(e *Encoder) error {
		zw := gzip.NewWriter(e.w)
		e.w = zw
		e.c = zw

		return nil
	}
}

// Encoder writes Harry's LapTimer xml data files.
type Encoder struct {
	compress bool
	w        io.Writer
	c        io.Closer
}

// NewEncoder returns a fully initialised encoder which writes its
// output to w.
func NewEncoder(w io.Writer, options ...EncoderOpt) (*Encoder, error) {
	e := &Encoder{w: w}
	for _, f := range options {
		if err := f(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}

// Encode encodes v as xml to the encoders output stream.
func (e *Encoder) Encode(v any) error {
	if _, err := e.w.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("encode write header: %w", err)
	}

	r, w := io.Pipe()
	enc := xml.NewEncoder(w)
	enc.Indent("", "\t")

	errs := e.run(r)
	if err := enc.Encode(v); err != nil {
		w.Close() // nolint: errcheck
		return fmt.Errorf("encode: %w", err)
	}

	// Close w so it triggers io.EOF in filter.
	if err := w.Close(); err != nil {
		return fmt.Errorf("encode close: %w", err)
	}

	if err := <-errs; err != nil {
		return err
	}

	if e.c == nil {
		// No compression so all done.
		return nil
	}

	// Close gzip Writer to ensure all data is flushed.
	return e.c.Close()
}

// run filters the output from r to our output channel w.
func (e *Encoder) run(r io.Reader) <-chan error {
	errs := make(chan error, 1)
	go func() {
		errs <- e.filter(r)
	}()

	return errs
}

// filter filters the data from r and writes it w.
func (e *Encoder) filter(r io.Reader) error {
	// rep updates our output so it matches LapTimer behaviour.
	rep := strings.NewReplacer(
		// golang uses #34 and #39 as they are shorter.
		"&#34;", "&quote;",
		"&#39;", "&apos;",
		// golang escapes newlines and tabs.
		"&#xA;", "\n",
		"&#x9;", "\t",
	)

	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		// We check err later as we can get data even with an error
		// in particular io.EOF.
		if _, werr := e.w.Write([]byte(rep.Replace(line))); werr != nil {
			return fmt.Errorf("encoder write: %w", werr)
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				// Last line.
				return nil
			}

			return fmt.Errorf("encoder read: %w", err)
		}
	}
}
