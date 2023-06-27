package main

import (
	"fmt"
	"syscall"
)

func main() {
	fmtPrint("hello")
	fmt.Println("hello")
}

func fmtPrint(message string) {
	message += "\n"
	// Get the file descriptor for standard output
	fd := syscall.Stdout

	// Convert the string message to a byte slice
	bytes := []byte(message)

	// Use syscall.Write to write the byte slice to standard output
	_, err := syscall.Write(fd, bytes)
	if err != nil {
		panic(err)
	}
}
