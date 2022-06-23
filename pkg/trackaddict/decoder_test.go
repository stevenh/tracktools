package trackaddict

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	f, err := os.Open("../../testdata/Log-20220531-085930 Goodwood Motorcircuit - 2.57.527.csv")
	require.NoError(t, err)
	defer f.Close() // nolint: errcheck

	d, err := NewDecoder(f)
	require.NoError(t, err)

	_, err = d.Decode()
	require.NoError(t, err)
}
