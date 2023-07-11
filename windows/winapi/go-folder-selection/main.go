package main

import (
	"log"
	"syscall"
	"unsafe"
)

type BROWSEINFO struct {
	hwndOwner      syscall.Handle
	pidlRoot       uintptr
	pszDisplayName *uint16
	lpszTitle      *uint16
	ulFlags        uint32
	lpfn           uintptr
	lParam         uintptr
	iImage         int32
}

func main() {
	var bi BROWSEINFO
	bi.lpszTitle, _ = syscall.UTF16PtrFromString("Select Folder")
	bi.ulFlags = 0x00000001 | 0x00000040 // BIF_RETURNONLYFSDIRS | BIF_USENEWUI

	// Get a handle to the shell32.dll library
	shell32 := syscall.NewLazyDLL("shell32.dll")

	// Get a handle to the SHBrowseForFolderW function
	shBrowseForFolder := shell32.NewProc("SHBrowseForFolderW")

	// Call the SHBrowseForFolderW function
	pidl, _, _ := shBrowseForFolder.Call(uintptr(unsafe.Pointer(&bi)))

	if pidl != 0 {
		var path [syscall.MAX_PATH]uint16

		// Get the selected folder path using SHGetPathFromIDListW function
		shGetPathFromIDList := shell32.NewProc("SHGetPathFromIDListW")
		_, _, _ = shGetPathFromIDList.Call(pidl, uintptr(unsafe.Pointer(&path[0])))

		// Convert the UTF-16 encoded path to a Go string
		folderPath := syscall.UTF16ToString(path[:])

		log.Println("Selected folder: ", folderPath)

		// Free the allocated PIDL using CoTaskMemFree function from ole32.dll
		ole32 := syscall.NewLazyDLL("ole32.dll")
		coTaskMemFree := ole32.NewProc("CoTaskMemFree")
		_, _, _ = coTaskMemFree.Call(pidl)
	} else {
		log.Println("Folder selection canceled.")
	}
}
