package gpmf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/stevenh/tracktools/pkg/gopro/gpmf/geo"
)

// Dump dumps data as JSON.
func Dump(data []*Element) error {
	d := struct {
		Data []*Element
	}{Data: data}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(d); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	fmt.Println("")

	return nil
}

// StatsDumper collates stats from slice of Elements using Walk.
type StatsDumper struct {
	lat, lon float64
	counts   map[string]int
	p        *geo.Processor
}

// NewStatsDumper returns a fully initialised StatsDumper.
func NewStatsDumper() *StatsDumper {
	return &StatsDumper{
		counts: make(map[string]int),
		p:      geo.NewProcessor(),
	}
}

// Walk is WalkFunc that collects stats.
func (s *StatsDumper) Walk(e *Element) {
	if d, ok := e.Data.(GPSData); ok {
		fmt.Println("fix:", e.Metadata["gps_fix_description"])
		fmt.Println("dop:", e.Metadata["gps_dilution_of_precision"])
		for _, v := range d {
			di := s.p.Distance(s.lat, s.lon, v.Latitude, v.Longitude)
			fmt.Printf("distance: %.2f %s\n", di, v)
		}
	}

	v := reflect.ValueOf(e.Data)
	if v.Type().Kind() == reflect.Slice {
		s.counts[e.Header.FourCC()] += v.Len()
	}
}

// Results prints the stats results.
func (s *StatsDumper) Results(w io.Writer) {
	for k, v := range s.counts {
		fmt.Fprintln(w, k, "=", v)
	}
}
