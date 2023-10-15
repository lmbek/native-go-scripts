// TODO: Rework it into Go embedding
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	webviewWindow "webview_go"
)

var count uint = 0

type Count struct {
	Value uint `json:"count"`
}

func main() {
	if err := runSingleInstance(); err != nil {
		fmt.Println(err)
		return
	}
}

func createLockFile(lockfile string) {
	file, err := os.Create(lockfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func removeLockFile(lockfile string) {
	if err := os.Remove(lockfile); err != nil {
		log.Println("Failed to remove the lock file:", err)
	}
}

func runSingleInstance() error {
	lockFile := "myapp.lock"

	_, err := os.Stat(lockFile)
	if err == nil {
		return fmt.Errorf("Another instance is already running.")
	}

	createLockFile(lockFile)

	// Handle application shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		// This code will be executed when the application is interrupted (e.g., Ctrl+C).
		removeLockFile(lockFile)
		os.Exit(1)
	}()
	defer removeLockFile(lockFile)

	// run webview
	err = startWebview("index.html", "My Webview", 1280, 720)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func startWebview(htmlFile, title string, width, height int) error {
	bytes, err := os.ReadFile(htmlFile)
	if err != nil {
		return fmt.Errorf("Error reading HTML file: %v", err)
	}

	html := string(bytes)

	webview := webviewWindow.NewWindow(true, nil)
	defer webview.Destroy()
	webview.SetTitle("My Webview")
	//w.SetSize(480, 320, webview.HintNone)
	webview.SetSize(1280, 720, webviewWindow.HintMax)

	// A binding that increments a value and immediately returns the new value.
	err = webview.Bind("increment", func() Count { return increment() })
	if err != nil {
		fmt.Println(err)
	}

	webview.SetHtml(html)
	//webview.Navigate("https://google.dk")

	webview.Run()
	return nil
}

func increment() Count {
	count++
	return Count{Value: count}
}
