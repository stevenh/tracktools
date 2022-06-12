package gopro

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMatcher(t *testing.T) {
	tests := []struct {
		name  string
		first string
		other []string
		err   bool
	}{
		{
			name:  "single",
			first: "^test$",
		},
		{
			name:  "double",
			first: "^test$",
			other: []string{"^test2$"},
		},
		{
			name:  "bad-regexp",
			first: "^test[0-9$",
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewMatcher(tc.name, tc.first, tc.other...)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, m)
		})
	}
}

func TestNewMatcherMust(t *testing.T) {
	require.Panics(t, func() {
		NewMatcherMust("test", "[0-9")
	})
}

func TestMatch(t *testing.T) {
	m, err := NewMatcher("test", "^([0-9])([0-9])([0-9])$")
	require.NoError(t, err)
	f, err := m.Match("012")
	require.Error(t, err)
	require.Nil(t, f)
}

func TestHero5(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected *File
		err      error
	}{
		{
			name: "first",
			file: "test/GOPR1234.mp4",
			expected: &File{
				Name:    "GOPR1234.mp4",
				Index:   "1234",
				Chapter: "00",
			},
		},
		{
			name: "chapter01",
			file: "test/GP011234.mp4",
			expected: &File{
				Name:    "GP011234.mp4",
				Index:   "1234",
				Chapter: "01",
			},
		},
		{
			name: "chapter02",
			file: "test/GP021234.mp4",
			expected: &File{
				Name:    "GP021234.mp4",
				Index:   "1234",
				Chapter: "02",
			},
		},
		{
			name: "no-match",
			file: "test/GH021234.mp4",
			err:  ErrNoMatch,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := Hero5.Match(tc.file)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, f)
		})
	}
}

func TestHero10(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected *File
		err      error
	}{
		{
			name: "first-hevc",
			file: "test/GX011234.mp4",
			expected: &File{
				Name:    "GX011234.mp4",
				Index:   "1234",
				Chapter: "01",
			},
		},
		{
			name: "second-hevc",
			file: "test/GX021234.mp4",
			expected: &File{
				Name:    "GX021234.mp4",
				Index:   "1234",
				Chapter: "02",
			},
		},
		{
			name: "first-avc",
			file: "test/GH011234.mp4",
			expected: &File{
				Name:    "GH011234.mp4",
				Index:   "1234",
				Chapter: "01",
			},
		},
		{
			name: "second-avc",
			file: "test/GH021234.mp4",
			expected: &File{
				Name:    "GH021234.mp4",
				Index:   "1234",
				Chapter: "02",
			},
		},
		{
			name: "no-match",
			err:  ErrNoMatch,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := Hero10.Match(tc.file)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, f)
		})
	}
}
