package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"syscall"
	"unsafe"
)

var width = 1920
var height = 1080

func removeFalseFlags(err error) error {
	if err.Error() == "The operation completed successfully." {
		return nil
	}
	return err
}

func main() {
	run()

	// Example usage of getColorAtPixel function
	screenshot, _ := loadImage("screenshot.png")
	r, g, b, a := getColorAtPixel(screenshot, 100, 100)
	fmt.Printf("Color at pixel (100, 100): R=%d, G=%d, B=%d, A=%d\n", r, g, b, a)
}

func run() {
	hdc, err := getDC()
	err = removeFalseFlags(err)

	if err != nil {
		fmt.Println("Failed to get device context:", err)
		return
	}
	defer releaseDC(0, hdc)

	hbitmap, _, err := createCompatibleBitmap(hdc, width, height)

	if err.Error() == "The operation completed successfully." {
		err = nil
	}

	if err != nil {
		fmt.Println("Failed to create compatible bitmap:", err)
		return
	}
	defer deleteObject(hbitmap)

	memDC, _, err := createCompatibleDC(hdc)
	err = removeFalseFlags(err)
	if err != nil {
		fmt.Println("Failed to create compatible device context:", err)
		return
	}
	defer deleteDC(memDC)

	oldObj, _, err := selectObject(memDC, hbitmap)
	err = removeFalseFlags(err)
	if err != nil {
		fmt.Println("Failed to select bitmap into device context:", err)
		return
	}
	defer selectObject(memDC, oldObj)

	success, _, err := bitblt(memDC, 0, 0, width, height, hdc, 0, 0, 0x00CC0020)
	err = removeFalseFlags(err)
	if err != nil || !success {
		fmt.Println("Failed to copy screen contents into bitmap:", err)
		return
	}

	image, err := createImageFromBitmap(hbitmap, width, height)
	if err != nil {
		fmt.Println("Failed to create image from bitmap:", err)
		return
	}

	err = saveImageAsPNG(image, "screenshot.png")
	if err != nil {
		fmt.Println("Failed to save image as PNG:", err)
		return
	}

	err = saveImageAsJPEG(image, "screenshot.jpg")
	if err != nil {
		fmt.Println("Failed to save image as JPEG:", err)
		return
	}

	fmt.Println("Screenshot saved successfully.")
}

func getDC() (syscall.Handle, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("GetDC")
	hdc, _, err := proc.Call(0)
	if hdc == 0 {
		return 0, fmt.Errorf("failed to get device context")
	}
	return syscall.Handle(hdc), err
}

func releaseDC(hwnd syscall.Handle, hdc syscall.Handle) (uintptr, uintptr, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("ReleaseDC")
	result, _, err := proc.Call(uintptr(hwnd), uintptr(hdc))
	if result == 0 {
		return 0, 0, fmt.Errorf("failed to release device context")
	}
	return result, 0, err
}

func createCompatibleBitmap(hdc syscall.Handle, width, height int) (syscall.Handle, uintptr, error) {
	gdi32 := syscall.NewLazyDLL("gdi32.dll")
	proc := gdi32.NewProc("CreateCompatibleBitmap")
	hbitmap, _, err := proc.Call(uintptr(hdc), uintptr(width), uintptr(height))
	if hbitmap == 0 {
		return 0, 0, fmt.Errorf("failed to create compatible bitmap")
	}
	return syscall.Handle(hbitmap), 0, err
}

func createCompatibleDC(hdc syscall.Handle) (syscall.Handle, uintptr, error) {
	gdi32 := syscall.NewLazyDLL("gdi32.dll")
	proc := gdi32.NewProc("CreateCompatibleDC")
	hdcPtr, _, err := proc.Call(uintptr(hdc))
	if hdcPtr == 0 {
		return 0, 0, fmt.Errorf("failed to create compatible device context")
	}
	return syscall.Handle(hdcPtr), 0, err
}

func selectObject(hdc syscall.Handle, hObject syscall.Handle) (syscall.Handle, uintptr, error) {
	gdi32 := syscall.NewLazyDLL("gdi32.dll")
	proc := gdi32.NewProc("SelectObject")
	oldObj, _, err := proc.Call(uintptr(hdc), uintptr(hObject))
	if oldObj == 0 {
		return 0, 0, fmt.Errorf("failed to select object into device context")
	}
	return syscall.Handle(oldObj), 0, err
}

func bitblt(hdcDest syscall.Handle, nXDest, nYDest, nWidth, nHeight int, hdcSrc syscall.Handle, nXSrc, nYSrc int, dwRop uint32) (bool, uintptr, error) {
	gdi32 := syscall.NewLazyDLL("gdi32.dll")
	proc := gdi32.NewProc("BitBlt")
	success, _, err := proc.Call(uintptr(hdcDest), uintptr(nXDest), uintptr(nYDest), uintptr(nWidth), uintptr(nHeight), uintptr(hdcSrc), uintptr(nXSrc), uintptr(nYSrc), uintptr(dwRop))
	if success == 0 {
		return false, 0, fmt.Errorf("failed to copy screen contents into bitmap")
	}
	return success != 0, 0, err
}

