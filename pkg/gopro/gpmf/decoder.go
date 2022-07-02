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

func NewDecoder() *Decoder {
	return &Decoder{
		reader: NewReader(),
	}
}

// Decode decodes metadata from the mp4 stream in rs.
func (d *Decoder) Decode(rs io.ReadSeeker) error {
	f, err := mp4.DecodeFile(rs)
	if err != nil {
		return fmt.Errorf("decode: mp4 %w", err)
	}

	for i, trak := range f.Moov.Traks {
		if trak.Mdia.Hdlr.HandlerType != handlerType {
			continue
		} else if !strings.Contains(trak.Mdia.Hdlr.Name, handlerName) {
			continue
		}

		stbl := trak.Mdia.Minf.Stbl

		timeUnits := time.Second / time.Duration(trak.Mdia.Mdhd.Timescale)
		if err := d.decodeTrak(rs, stbl, timeUnits); err != nil {
			return fmt.Errorf("decode: trak %d: %w", i, err)
		}
	}

	return nil
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

func (d *Decoder) readChunk(rs io.ReadSeeker, offset, size int64, startDec, endDec uint64, timeUnits time.Duration) error {
	if _, err := rs.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("seek: %w", err)
	}

	data, err := d.reader.Read(io.LimitReader(rs, size))
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	// TODO(steve): remove
	d.dumpStats(data, startDec, endDec, timeUnits)
	//return dump(data)
	return nil
}

// TODO(steve): remove
func (d *Decoder) dumpStats(data []*Element, startDec, endDec uint64, timeUnits time.Duration) {
	start := time.Duration(startDec) * timeUnits
	end := time.Duration(endDec) * timeUnits
	dur := float64(end-start) / float64(time.Second)
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
						for _, v := range s {
							d := p.Distance(50.857933, -0.752594, v.Latitude, v.Longitude)
							fmt.Printf("distance: %.2f pos: %.7f,%.7f\n", d, v.Latitude, v.Longitude)
						}
					}
				}
			}
		}
	}

	// Output how many of each FourCC we have and how long they represent.
	for k, v := range counts {
		fmt.Println(k, "=", v, float64(v)/dur)
	}
}

func (d *Decoder) decodeTrak(rs io.ReadSeeker, stbl *mp4.StblBox, timeUnits time.Duration) error {
	chunkOffsets, err := d.chunkOffsets(stbl)
	if err != nil {
		return err
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

			if err := d.readChunk(rs, int64(offset), chunkSize, start, dec, timeUnits); err != nil {
				return err
			}

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

	return nil
}
