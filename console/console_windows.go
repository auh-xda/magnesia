//go:build windows
// +build windows

package console

import (
	"os"
	"syscall"
	"unsafe"
)

func EnableVirtualTerminal() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")

	handle := syscall.Handle(os.Stdout.Fd())
	var mode uint32
	getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))

	// ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	setConsoleMode.Call(uintptr(handle), uintptr(mode|ENABLE_VIRTUAL_TERMINAL_PROCESSING))
}
