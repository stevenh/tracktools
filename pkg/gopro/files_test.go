package gopro

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileSlice(t *testing.T) {
	tests := []struct {
		name     string
		data     FileSlice
		expected FileSlice
		err      error
	}{
		{
			name: "valid",
			data: FileSlice{
				File{Chapter: "02"},
				File{Chapter: "01"},
			},
			expected: FileSlice{
				File{Chapter: "01"},
				File{Chapter: "02"},
			},
		},
		{
			name: "unexpected-first-chapter",
			data: FileSlice{
				File{Chapter: "02"},
			},
			expected: FileSlice{
				File{Chapter: "02"},
			},
			err: chapterError{chapter: "02"},
		},
		{
			name: "unexpected-chapter",
			data: FileSlice{
				File{Chapter: "03"},
				File{Chapter: "01"},
			},
			expected: FileSlice{
				File{Chapter: "01"},
				File{Chapter: "03"},
			},
			err: chapterError{chapter: "03", idx: 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sort.Sort(tc.data)
			require.Equal(t, tc.expected, tc.data)

			err := tc.data.Validate()
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
