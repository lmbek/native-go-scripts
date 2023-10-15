package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	WM_CAP_START              = 0x400
	WM_CAP_DRIVER_CONNECT     = WM_CAP_START + 10
	WM_CAP_DRIVER_DISCONNECT  = WM_CAP_START + 11
	WM_CAP_SINGLE_FRAME_OPEN  = WM_CAP_START + 70
	WM_CAP_SINGLE_FRAME_CLOSE = WM_CAP_START + 71
	WM_CAP_SINGLE_FRAME       = WM_CAP_START + 72
	WM_CAP_FILE_SAVEDIB       = WM_CAP_START + 25
	WS_VISIBLE                = 0x10000000
	WS_CHILD                  = 0x40000000
)

var (
	avicap32 = syscall.NewLazyDLL("avicap32.dll")

	// AVICap32 API functions
	capCreateCaptureWindowW = avicap32.NewProc("capCreateCaptureWindowW")
	SendMessage             = avicap32.NewProc("SendMessageW")
	DestroyWindow           = avicap32.NewProc("DestroyWindow")
	capFileSaveDIB          = avicap32.NewProc("capFileSaveDIB")
)

func main() {
	// Initialize the capture window
	captureWindow := initCaptureWindow()

	fmt.Println(captureWindow)
	if captureWindow == 0 {
		fmt.Println("Failed to initialize capture window")
		return
	}
	defer DestroyWindow.Call(captureWindow)

	// Start capturing
	_, _, _ = SendMessage.Call(captureWindow, WM_CAP_DRIVER_CONNECT, 0, 0)

	// Capture a single frame
	_, _, _ = SendMessage.Call(captureWindow, WM_CAP_SINGLE_FRAME_OPEN, 0, 0)
	_, _, _ = SendMessage.Call(captureWindow, WM_CAP_SINGLE_FRAME_CLOSE, 0, 0)

	// Save the captured image
	imageFile := "captured_image.bmp"
	_, _, _ = SendMessage.Call(captureWindow, WM_CAP_FILE_SAVEDIB, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(imageFile))))

	// Stop capturing
	_, _, _ = SendMessage.Call(captureWindow, WM_CAP_DRIVER_DISCONNECT, 0, 0)

	fmt.Println("Image saved as", imageFile)
}

func initCaptureWindow() uintptr {
	// Create a capture window
	title := "Live! Cam Sync 1080p"
	titlePtr, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		fmt.Println("Error converting title to UTF16:", err)
		return 0
	}

	captureWindow, _, _ := capCreateCaptureWindowW.Call(
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(WS_VISIBLE|WS_CHILD),
		0, 0, 320, 240,
		0, 0,
		0, 0,
	)

	return captureWindow
}
