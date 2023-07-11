package main

/*
#cgo LDFLAGS: -lX11 -lXtst
#include <X11/Xlib.h>
#include <X11/extensions/XTest.h>
*/
import "C"
import "fmt"

func main() {
	MouseControl()
}

func MouseControl() {
	display := openDisplay()
	defer closeDisplay(display)

	rootWindow := getDefaultRootWindow(display)
	rootX, rootY := getMouseCoordinates(display, rootWindow)

	fmt.Printf("Mouse coordinates: X=%d, Y=%d\n", rootX, rootY)

	moveX := rootX + 250
	moveY := rootY + 250

	moveMouse(display, moveX, moveY)

	fmt.Println("Mouse moved successfully")
}

// Open the X display
func openDisplay() *C.Display {
	display := C.XOpenDisplay(nil)
	if display == nil {
		fmt.Println("Failed to open display")
		return nil
	}
	return display
}

// Close the X display
func closeDisplay(display *C.Display) {
	C.XCloseDisplay(display)
}

// Get the default root window of the X display
func getDefaultRootWindow(display *C.Display) C.Window {
	screen := C.XDefaultScreen(display)
	return C.XRootWindow(display, C.int(screen))
}

// Get the current mouse coordinates relative to the root window
func getMouseCoordinates(display *C.Display, rootWindow C.Window) (int, int) {
	var rootX, rootY, winX, winY C.int
	var mask C.uint

	result := C.XQueryPointer(
		display,
		rootWindow,
		&rootWindow,
		&rootWindow,
		&rootX,
		&rootY,
		&winX,
		&winY,
		&mask,
	)

	if result == 0 {
		fmt.Println("Failed to query pointer")
		return 0, 0
	}

	return int(rootX), int(rootY)
}

// Move the mouse cursor to the specified coordinates
func moveMouse(display *C.Display, x, y int) {
	C.XTestFakeMotionEvent(display, -1, C.int(x), C.int(y), 0)
	C.XFlush(display)
}
