package gpmf

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Eyevinn/mp4ff/mp4"
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
			// Not our handler type.
			continue
		} else if !strings.Contains(trak.Mdia.Hdlr.Name, handlerName) {
			// Not our handler name.
			continue
		}

		units := time.Second / time.Duration(trak.Mdia.Mdhd.Timescale)
		data, err := d.decodeTrak(rs, trak.Mdia.Minf.Stbl, units)
		if err != nil {
			return nil, fmt.Errorf("decode: trak %d: %w", i, err)
		}

		return data, nil
	}

	return nil, fmt.Errorf("decode: no metadata for %q found", handlerName)
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

// readChunk reads a single chunk from rs and amends it with offset information.
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

	if err := Walk(data, newOffsetWalker(startDec, endDec, units).walk); err != nil {
		return nil, fmt.Errorf("offsets: %w", err)
	}

	return data, nil
}

// decodeTrak decodes all chunks from single tracks data as detailed in stbl from rs.
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
	entries := len(stsc.Entries)
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
	for i, entry := range stsc.Entries {
		chunkNr = entry.FirstChunk
		chunkLen := entry.SamplesPerChunk

		// Used to change group of chunks.
		var nextEntryStart uint32
		if i < entries-1 {
			nextEntryStart = stsc.Entries[i+1].FirstChunk
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

			cd, err := d.readChunk(rs, int64(offset), chunkSize, start, dec, units) //nolint: gosec
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
