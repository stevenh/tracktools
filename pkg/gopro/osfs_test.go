package gopro

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_ReadDir(t *testing.T) {
	tf, err := os.CreateTemp(t.TempDir(), "os-test")
	require.NoError(t, err)

	name := tf.Name()
	require.NoError(t, tf.Close())

	dirs, err := os.ReadDir(name)
	require.ErrorIs(t, err, syscall.ENOTDIR)
	require.Empty(t, dirs)
}

func TestOSFS(t *testing.T) {
	o := osFS{}
	tf, err := o.CreateTemp("", "gopro-test")
	require.NoError(t, err)
	require.NoError(t, tf.Close())

	name := tf.Name()

	t.Cleanup(func() { os.Remove(name) }) //nolint: errcheck

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
	require.Empty(t, dirs)

	err = o.Remove(name)
	require.NoError(t, err)
}
