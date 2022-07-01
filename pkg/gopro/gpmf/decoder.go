package gpmf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/edgeware/mp4ff/mp4"
	"github.com/tidwall/geodesic"
)

const (
	handlerType = "meta"
	handlerName = "GoPro MET"
)

// Decoder is a GoPro mp4 metadata decoder.
type Decoder struct {
}

// Decode decodes metadata from the mp4 stream in rs.
func (d *Decoder) Decode(rs io.ReadSeeker) error {
	line := geodesic.WGS84.DirectLine(50.857950, -0.752633, 173, 100, 0)
	fmt.Println("line:", line)
	f, err := mp4.DecodeFile(rs)
	if err != nil {
		return fmt.Errorf("decode: mp4 %w", err)
	}

	reader := NewReader()
	for _, trak := range f.Moov.Traks {
		if trak.Mdia.Hdlr.HandlerType != handlerType {
			continue
		} else if !strings.Contains(trak.Mdia.Hdlr.Name, handlerName) {
			continue
		}

		// Chunks contain one or more contiguous samples.
		// stts: sample-to-time table.
		// stsc: sample-to-chunk table - chunk offset?
		// stsz: sample sizes (framing) - size of each sample.
		stbl := trak.Mdia.Minf.Stbl
		var chunkOffsets []uint64
		switch {
		case stbl.Stco != nil:
			chunkOffsets = make([]uint64, len(stbl.Stco.ChunkOffset))
			for i := range stbl.Stco.ChunkOffset {
				chunkOffsets[i] = uint64(stbl.Stco.ChunkOffset[i])
			}
		case stbl.Co64 != nil:
			chunkOffsets = stbl.Co64.ChunkOffset
		default:
			return fmt.Errorf("decode: neither stco nor co64 available")
		}

		stsc := stbl.Stsc
		stsz := stbl.Stsz
		stts := stbl.Stts
		entries := len(stsc.FirstChunk) // Entries in stsc box.
		lastSampleNr := stbl.Stsz.GetNrSamples() - 1

		var (
			timeIdx            int
			dec                uint64
			chunkNr            uint32
			firstSampleInChunk uint32 = 1
		)

		timeUnits := time.Second / time.Duration(trak.Mdia.Mdhd.Timescale)
		timeNext := stts.SampleCount[timeIdx]
		dur := stts.SampleTimeDelta[timeIdx]
		for i := 0; i < entries; i++ {
			chunkNr = stsc.FirstChunk[i]
			chunkLen := stsc.SamplesPerChunk[i]

			var nextEntryStart uint32 // Used to change group of chunks
			if i < entries-1 {
				nextEntryStart = stsc.FirstChunk[i+1]
			}

			for {
				nextChunkStart := firstSampleInChunk + chunkLen
				offset := chunkOffsets[chunkNr-1]
				start := time.Duration(dec) * timeUnits
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

				if _, err := rs.Seek(int64(offset), io.SeekStart); err != nil {
					return fmt.Errorf("decode: seek: %w", err)
				}

				data, err := reader.Read(io.LimitReader(rs, chunkSize))
				if err != nil {
					return fmt.Errorf("decode: read: %w", err)
				}

				end := time.Duration(dec) * timeUnits
				fmt.Println("start:", start, "end:", end, "len:", len(data))

				counts := make(map[string]int)
				for _, e := range data {
					counts[e.Header.FourCC()]++
					for _, e := range e.Nested {
						counts[e.Header.FourCC()]++
						for _, e := range e.Nested {
							if e.Data == nil {
								continue
							}
							v := reflect.ValueOf(e.Data)
							t := v.Type()
							if t.Kind() == reflect.Slice {
								counts[e.Header.FourCC()] += v.Len()
							}
						}
					}
				}
				fmt.Println("counts:", counts)
				dur := float64(end-start) / float64(time.Second)
				fmt.Println("dur:", dur)
				for k, v := range counts {
					fmt.Println(k, "=", v, float64(v)/dur)
				}

				// TODO(steve): remove
				d := struct {
					Data []*Element
				}{Data: data}

				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				if err := enc.Encode(d); err != nil {
					return err
				}
				return nil
				// END REMOVE

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
	}

	return nil
}
