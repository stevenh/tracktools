package gpmf

import (
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.reader(t)
			if rc, ok := r.(io.Closer); ok {
				defer rc.Close() // nolint: errcheck
			}

			dec := &Decoder{}
			_, err := dec.Decode(r)
			require.NoError(t, err)

			// TODO(steve): validate data.
		})
	}
}
