package laptimer

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	// fixInput is a replacer which updates LapTimer's xml output so we can
	// get a match.
	fixInput = strings.NewReplacer(
		// LapTimer doesn't escape quotes in all fields e.g. vehicle->name.
		"'", "&apos;",
		// LapTimer uses windows-1252 not UTF-8.
		`<?xml version="1.0" encoding="windows-1252"?>`,
		`<?xml version="1.0" encoding="UTF-8"?>`,
	)
)

func TestDecoder(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{
			name: "single-lap",
			file: "../../testdata/LapTimer-0009-20220607-110056.hlptr",
		},
		{
			name: "multi-lap-obd",
			file: "../../testdata/LapTimer-All-20220621-092410.hlptr",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.file)
			require.NoError(t, err)
			defer f.Close()

			buf, err := ioutil.ReadAll(f)
			require.NoError(t, err)

			in := bytes.NewBuffer(buf)

			var db DB
			d := NewDecoder(in)
			err = d.Decode(&db)
			require.NoError(t, err)

			var out bytes.Buffer
			e, err := NewEncoder(&out)
			require.NoError(t, err)

			err = e.Encode(db)
			require.NoError(t, err)

			// Fix known issues in LapTimer xml before comparing.
			expected := strings.TrimSpace(fixInput.Replace(string(buf)))
			require.Equal(t, expected, out.String())
		})
	}
}
