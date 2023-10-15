package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	winmm                   = syscall.NewLazyDLL("winmm.dll")
	capCreateCaptureWindowW = winmm.NewProc("capCreateCaptureWindowW")
	SendMessage             = winmm.NewProc("SendMessageW")
)

const (
	WM_CAP_START             = 0x0400
	WM_CAP_DRIVER_CONNECT    = WM_CAP_START + 10
	WM_CAP_DRIVER_DISCONNECT = WM_CAP_START + 11
	WM_CAP_SAVEDIB           = WM_CAP_START + 25
	WM_CAP_EDIT_COPY         = WM_CAP_START + 30
)

func main() {
	// Initialize the Windows Multimedia API
	hwnd := createCaptureWindow(0, 0, 0, 0, 320, 240)
	if hwnd == 0 {
		fmt.Println("Failed to create capture window")
		os.Exit(1)
	}
	defer sendMessage(hwnd, WM_CAP_DRIVER_DISCONNECT, 0, 0)

	// Connect to the default camera
	if !sendMessage(hwnd, WM_CAP_DRIVER_CONNECT, 0, 0) {
		fmt.Println("Failed to connect to the camera")
		os.Exit(1)
	}

	// Start capturing
	if !sendMessage(hwnd, WM_CAP_SAVEDIB, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("output.bmp")))) {
		fmt.Println("Failed to capture video")
		os.Exit(1)
	}

	fmt.Println("Press Enter to stop capturing...")
	fmt.Scanln()

	// Disconnect from the camera
	sendMessage(hwnd, WM_CAP_DRIVER_DISCONNECT, 0, 0)
}

func createCaptureWindow(x, y, width, height, parentHwnd int, id int) uintptr {
	title := "CaptureWindow"
	ret, _, _ := capCreateCaptureWindowW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(0x40000000|0x10000000), // WS_VISIBLE|WS_CHILD
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(parentHwnd),
		uintptr(id),
	)

	return ret
}

func sendMessage(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) bool {
	ret, _, err := syscall.Syscall6(
		SendMessage.Addr(),
		4,
		hwnd,
		uintptr(msg),
		wParam,
		lParam,
		0,
		0,
	)
	if err != 0 {
		fmt.Printf("SendMessage failed with error code: %d\n", err)
		return false
	}
	return ret != 0
}
