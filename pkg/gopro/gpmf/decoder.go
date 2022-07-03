package gpmf

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/edgeware/mp4ff/mp4"
	"github.com/stevenh/tracktools/pkg/gopro/gpmf/geo"
)

const (
	handlerType = "meta"
	handlerName = "GoPro MET"
)

// Decoder is a GoPro mp4 metadata decoder.
type Decoder struct {
	reader *Reader
}

// NewDecoder returns a new Decoder.
func NewDecoder() *Decoder {
	return &Decoder{
		reader: NewReader(),
	}
}

// Decode decodes metadata from the mp4 stream in rs.
func (d *Decoder) Decode(rs io.ReadSeeker) ([]*Element, error) {
	f, err := mp4.DecodeFile(rs)
	if err != nil {
		return nil, fmt.Errorf("decode: mp4 %w", err)
	}

	for i, trak := range f.Moov.Traks {
		if trak.Mdia.Hdlr.HandlerType != handlerType {
			continue
		} else if !strings.Contains(trak.Mdia.Hdlr.Name, handlerName) {
			continue
		}

		stbl := trak.Mdia.Minf.Stbl

		units := time.Second / time.Duration(trak.Mdia.Mdhd.Timescale)
		data, err := d.decodeTrak(rs, stbl, units)
		if err != nil {
			return nil, fmt.Errorf("decode: trak %d: %w", i, err)
		}

		return data, nil
	}

	return nil, fmt.Errorf("no metadata %q trak found", handlerName)
}

// chunkOffsets returns the chunk offsets for stbl.
func (d *Decoder) chunkOffsets(stbl *mp4.StblBox) ([]uint64, error) {
	switch {
	case stbl.Stco != nil:
		res := make([]uint64, len(stbl.Stco.ChunkOffset))
		for i := range stbl.Stco.ChunkOffset {
			res[i] = uint64(stbl.Stco.ChunkOffset[i])
		}

		return res, nil
	case stbl.Co64 != nil:
		return stbl.Co64.ChunkOffset, nil
	default:
		return nil, errors.New("no stco or co64 available")
	}
}

func (d *Decoder) readChunk(rs io.ReadSeeker,
	offset, size int64,
	startDec, endDec uint64,
	units time.Duration,
) ([]*Element, error) {
	if _, err := rs.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek: %w", err)
	}

	data, err := d.reader.Read(io.LimitReader(rs, size))
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	// TODO(steve): remove
	d.dumpStats(data, startDec, endDec, units)
	return data, nil

	if err := dump(data); err != nil {
		return nil, err
	}

	return data, nil
}

// TODO(steve): remove
func (d *Decoder) dumpStats(data []*Element, startDec, endDec uint64, units time.Duration) {
	start := time.Duration(startDec) * units
	end := time.Duration(endDec) * units
	dur := end - start

	fmt.Println("chunk start:", start, "end:", end, "dur:", dur)

	p := geo.NewProcessor()

	counts := make(map[string]int)
	for _, e := range data {
		counts[e.Header.FourCC()]++
		for _, e := range e.Nested {

			key := e.Header.FourCC()
			counts[key]++
			for _, e := range e.Nested {
				if e.Data == nil {
					continue
				}
				v := reflect.ValueOf(e.Data)
				t := v.Type()
				if t.Kind() == reflect.Slice {
					counts[e.Header.FourCC()] += v.Len()
					if s, ok := e.Data.([]GPS); ok {
						fmt.Println("fix:", e.Metadata["gps_fix_description"])
						fmt.Println("dop:", e.Metadata["gps_dilution_of_precision"])
						inc := dur / time.Duration(len(s))
						offset := start
						for i, v := range s {
							d := p.Distance(50.857933, -0.752594, v.Latitude, v.Longitude)
							fmt.Printf("distance: %.2f pos: %.7f,%.7f speed: %.2f off: %s\n", d, v.Latitude, v.Longitude, v.Speed, offset)
							v.Offset = offset
							offset += inc
							s[i] = v
						}
					}
				}
			}
		}
	}

	// Output how many of each FourCC we have and how long they represent.
	durSec := float64(dur) / float64(time.Second)
	for k, v := range counts {
		fmt.Println(k, "=", v, float64(v)/durSec)
	}
}

func (d *Decoder) decodeTrak(rs io.ReadSeeker, stbl *mp4.StblBox, units time.Duration) ([]*Element, error) {
	chunkOffsets, err := d.chunkOffsets(stbl)
	if err != nil {
		return nil, err
	}

	// Chunks contain one or more contiguous samples.
	// Sample to time table.
	stts := stbl.Stts
	// Sample to chunk table - chunk offset
	stsc := stbl.Stsc
	// Sample sizes (framing) - size of each sample.
	stsz := stbl.Stsz

	// Entries in stsc box.
	entries := len(stsc.FirstChunk)
	lastSampleNr := stbl.Stsz.GetNrSamples() - 1

	var (
		timeIdx            int
		dec                uint64
		chunkNr            uint32
		firstSampleInChunk uint32 = 1
	)

	var data []*Element
	timeNext := stts.SampleCount[timeIdx]
	dur := stts.SampleTimeDelta[timeIdx]
	for i := 0; i < entries; i++ {
		chunkNr = stsc.FirstChunk[i]
		chunkLen := stsc.SamplesPerChunk[i]

		// Used to change group of chunks.
		var nextEntryStart uint32
		if i < entries-1 {
			nextEntryStart = stsc.FirstChunk[i+1]
		}

		for {
			nextChunkStart := firstSampleInChunk + chunkLen
			offset := chunkOffsets[chunkNr-1]
			start := dec
			var chunkSize int64
			for s, l := firstSampleInChunk, firstSampleInChunk+chunkLen; s < l; s++ {
				size := stsz.GetSampleSize(int(s))
				if s > timeNext {
					timeIdx++
					timeNext = stts.SampleCount[timeIdx]
					dur = stts.SampleTimeDelta[timeIdx]
				}
				dec += uint64(dur)
				chunkSize += int64(size)
			}

			cd, err := d.readChunk(rs, int64(offset), chunkSize, start, dec, units)
			if err != nil {
				return nil, err
			}

			data = append(data, cd...)
			if lastSampleNr < firstSampleInChunk {
				break
			}

			chunkNr++
			firstSampleInChunk = nextChunkStart
			if chunkNr == nextEntryStart {
				break
			}
		}
	}

	return data, nil
}
