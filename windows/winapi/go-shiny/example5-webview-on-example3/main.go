package main

import (
	"fmt"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"image"
	"log"
	"syscall"
	"unsafe"
)

const dllPath = "WebView2Loader.dll"
const functionName = "CreateWebView2Environment"

type (
	IWebView2Environment    uintptr
	IWebView2CreateWebView2 uintptr
	IWebView2WebView        uintptr
)

func main() {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  "Hello World",
			Width:  1000,
			Height: 600,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		var t screen.Texture

		buffer, err := s.NewBuffer(image.Point{X: 1000, Y: 600})
		if err != nil {
			log.Fatal(err)
		}
		defer buffer.Release()
		size2 := buffer.Bounds().Size()
		fmt.Println(size2)

		// Load the DLL
		dllHandle, err := syscall.LoadLibrary(dllPath)
		if err != nil {
			log.Fatal(err)
		}
		defer syscall.FreeLibrary(dllHandle)

		// Get the function address from the DLL
		functionAddress, err := syscall.GetProcAddress(dllHandle, functionName)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("functionAddress: %v \n", functionAddress)

		// Convert the function address to a callable function pointer
		createEnv := func() IWebView2Environment { return IWebView2Environment(functionAddress) }
		fmt.Printf("createEnv: %v \n", createEnv)

		// Create the WebView2 environment
		env := createEnv()
		fmt.Printf("env: %v \n", env)

		// Create the WebView2 control
		createWebViewFn := func(env IWebView2Environment, a, b uintptr) IWebView2WebView {
			return IWebView2WebView(unsafe.Pointer(uintptr(env) + a + b))
		}

		fmt.Printf("createWebViewFn: %v \n", createWebViewFn)
		webview := createWebViewFn(env, uintptr(0), uintptr(0))

		fmt.Printf("webview: %v \n", webview)

		// Set the WebView2 window handle to the Shiny window's handle
		setWindowHandle := func(webview IWebView2WebView, handle uintptr) { /* implementation */ }
		setWindowHandle(webview, uintptr(40))

		fmt.Printf("setWindowHandle: %v \n", setWindowHandle)

		//setWindowHandle(webview, uintptr(w.DriverWindow().PlatformWindow()))

		// Navigate to a URL
		navigate := func(webview IWebView2WebView, url uintptr) { /* implementation */ }
		//navigate(webview, uintptr(24))

		navigate(webview, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("https://google.dk"))))

		for {
			e := w.NextEvent()

			format := "got %#v\n"
			if _, ok := e.(fmt.Stringer); ok {
				format = "got %v\n"
			}
			fmt.Printf(format, e)

			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}

			case paint.Event:
				if t != nil {
					t.Release()
				}

				fmt.Println("here is: --> ", w, " <---- ")

				w.Publish()

			case size.Event:
				w.Send(paint.Event{})
			}
		}
	})
}
