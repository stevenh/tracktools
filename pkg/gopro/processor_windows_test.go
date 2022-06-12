package gopro

import (
	"fmt"
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
