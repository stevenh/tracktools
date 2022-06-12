package gopro

import (
	"io/fs"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOSFS(t *testing.T) {
	o := osFS{}
	tf, err := o.CreateTemp("", "gopro-test")
	require.NoError(t, err)

	name := tf.Name()

	defer os.Remove(name) // nolint: errcheck

	now := time.Now().Round(0)
	err = o.Chtimes(name, now, now)
	require.NoError(t, err)

	fi, err := o.Stat(name)
	require.NoError(t, err)
	require.Equal(t, now, fi.ModTime())

	f, err := o.Open(name)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	dirs, err := o.ReadDir(name)
	require.ErrorIs(t, err, syscall.ENOTDIR)
	require.Equal(t, []fs.DirEntry{}, dirs)

	err = o.Remove(name)
	require.NoError(t, err)
}
