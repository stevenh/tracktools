package gopro

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigERror(t *testing.T) {
	e := configError("test")
	require.Equal(t, `missing config "test"`, e.Error())
}

func TestMatchError(t *testing.T) {
	e := matchError(1)
	require.Equal(t, "unexpected match len: 1", e.Error())
}

func TestChapterError(t *testing.T) {
	e := chapterError{
		chapter: "test",
		idx:     1,
	}
	require.Equal(t, `unexpected chapter "test" at index 1`, e.Error())
}
