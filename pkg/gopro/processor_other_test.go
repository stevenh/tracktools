//go:build !windows

package gopro

import (
	"os/exec"
	"syscall"
)

// sig sends a SIGINT to pid.
func sig(pid int) error {
	return syscall.Kill(pid, syscall.SIGINT)
}

// cmdSetup for other OSes does nothing.
func cmdSetup(_ *exec.Cmd) {
}
