package gpmf

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	tests := []struct {
		name   string
		reader func(*testing.T) io.Reader
	}{
		/*
							{
								name: "basic-nested",
								reader: func(t *testing.T) io.Reader {
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
							name: "hero5-raw",
							reader: func(t *testing.T) io.Reader {
			t.Helper()
								f, err := os.Open("../../../test/hero5.raw")
								require.NoError(t, err)
								return f
							},
						},
						{
							name: "fusion-raw",
							reader: func(t *testing.T) io.Reader {
			t.Helper()
								f, err := os.Open("../../../test/fusion.raw")
								require.NoError(t, err)
								return f
							},
						},
		*/
		{
			name: "hero6-raw",
			reader: func(t *testing.T) io.Reader {
				t.Helper()
				f, err := os.Open("../../../test/hero6.raw")
				require.NoError(t, err)
				return f
			},
		},
		/*
			{
				name: "goodwood-raw",
				reader: func(t *testing.T) io.Reader {
					f, err := os.Open("../../../test/GOPR1109.raw")
					require.NoError(t, err)
					return f
				},
			},
		*/
	}

	reader := NewReader()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.reader(t)
			if rc, ok := r.(io.Closer); ok {
				defer rc.Close() // nolint: errcheck
			}

			data, err := reader.Read(r)
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
