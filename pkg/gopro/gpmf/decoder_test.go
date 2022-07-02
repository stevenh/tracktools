package gpmf

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	tests := []struct {
		name   string
		reader func(*testing.T) io.ReadSeeker
	}{
		/*
			{
				name: "basic-nested",
				reader: func(t *testing.T) io.ReadSeeker {
					t.Helper()
					klv := []byte{
						0x44, 0x45, 0x56, 0x43, 0x00, 0x04, 0x00, 0x07,
						0x44, 0x56, 0x49, 0x44, 0x4c, 0x04, 0x00, 0x01,
						0x00, 0x00, 0x10, 0x01, 0x44, 0x56, 0x4e, 0x4d,
						0x63, 0x01, 0x00, 0x06, 0x43, 0x61, 0x6d, 0x65,
						0x72, 0x61, 0x00, 0x00,
					}
					return bytes.NewBuffer(klv)
				},
			},
			{
				name: "hero5-mp4",
				reader: func(t *testing.T) io.ReadSeeker {
					t.Helper()
					f, err := os.Open("../../../test/hero5.mp4")
					require.NoError(t, err)
					return f
				},
			},
			{
				name: "fusion-mp4",
				reader: func(t *testing.T) io.ReadSeeker {
					t.Helper()
					f, err := os.Open("../../../test/fusion.mp4")
					require.NoError(t, err)
					return f
				},
			},
			{
				name: "hero6-mp4",
				reader: func(t *testing.T) io.ReadSeeker {
					t.Helper()
					f, err := os.Open("../../../test/hero6.mp4")
					require.NoError(t, err)
					return f
				},
			},
		*/
		{
			name: "goodwood-mp4",
			reader: func(t *testing.T) io.ReadSeeker {
				t.Helper()
				f, err := os.Open("../../../test/GOPR1100-JOINED.mp4")
				//f, err := os.Open("../../../test/GOPR1109.mp4")
				require.NoError(t, err)
				return f
			},
		},
	}

	reader := NewReader()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.reader(t)
			if rc, ok := r.(io.Closer); ok {
				defer rc.Close() // nolint: errcheck
			}

			dec := &Decoder{}
			err := dec.Decode(r)
			require.NoError(t, err)

			return

			_, err = r.Seek(int64(381620), io.SeekStart)
			require.NoError(t, err)
			//4568 dur: 1001 offset: 381620

			data, err := reader.Read(io.LimitReader(r, 4568))
			//data, err := reader.Read(r)
			require.NoError(t, err)

			d := struct {
				Data []*Element
			}{Data: data}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			err = enc.Encode(d)
			require.NoError(t, err)
		})
	}
}
