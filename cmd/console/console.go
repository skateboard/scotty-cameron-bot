package console

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

var (
	consoleMutex = sync.Mutex{}
)

func SetConsoleTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}
	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err
}

func SetBase(ver string) {
	consoleMutex.Lock()
	defer consoleMutex.Unlock()
	_, err := SetConsoleTitle(fmt.Sprintf("Scotty %v", ver))
	if err != nil {
		return
	}
}
