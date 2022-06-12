package gopro

import (
	"fmt"
	"os/exec"
	"syscall"
)

// sig sends a CTRL_BREAK_EVENT to pid.
func sig(pid int) error {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return fmt.Errorf("sig load: %w", err)
	}

	f, err := kernel32.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return fmt.Errorf("sig find proc: %w", err)
	}

	if r, _, err := f.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid)); r == 0 {
		return fmt.Errorf("sig call: %w", err)
	}

	return nil
}

// cmdSetup for windows configures a new process group so that
// sig only effects out helper process.
func cmdSetup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
