// Log the panic under windows to the log file
//
// Code from minix, via
//
// http://play.golang.org/p/kLtct7lSUg

//go:build windows

package service

import (
	"os"
	"syscall"
	"turboengine/common/log"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func setStdHandle(stdHandle int32, handle syscall.Handle) error {
	r0, _, e1 := syscall.SyscallN(procSetStdHandle.Addr(), 2, uintptr(stdHandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return e1
		}
		return syscall.EINVAL
	}
	return nil
}

// redirectStderr to the file passed in
func redirectStderr(f *os.File) {
	err := setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		log.Errorf("Failed to redirect stderr to file: %v", err)
	}
	// SetStdHandle does not affect prior references to stderr
	os.Stderr = f
}