func createImageFromBitmap(hBitmap syscall.Handle, width, height int) (*image.RGBA, error) {
	gdiplus := syscall.NewLazyDLL("gdiplus.dll")

	var gdiplusToken uintptr
	startupInput := newGdiplusStartupInput()
	ret, _, _ := gdiplus.NewProc("GdiplusStartup").Call(uintptr(unsafe.Pointer(&gdiplusToken)), uintptr(unsafe.Pointer(startupInput)), 0)
	if ret != 0 {
		return nil, fmt.Errorf("failed to initialize GDI+: %d", ret)
	}
	defer gdiplus.NewProc("GdiplusShutdown").Call(gdiplusToken)

	var imagePtr uintptr
	ret, _, _ = gdiplus.NewProc("GdipCreateBitmapFromHBITMAP").Call(uintptr(hBitmap), 0, uintptr(unsafe.Pointer(&imagePtr)))
	if ret != 0 {
		return nil, fmt.Errorf("failed to create GDI+ bitmap: %d", ret)
	}
	defer gdiplus.NewProc("GdipDisposeImage").Call(imagePtr)

	var bitmapData gdiplusBitmapData
	ret, _, _ = gdiplus.NewProc("GdipBitmapLockBits").Call(imagePtr, uintptr(unsafe.Pointer(&imageRect)), 3, 0x26200A, uintptr(unsafe.Pointer(&bitmapData)))
	if ret != 0 {
		return nil, fmt.Errorf("failed to lock GDI+ bitmap bits: %d", ret)
	}
	defer gdiplus.NewProc("GdipBitmapUnlockBits").Call(imagePtr, uintptr(unsafe.Pointer(&bitmapData)))

	imageBytes := make([]byte, int(bitmapData.Stride)*int(bitmapData.Height))
	for y := 0; y < int(bitmapData.Height); y++ {
		row := (*[1 << 30]byte)(unsafe.Pointer(bitmapData.Scan0 + uintptr(y)*uintptr(bitmapData.Stride)))[:int(bitmapData.Stride)]
		copy(imageBytes[y*int(bitmapData.Stride):(y+1)*int(bitmapData.Stride)], row)
	}

	// Create a new RGBA image and set alpha channel to 255 (fully opaque)
	rgbaImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < len(imageBytes); i += 4 {
		rgbaImage.Pix[i] = imageBytes[i+2]   // Blue channel
		rgbaImage.Pix[i+1] = imageBytes[i+1] // Green channel
		rgbaImage.Pix[i+2] = imageBytes[i]   // Red channel
		rgbaImage.Pix[i+3] = 255             // Set alpha channel to 255
	}

	return rgbaImage, nil
}

func saveImageAsPNG(img *image.RGBA, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func saveImageAsJPEG(img *image.RGBA, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, nil)
}

type gdiplusStartupInput struct {
	GdiplusVersion           uint32
	DebugEventCallback       uintptr
	SuppressBackgroundThread uint32
	SuppressExternalCodecs   uint32
}

func newGdiplusStartupInput() *gdiplusStartupInput {
	return &gdiplusStartupInput{
		GdiplusVersion: 1,
	}
}

type gdiplusRect struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

var imageRect = gdiplusRect{
	X:      0,
	Y:      0,
	Width:  int32(width),
	Height: int32(height),
}

type gdiplusBitmapData struct {
	Width       uint32
	Height      uint32
	Stride      int32
	PixelFormat uint32
	Scan0       uintptr
	Reserved    uint32
}

func deleteObject(hObject syscall.Handle) (uintptr, uintptr, error) {
	gdi32 := syscall.NewLazyDLL("gdi32.dll")
	proc := gdi32.NewProc("DeleteObject")
	result, _, err := proc.Call(uintptr(hObject))
	if result == 0 {
		return 0, 0, fmt.Errorf("failed to delete object")
	}
	return result, 0, err
}

func deleteDC(hdc syscall.Handle) (uintptr, uintptr, error) {
	gdi32 := syscall.NewLazyDLL("gdi32.dll")
	proc := gdi32.NewProc("DeleteDC")
	result, _, err := proc.Call(uintptr(hdc))
	if result == 0 {
		return 0, 0, fmt.Errorf("failed to delete device context")
	}
	return result, 0, err
}

func getColorAtPixel(img *image.RGBA, x, y int) (uint8, uint8, uint8, uint8) {
	idx := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)*4
	return img.Pix[idx], img.Pix[idx+1], img.Pix[idx+2], img.Pix[idx+3]
}

func loadImage(filename string) (*image.RGBA, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	rgbaImage := image.NewRGBA(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rgbaImage.SetRGBA(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}

	return rgbaImage, nil
}
