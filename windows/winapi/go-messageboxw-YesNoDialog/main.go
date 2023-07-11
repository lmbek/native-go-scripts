package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func main() {
	CustomYesNoDialog()
}

func CustomYesNoDialog() {
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
	lpCaption, _ := syscall.UTF16PtrFromString("What do you want?") // LPCWSTR
	lpText, _ := syscall.UTF16PtrFromString("Press Yes or No")      // LPCWSTR

	buttonPressed, _, _ := syscall.SyscallN(procMessageBox.Addr(),
		0,
		uintptr(unsafe.Pointer(lpText)),
		uintptr(unsafe.Pointer(lpCaption)),
		MB_YESNO|MB_ICONQUESTION, // Let the window TOPMOST.
	)

	const yes = 6
	const no = 7

	if buttonPressed == yes {
		fmt.Println("You pressed yes")
	} else if buttonPressed == no {
		fmt.Println("You pressed no")
		os.Exit(0)
	}

}
