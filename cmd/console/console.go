package console

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

var (
	consoleMutex = sync.Mutex{}

	version   = ""
	checkouts = 0
	carts     = 0
	errors    = 0
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

	version = ver
	_, err := SetConsoleTitle(fmt.Sprintf("Scotty %v | Checkouts: %v | Carts: %v | Errors: %v", ver, checkouts, carts, errors))
	if err != nil {
		return
	}
}

func IncreaseCheckouts() {
	consoleMutex.Lock()
	defer consoleMutex.Unlock()

	checkouts++
	_, err := SetConsoleTitle(fmt.Sprintf("Scotty %v | Checkouts: %v | Carts: %v | Errors: %v", version, checkouts, carts, errors))
	if err != nil {
		return
	}
}

func IncreaseCarts() {
	consoleMutex.Lock()
	defer consoleMutex.Unlock()

	carts++
	_, err := SetConsoleTitle(fmt.Sprintf("Scotty %v | Checkouts: %v | Carts: %v | Errors: %v", version, checkouts, carts, errors))
	if err != nil {
		return
	}
}

func IncreaseErrors() {
	consoleMutex.Lock()
	defer consoleMutex.Unlock()

	errors++
	_, err := SetConsoleTitle(fmt.Sprintf("Scotty %v | Checkouts: %v | Carts: %v | Errors: %v", version, checkouts, carts, errors))
	if err != nil {
		return
	}
}

func DecreaseCarts() {
	consoleMutex.Lock()
	defer consoleMutex.Unlock()

	checkouts--
	_, err := SetConsoleTitle(fmt.Sprintf("Scotty %v | Checkouts: %v | Carts: %v | Errors: %v", version, checkouts, carts, errors))
	if err != nil {
		return
	}
}
