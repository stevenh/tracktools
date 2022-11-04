package gpmf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	tests := []string{
		"hero5",
		"fusion",
		"hero6",
		"hero6+ble",
		"hero7",
		"hero8",
	}

	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			f, err := os.Open(filepath.Join("../../../test", tc+".mp4"))
			require.NoError(t, err)
			defer f.Close() //nolint: errcheck

			dec := &Decoder{}
			_, err = dec.Decode(f)
			require.NoError(t, err)

			// TODO(steve): validate data.
		})
	}
}
