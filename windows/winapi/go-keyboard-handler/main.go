package main

import (
	"fmt"
	"strings"
	"syscall"
)

var (
	user32DLL           = syscall.NewLazyDLL("user32.dll")
	keybdEvent          = user32DLL.NewProc("keybd_event")
	getForegroundWindow = user32DLL.NewProc("GetForegroundWindow")
	setForegroundWindow = user32DLL.NewProc("SetForegroundWindow")

	kernel32DLL              = syscall.NewLazyDLL("kernel32.dll")
	getCurrentThreadId       = kernel32DLL.NewProc("GetCurrentThreadId")
	getWindowThreadProcessId = user32DLL.NewProc("GetWindowThreadProcessId")
	attachThreadInput        = user32DLL.NewProc("AttachThreadInput")
)

const (
	keyEventKeyDown = 0
	keyEventKeyUp   = 2
)

func sendKey(keyCode int, keyEvent int) {
	// Get the foreground window handle
	hWnd, _, _ := getForegroundWindow.Call()

	// Get the thread identifier of the foreground window
	foregroundThreadId, _, _ := getWindowThreadProcessId.Call(hWnd, 0)

	// Get the current thread identifier
	currentThreadId, _, _ := getCurrentThreadId.Call()

	// Attach the input processing mechanism of the foreground window's thread to the current thread
	_, _, _ = attachThreadInput.Call(currentThreadId, foregroundThreadId, 1)

	// Simulate the key event
	_, _, _ = keybdEvent.Call(uintptr(keyCode), 0, uintptr(keyEvent), 0)

	// Detach the input processing mechanism
	_, _, _ = attachThreadInput.Call(currentThreadId, foregroundThreadId, 0)
}

func sendKeySequence(keys ...int) {
	for _, key := range keys {
		sendKey(key, keyEventKeyDown)
		sendKey(key, keyEventKeyUp)
	}
}

func typeParagraph(paragraph string) {
	// Convert the paragraph to uppercase for simplicity
	paragraph = strings.ToUpper(paragraph)

	//this is a sample paragraph that can do ask i ask i like and cake is very nice thanks
	//this is a sample paragraph

	//Iterate over each character in the paragraph and simulate typing
	for _, ch := range paragraph {
		if ch == ' ' {
			// Simulate pressing the space key
			sendKey(0x20, keyEventKeyDown)
			sendKey(0x20, keyEventKeyUp)
		} else {
			// Get the virtual key code for the character
			keyCode := int(ch)

			// Simulate pressing and releasing the key
			sendKey(keyCode, keyEventKeyDown)
			sendKey(keyCode, keyEventKeyUp)
		}
	}
}

/*
	0x41, // A
	0x42, // B
	0x43, // C
	0x44, // D
	0x45, // E
	0x46, // F
	0x47, // G
	0x48, // H
	0x49, // I
	0x4A, // J
	0x4B, // K
	0x4C, // L
	0x4D, // M
	0x4E, // N
	0x4F, // O
	0x50, // P
	0x51, // Q
	0x52, // R
	0x53, // S
	0x54, // T
	0x55, // U
	0x56, // V
	0x57, // W
	0x58, // X
	0x59, // Y
	0x5A, // Z
	0x30, // 0
	0x31, // 1
	0x32, // 2
	0x33, // 3
	0x34, // 4
	0x35, // 5
	0x36, // 6
	0x37, // 7
	0x38, // 8
	0x39, // 9
*/

// PLEASE NOTE: This script is not perfect, it works with certain characters, but use with caution.
// You can make a map of accepted keys and ensure no wrong output like a delete button is pressed.

func main() {
	// Delay for a few seconds to allow time to focus on the desired window
	// You can remove or adjust this delay as per your needs
	//time.Sleep(1 * time.Second)

	// Example on when used:
	// 1995 is the year i was born in
	//

	// space
	sendKey(0x20, keyEventKeyDown)
	sendKey(0x20, keyEventKeyUp)

	// 9
	sendKey(0x31, keyEventKeyDown)
	sendKey(0x31, keyEventKeyUp)

	// 9
	sendKey(0x39, keyEventKeyDown)
	sendKey(0x39, keyEventKeyUp)

	// 9
	sendKey(0x39, keyEventKeyDown)
	sendKey(0x39, keyEventKeyUp)

	// 5
	sendKey(0x35, keyEventKeyDown)
	sendKey(0x35, keyEventKeyUp)

	// space
	sendKey(0x20, keyEventKeyDown)
	sendKey(0x20, keyEventKeyUp)

	// Specify the paragraph to type
	paragraph := "is the year i was born in"

	// Type the paragraph
	typeParagraph(paragraph)

	fmt.Println("Typed the paragraph.")
}
