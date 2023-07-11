package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procSetCursorPos     = user32.NewProc("SetCursorPos")
	procGetCursorPos     = user32.NewProc("GetCursorPos")
	procMouseEvent       = user32.NewProc("mouse_event")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

type POINT struct {
	X, Y int32
}

func main() {
	// Set a delay to give you time to position the mouse
	//time.Sleep(3 * time.Second)

	// Get the current mouse coordinates
	var pt POINT

	//go func() {
	//	for {
	//		//fmt.Printf("Current mouse position: X=%d, Y=%d\n", pt.X, pt.Y)
	//		time.Sleep(10 * time.Millisecond)
	//	}
	//}()

	for {
		GetCursorPos(&pt)
		time.Sleep(10 * time.Millisecond)

		if (pt.X >= 100 && pt.X <= 200) && (pt.Y >= 100 && pt.Y <= 200) {
			fmt.Println("inside rectangle!")
		} else {
			fmt.Printf("Current mouse position: X=%d, Y=%d\n", pt.X, pt.Y)
		}
	}

	// Move the mouse to a new position
	//newX, newY := pt.X+100, pt.Y+100
	//SetCursorPos(newX, newY)

	// Get the updated mouse coordinates
	//GetCursorPos(&pt)
	//fmt.Printf("New mouse position: X=%d, Y=%d\n", pt.X, pt.Y)
}

func SetCursorPos(x, y int32) {
	procSetCursorPos.Call(uintptr(x), uintptr(y))
}

func GetCursorPos(pt *POINT) {
	procGetCursorPos.Call(uintptr(unsafe.Pointer(pt)))
}

func mouseEvent(dwFlags, dx, dy, dwData, dwExtraInfo uintptr) {
	procMouseEvent.Call(dwFlags, dx, dy, dwData, dwExtraInfo)
}

func GetSystemMetrics(nIndex int32) int32 {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(nIndex))
	return int32(ret)
}
