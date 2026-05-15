package title

import (
	"syscall"
	"unsafe"
)

func SetTitle(title string) error {
	h, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return err
	}
	defer syscall.FreeLibrary(h)
	proc, err := syscall.GetProcAddress(h, "SetConsoleTitleW")
	if err != nil {
		return err
	}
	_, _, err = syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return err
}
