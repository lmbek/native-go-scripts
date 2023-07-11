package main

import (
	"syscall"
	"unsafe"
)

func main() {
	CustomWarning()
}

func CustomWarning() {
	// https://stackoverflow.com/questions/46705163/how-to-alert-in-go-to-show-messagebox
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw
	var user32DLL = syscall.NewLazyDLL("user32.dll")
	var procMessageBox = user32DLL.NewProc("MessageBoxW") // Return value: Type int
	const (
		MB_OK           = 0x00000000
		MB_OKCANCEL     = 0x00000001
		MB_YESNO        = 0x00000004
		MB_SYSTEMMODAL  = 0x00001000
		MB_ICONQUESTION = 0x00000020
		MB_ICONWARNING  = 0x00000030
		MB_ICONASTERISK = 0x00000040
	)

	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw#return-value
	lpCaption, _ := syscall.UTF16PtrFromString("This is a custom warning")                                              // LPCWSTR
	lpText, _ := syscall.UTF16PtrFromString("As you are aware\n\nThis is a native solution\n\nThank you for trying it") // LPCWSTR

	_, _, _ = syscall.SyscallN(procMessageBox.Addr(),
		0,
		uintptr(unsafe.Pointer(lpText)),
		uintptr(unsafe.Pointer(lpCaption)),
		MB_OK|MB_ICONWARNING, // Let the window TOPMOST.
	)
}
