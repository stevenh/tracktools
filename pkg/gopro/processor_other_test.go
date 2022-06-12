//go:build !windows

package gopro

import (
	"syscall"
)

// sig sends a SIGINT to pid.
func sig(pid int) error {
	return syscall.Kill(pid, syscall.SIGINT)
}
