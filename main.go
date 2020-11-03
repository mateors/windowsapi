package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32DLL           = windows.NewLazyDLL("user32.dll")
	procSystemParamInfo = user32DLL.NewProc("SystemParametersInfoW")

	//modKernel32     = syscall.NewLazyDLL("kernel32.dll")
	//modComctl = syscall.NewLazyDLL("comctl32.dll")
	//editGetText = modComctl.NewProc("Edit_GetText")
	modUICore        = syscall.NewLazyDLL("UIAutomationCore.dll")
	findFirst        = modUICore.NewProc("GetSupportedPatterns")
	enumChildWindows = user32DLL.NewProc("EnumChildWindows")
)

func main() {

	fmt.Println(enumChildWindows.Find())

}

func changeDesktopBackground() {

	//p := unsafe.Pointer(&imagePath)
	//SPI_SETDESKWALLPAPER=
	// BOOL SystemParametersInfoW(
	// 	UINT  uiAction,
	// 	UINT  uiParam,
	// 	PVOID pvParam,
	// 	UINT  fWinIni
	//   );
	//uiAction=Desktop parameters=SPI_SETDESKWALLPAPER=0x0014
	//http://redgrittybrick.org/ascii.html decimal 20 = 0x0014
	//fWinIni=decimal 26 = 0x001A
	//procSystemParamInfo.Call(20, 0, uintptr(unsafe.Pointer(imagePath)), 0x001A)
	//x := unsafe.Pointer(imagePath)
	//fmt.Println(x)
	// u := unsafe.Pointer(nil)
	// p := unsafe.Pointer(uintptr(u) + imagePath)
	//p := (*int)(unsafe.Pointer(u))

	imagePath, _ := windows.UTF16PtrFromString(`E:\GOLANG\src\mateors\windowsapi\winback.jpg`)
	_, _, err := procSystemParamInfo.Call(20, 0, uintptr(unsafe.Pointer(imagePath)), 0x001A)
	fmt.Println(err)

}
