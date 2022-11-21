// part of a complementary module goslowc2 as a module in a blog wrapping as a subprocess and IOCs

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32dll_obj   = syscall.NewLazyDLL("user32.dll")
	messageboxa_obj = user32dll_obj.NewProc("MessageBoxA")
)

func main() {
	//Win11x64 Go 1.18.4 ptr call grabs first byte?
	foo_obj, _ := syscall.UTF16PtrFromString("foo")
	bar_obj, _ := syscall.UTF16PtrFromString("bar")
	//barj, _ := syscall.UTF16PtrFromString.('bar')
	fmt.Println("calling user32.dll function now..")
	//messageboxa_obj.Call(0, uintptr(unsafe.Pointer(foo_obj)), uintptr(unsafe.Pointer(bar)), 0x1L)
	//messageboxa_obj.Call(0, 0, 0, 0)
	messageboxa_obj.Call(0, uintptr(unsafe.Pointer(foo_obj)), uintptr(unsafe.Pointer(bar_obj)), 0)
	//if err != nil {
	//	fmt.Println(err)
	//}

}

// Very useful references that helped me understand GoLang structures and future ideas
// https://anubissec.github.io/How-To-Call-Windows-APIs-In-Golang/#
// https://justen.codes/breaking-all-the-rules-using-go-to-call-windows-api-2cbfd8c79724
// https://www.thesubtlety.com/post/getting-started-golang-windows-apis/
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messagebox
